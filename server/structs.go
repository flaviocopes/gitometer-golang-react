package main

// week represents the summary of a week of activity
// on a repository
type week struct {
	ID           int
	RepositoryID int
	WeekNumber   int
	Year         int
	CreatedOn    string
	IssuesClosed int
	IssuesOpened int
	Stars        int
	Commits      int
	WeekStart    string
	WeekEnd      string
	PrOpened     int
	PrMerged     int
	PrClosed     int
}

// timeline represents important events happened on a
// repository, which will be displayed on the repo timeline
type timeline struct {
	ID           int
	RepositoryID int
	Title        string
	Description  string
	Emoji        string
	Date         string
}

// repository contains the details of a repository
type repository struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	OwnerName                string
	RepoAge                  int    `json:"repository_created_months_ago"`
	Initialized              bool   `json:"initialized"`
	CommitsPerMonth          string `json:"commits_per_month"`
	StarsPerMonth            string `json:"stars_per_month"`
	TotalStars               int    `json:"total_stars"`
	TotalCommits             int    `json:"total_commits"`
	Description              string `json:"description"`
	CommitsCountLast12Months int    `json:"commits_count_last_12_months"`
	CommitsCountLast4Weeks   int    `json:"commits_count_last_4_weeks"`
	CommitsCountLastWeek     int    `json:"commits_count_last_week"`
	StarsCountLast12Months   int    `json:"stars_count_last_12_months"`
	StarsCountLast4Weeks     int    `json:"stars_count_last_4_weeks"`
	StarsCountLastWeek       int    `json:"stars_count_last_week"`
}

// owner contains the details of an owner or a repo
type owner struct {
	ID                  int
	Name                string
	Description         string
	Avatar              string `json:"avatar"`
	GitHubID            string
	AddedBy             string
	Enabled             bool
	InstallationID      string
	RepositorySelection string
}

// repoData contains the aggregate repository data returned
// by the API call
type repoData struct {
	MonthlyData monthlyData  `json:"monthly_data"`
	WeeklyData  []week       `json:"weekly_data"`
	Years       map[int]bool `json:"years"`
	Timeline    []timeline   `json:"timeline"`
	Repository  repository   `json:"repository"`
	Owner       owner        `json:"owner"`
}

// monthlyData contains the monthly activity of a repo
type monthlyData struct {
	CommitsPerMonth string `json:"commits_per_month"`
	StarsPerMonth   string `json:"stars_per_month"`
}

// Error handling types

type errRepoNotInitialized string

func (e errRepoNotInitialized) Error() string {
	return string(e)
}

type errRepoNotFound string

func (e errRepoNotFound) Error() string {
	return string(e)
}

// repositorySummary contains the details of a repository
type repositorySummary struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	OwnerName  string `json:"ownerName"`
	TotalStars int    `json:"totalStars"`
}

// repositories contains a slice of repositories
type repositories struct {
	Repositories []repositorySummary `json:"repositories"`
}
