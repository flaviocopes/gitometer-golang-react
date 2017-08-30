package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Postgres drivers
	_ "github.com/lib/pq"

	"github.com/flaviocopes/gitometer/server/common"
)

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

var db *sql.DB

// InitDb ...
func InitDb() {
	if db == nil {
		config := dbConfig()
		var err error
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password='%s' dbname=%s sslmode=disable",
			config[dbhost], config[dbport],
			config[dbuser], config[dbpass], config[dbname])

		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully connected!")
	}
}

// Close closes the db connection (called via defer)
func Close() {
	db.Close()
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

// QueryRepos first fetches the repositories data from the db
func QueryRepos(repos *common.Repositories) error {
	rows, err := db.Query(`
		SELECT
			id,
			repository_owner,
			repository_name,
			total_stars
		FROM repositories
		ORDER BY total_stars DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		repo := common.RepositorySummary{}
		err = rows.Scan(
			&repo.ID,
			&repo.OwnerName,
			&repo.Name,
			&repo.TotalStars,
		)
		if err != nil {
			return err
		}
		repos.Repositories = append(repos.Repositories, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// AddNewRepo adds a repository to the db
func AddNewRepo(owner, name string, repo *common.Repository) error {
	var id int

	err := db.QueryRow("SELECT id_of_repository_on_github FROM repositories WHERE repository_owner=$1 AND repository_name=$2", owner, name).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		// Repo is new

		sqlStatement := `
			INSERT INTO repositories (
				id_of_repository_on_github,
				repository_name,
				repository_owner,
				default_branch,
				created_at,
				description,
				repository_created_months_ago,
				total_stars,
				total_commits,
				commits_count_last_12_months,
				commits_count_last_4_weeks,
				commits_count_last_week,
				stars_count_last_12_months,
				stars_count_last_4_weeks,
				stars_count_last_week,
				stars_per_month
				)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
		_, err := db.Exec(
			sqlStatement,
			repo.ID,
			repo.Name,
			repo.OwnerName,
			repo.DefaultBranch,
			repo.CreatedAt,
			repo.Description,
			repo.RepoAge,
			repo.TotalStars,
			repo.TotalCommits,
			repo.CommitsCountLast12Months,
			repo.CommitsCountLast4Weeks,
			repo.CommitsCountLastWeek,
			repo.StarsCountLast12Months,
			repo.StarsCountLast4Weeks,
			repo.StarsCountLastWeek,
			repo.StarsPerMonth,
		)

		if err != nil {
			log.Fatal(err)
		}

	case err != nil:
		log.Fatal(err)
	default:
		sqlStatement := `
			UPDATE repositories SET
				stars_per_month = $1,
				default_branch = $2,
				description = $3,
				repository_created_months_ago = $4,
				total_stars = $5,
				total_commits = $6,
				commits_count_last_12_months = $7,
				commits_count_last_4_weeks = $8,
				commits_count_last_week = $9,
				stars_count_last_12_months = $10,
				stars_count_last_4_weeks = $11,
				stars_count_last_week = $12
			WHERE id_of_repository_on_github = $13`
		_, err := db.Exec(
			sqlStatement,
			repo.StarsPerMonth,
			repo.DefaultBranch,
			repo.Description,
			repo.RepoAge,
			repo.TotalStars,
			repo.TotalCommits,
			repo.CommitsCountLast12Months,
			repo.CommitsCountLast4Weeks,
			repo.CommitsCountLastWeek,
			repo.StarsCountLast12Months,
			repo.StarsCountLast4Weeks,
			repo.StarsCountLastWeek,
			id,
		)

		if err != nil {
			log.Fatal(err)
		}

	}

	return nil
}

// FetchRepo given a Repository value with name and owner of the repo
// fetches more details from the database and fills the value with more
// data
func FetchRepo(repo *common.Repository, data *common.RepoData) error {
	if len(repo.Name) == 0 {
		return fmt.Errorf("Repository name not correctly set")
	}
	if len(repo.OwnerName) == 0 {
		return fmt.Errorf("Repository owner not correctly set")
	}
	sqlStatement := `
		SELECT
			id,
			initialized,
			repository_created_months_ago,
			total_stars,
			total_commits,
			description,
			commits_count_last_12_months,
			commits_count_last_4_weeks,
			commits_count_last_week,
			stars_count_last_12_months,
			stars_count_last_4_weeks,
			stars_count_last_week,
			stars_per_month
		FROM repositories
		WHERE repository_owner=$1 and repository_name=$2
		LIMIT 1;`
	row := db.QueryRow(sqlStatement, repo.OwnerName, repo.Name)
	err := row.Scan(
		&repo.ID,
		&repo.Initialized,
		&repo.RepoAge,
		&repo.TotalStars,
		&repo.TotalCommits,
		&repo.Description,
		&repo.CommitsCountLast12Months,
		&repo.CommitsCountLast4Weeks,
		&repo.CommitsCountLastWeek,
		&repo.StarsCountLast12Months,
		&repo.StarsCountLast4Weeks,
		&repo.StarsCountLastWeek,
		&repo.StarsPerMonth)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			//locally handle SQL error, abstract for caller
			return common.ErrRepoNotFound("Repository not found")
		default:
			return err
		}
	}
	if !repo.Initialized {
		return common.ErrRepoNotInitialized("Repository not initialized")
	}
	if repo.RepoAge < 3 {
		return common.ErrRepoNotInitialized("Repository not initialized")
	}

	// assign to data
	data.Repository = *repo

	return nil
}
