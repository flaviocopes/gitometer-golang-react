package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/flaviocopes/gitometer/server/common"
	"github.com/flaviocopes/gitometer/server/db"
	"github.com/flaviocopes/gitometer/server/github"
)

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			//handle preflight in here
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func main() {
	db.InitDb()
	defer db.Close()

	http.HandleFunc("/api/index", indexHandler)
	http.HandleFunc("/api/repo/", getRepoHandler)
	http.HandleFunc("/api/repo", addRepoHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// parseParams accepts a req and returns the `num` path tokens found after the `prefix`.
// returns an error if the number of tokens are less or more than expected
func parseParams(req *http.Request, prefix string, num int) ([]string, error) {
	url := strings.TrimPrefix(req.URL.Path, prefix)
	params := strings.Split(url, "/")
	if len(params) != num || len(params[0]) == 0 || len(params[1]) == 0 {
		return nil, fmt.Errorf("Bad format. Expecting exactly %d params", num)
	}
	return params, nil
}

// indexHandler calls `queryRepos()` and marshals the result as JSON
func indexHandler(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
	if (*req).Method == "OPTIONS" {
		return
	}

	repos := common.Repositories{}

	err := db.QueryRepos(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

// queryRepo first fetches the repository, and if nothing is wrong
// it returns the result of fetchData()
func queryRepo(repo *common.Repository) (*common.RepoData, error) {
	data := common.RepoData{}
	err := db.FetchRepo(repo, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// getRepoHandler processes the response by parsing the params, then calling
// `query()`, and marshaling the result in JSON format, sending it to
// `http.ResponseWriter`.
func getRepoHandler(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
	if (*req).Method == "OPTIONS" {
		return
	}
	switch req.Method {
	case "GET":
		handleGetRepo(w, req)
	}
}

func addRepoHandler(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
	if (*req).Method == "OPTIONS" {
		return
	}
	handleAddNewRepo(w, req)
}

type newRepoData struct {
	Owner string
	Name  string
}

func handleAddNewRepo(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var data newRepoData
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	owner := data.Owner
	name := data.Name

	if owner == "" || name == "" {
		http.Error(w, "Missing parameter name or owner", 500)
	}

	github.AddRepoToDb(owner, name)

	fmt.Fprintf(w, string("ok"))
}

func handleGetRepo(w http.ResponseWriter, req *http.Request) {
	repo := common.Repository{}
	params, err := parseParams(req, "/api/repo/", 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	repo.OwnerName = params[0]
	repo.Name = params[1]

	data, err := queryRepo(&repo)
	if err != nil {
		switch err.(type) {
		case common.ErrRepoNotFound:
			http.Error(w, err.Error(), 404)
		case common.ErrRepoNotInitialized:
			http.Error(w, err.Error(), 401)
		default:
			http.Error(w, err.Error(), 500)
		}
		return
	}

	out, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}
