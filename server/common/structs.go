package common

// Repository contains the details of a repository
type Repository struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	OwnerName                string `json:"ownerName"`
	RepoAge                  int    `json:"repository_created_months_ago"`
	Initialized              bool   `json:"initialized"`
	TotalStars               int    `json:"total_stars"`
	TotalCommits             int    `json:"total_commits"`
	Description              string `json:"description"`
	CreatedAt                string `json:"created_at"`
	DefaultBranch            string `json:"default_branch"`
	CommitsCountLast12Months int    `json:"commits_count_last_12_months"`
	CommitsCountLast4Weeks   int    `json:"commits_count_last_4_weeks"`
	CommitsCountLastWeek     int    `json:"commits_count_last_week"`
	StarsCountLast12Months   int    `json:"stars_count_last_12_months"`
	StarsCountLast4Weeks     int    `json:"stars_count_last_4_weeks"`
	StarsCountLastWeek       int    `json:"stars_count_last_week"`
	StarsPerMonth            string `json:"stars_per_month"`
}

// RepoData contains the aggregate repository data returned
// by the API call
type RepoData struct {
	Repository Repository `json:"repository"`
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
