package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	initDb()
	defer db.Close()
	http.HandleFunc("/api/index", indexHandler)
	http.HandleFunc("/api/repo/", repoHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// indexHandler calls `queryRepos()` and marshals the result as JSON
func indexHandler(w http.ResponseWriter, req *http.Request) {
	repos := repositories{}

	err := queryRepos(&repos)
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
	repo := repository{}
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
		case errRepoNotFound:
			http.Error(w, err.Error(), 404)
		case errRepoNotInitialized:
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
