package common

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

// Repository contains the details of a repository
type Repository struct {
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

// RepoData contains the aggregate repository data returned
// by the API call
type RepoData struct {
	MonthlyData MonthlyData  `json:"monthly_data"`
	WeeklyData  []Week       `json:"weekly_data"`
	Years       map[int]bool `json:"years"`
	Repository  Repository   `json:"repository"`
	Owner       owner        `json:"owner"`
}

// MonthlyData contains the monthly activity of a repo
type MonthlyData struct {
	CommitsPerMonth string `json:"commits_per_month"`
	StarsPerMonth   string `json:"stars_per_month"`
}

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
