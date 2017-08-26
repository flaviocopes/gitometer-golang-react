package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

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

// queryRepo first fetches the repository, and if nothing is wrong
// it returns the result of fetchData()
func queryRepo(repo *repository) (*repoData, error) {
	data := repoData{}
	err := fetchRepo(repo, &data)
	if err != nil {
		return nil, err
	}
	return fetchData(repo, &data)
}

// fetchData calls utility functions to collect data from
// the database, builds and returns the `RepoData` value
func fetchData(repo *repository, data *repoData) (*repoData, error) {
	err := fetchMonthlyData(repo, data)
	if err != nil {
		return nil, err
	}
	err = fetchWeeklyData(repo, data)
	if err != nil {
		return nil, err
	}
	err = fetchYearlyData(repo, data)
	if err != nil {
		return nil, err
	}
	err = fetchTimelineData(repo, data)
	if err != nil {
		return nil, err
	}
	err = fetchOwnerData(repo, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// fetchRepo given a Repository value with name and owner of the repo
// fetches more details from the database and fills the value with more
// data
func fetchRepo(repo *repository, data *repoData) error {
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
			return errRepoNotFound("Repository not found")
		default:
			return err
		}
	}
	if !repo.Initialized {
		return errRepoNotInitialized("Repository not initialized")
	}
	if repo.RepoAge < 3 {
		return errRepoNotInitialized("Repository not initialized")
	}

	// assign to data
	data.Repository = *repo

	return nil
}

// fetchOwnerData given a Repository object with the `Owner` value
// it fetches information about it from the database
func fetchOwnerData(repo *repository, data *repoData) error {
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

// fetchMonthlyData given a repository ID, it fetches the monthly
// data information
func fetchMonthlyData(repo *repository, data *repoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	data.MonthlyData = monthlyData{}
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

// fetchWeeklyData given a repository ID, it fetches the weekly
// data information
func fetchWeeklyData(repo *repository, data *repoData) error {
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
		week := week{}
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

// fetchYearlyData returns the list of years for which we have weekly data
// available
func fetchYearlyData(repo *repository, data *repoData) error {
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

// fetchTimelineData returns all the timeline data we have in the db about
// the repo
func fetchTimelineData(repo *repository, data *repoData) error {
	if repo.ID == 0 {
		return fmt.Errorf("Repository ID not correctly set")
	}
	rows, err := db.Query(`
        SELECT
            id,
            repository_id,
            title,
            description,
            emoji,
            date
        FROM repositories_timelines
        WHERE repository_id=$1
        ORDER BY date ASC`, repo.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		timeline := timeline{}
		err = rows.Scan(
			&timeline.ID,
			&timeline.RepositoryID,
			&timeline.Title,
			&timeline.Description,
			&timeline.Emoji,
			&timeline.Date)
		if err != nil {
			return err
		}
		data.Timeline = append(data.Timeline, timeline)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
