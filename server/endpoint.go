package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/flaviocopes/gitometer/server/common"
	"github.com/flaviocopes/gitometer/server/db"
)

func main() {
	db.InitDb()
	defer db.Close()
	http.HandleFunc("/api/index", indexHandler)
	http.HandleFunc("/api/repo/", repoHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

	setupResponse(&w)
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
	return fetchData(repo, &data)
}

// fetchYearlyData returns the list of years for which we have weekly data
// available
func fetchYearlyData(repo *common.Repository, data *common.RepoData) error {
	if data.WeeklyData == nil {
		return fmt.Errorf("Repository weekly data not correctly set")
	}
	data.Years = make(map[int]bool)
	for i := 0; i < len(data.WeeklyData); i++ {
		year := data.WeeklyData[i].Year
		data.Years[year] = true
	}
	return nil
}

// fetchData calls utility functions to collect data from
// the database, builds and returns the `RepoData` value
func fetchData(repo *common.Repository, data *common.RepoData) (*common.RepoData, error) {
	err := db.FetchMonthlyData(repo, data)
	if err != nil {
		return nil, err
	}
	err = db.FetchWeeklyData(repo, data)
	if err != nil {
		return nil, err
	}
	err = fetchYearlyData(repo, data)
	if err != nil {
		return nil, err
	}
	err = db.FetchOwnerData(repo, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// repoHandler processes the response by parsing the params, then calling
// `query()`, and marshaling the result in JSON format, sending it to
// `http.ResponseWriter`.
func repoHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		handleGetRepo(w, req)
	case "POST":
		handleAddNewRepo(w, req)
	}
}

func handleAddNewRepo(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w)

	fmt.Fprintf(w, "ADD")
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

	setupResponse(&w)
	fmt.Fprintf(w, string(out))
}
