package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/go-github/github"
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

// InitDb ...
func InitDb() {
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
func AddNewRepo(repo *github.Repository) error {
	panic("Added")
	// sqlStatement := `
	// 	INSERT INTO repositories (age, email, first_name, last_name)
	// 	VALUES ($1, $2, $3, $4)`
	// _, err := db.Exec(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun")
	// if err != nil {
	// 	return err
	// }

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
			stars_count_last_week
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
		&repo.StarsCountLastWeek)
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

// FetchOwnerData given a Repository object with the `Owner` value
// it fetches information about it from the database
func FetchOwnerData(repo *common.Repository, data *common.RepoData) error {
	if len(repo.OwnerName) == 0 {
		return fmt.Errorf("Repository owner not correctly set")
	}
	sqlStatement := `
        SELECT
            id,
            name,
            COALESCE(description, ''),
            COALESCE(avatar_url, ''),
            COALESCE(github_id, ''),
            added_by,
            enabled,
            COALESCE(installation_id, ''),
            repository_selection
        FROM organizations
        WHERE name=$1
        ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(sqlStatement, repo.OwnerName)
	err := row.Scan(&data.Owner.ID,
		&data.Owner.Name,
		&data.Owner.Description,
		&data.Owner.Avatar,
		&data.Owner.GitHubID,
		&data.Owner.AddedBy,
		&data.Owner.Enabled,
		&data.Owner.InstallationID,
		&data.Owner.RepositorySelection)
	if err != nil {
		return err
	}
	return nil
}

// FetchMonthlyData given a repository ID, it fetches the monthly
// data information
func FetchMonthlyData(repo *common.Repository, data *common.RepoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	data.MonthlyData = common.MonthlyData{}
	sqlStatement := `
        SELECT
            commits_per_month,
            stars_per_month
        FROM repositories_historic_data
        WHERE repository_id=$1
        ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(sqlStatement, repo.ID)
	err := row.Scan(
		&data.MonthlyData.CommitsPerMonth,
		&data.MonthlyData.StarsPerMonth)
	if err != nil {
		return err
	}

	return nil
}

// FtchWeeklyData given a repository ID, it fetches the weekly
// data information
func FetchWeeklyData(repo *common.Repository, data *common.RepoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	rows, err := db.Query(`
        SELECT
            id,
            repository_id,
            week_number,
            year,
            created_on,
            issues_closed,
            issues_opened,
            stars,
            commits,
            week_start,
            week_end,
            pr_opened,
            pr_merged,
            pr_closed
        FROM repositories_weekly_data
        WHERE repository_id=$1
        ORDER BY id ASC`, repo.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		week := common.Week{}
		err = rows.Scan(
			&week.ID,
			&week.RepositoryID,
			&week.WeekNumber,
			&week.Year,
			&week.CreatedOn,
			&week.IssuesClosed,
			&week.IssuesOpened,
			&week.Stars,
			&week.Commits,
			&week.WeekStart,
			&week.WeekEnd,
			&week.PrOpened,
			&week.PrMerged,
			&week.PrClosed)
		if err != nil {
			return err
		}
		data.WeeklyData = append(data.WeeklyData, week)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
