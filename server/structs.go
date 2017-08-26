package main

// RepositorySummary contains the details of a repository
type RepositorySummary struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	OwnerName  string `json:"ownerName"`
	TotalStars int    `json:"totalStars"`
}

// Repositories contains a slice of repositories
type Repositories struct {
	Repositories []RepositorySummary `json:"repositories"`
}

// Week represents the summary of a week of activity
// on a repository
type Week struct {
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

// Timeline represents important events happened on a
// repository, which will be displayed on the repo timeline
type Timeline struct {
	ID           int
	RepositoryID int
	Title        string
	Description  string
	Emoji        string
	Date         string
}

// Owner contains the details of an owner or a repo
type Owner struct {
	ID                  int
	Name                string
	Description         string
	Avatar              string
	GitHubID            string
	AddedBy             string
	Enabled             bool
	InstallationID      string
	RepositorySelection string
}

// MonthlyData contains the monthly activity of a repo
type MonthlyData struct {
	CommitsPerMonth string
	StarsPerMonth   string
}
