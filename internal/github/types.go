package github

import "time"

// Repository contains the subset of GitHub repository metadata RepoRadar displays.
type Repository struct {
	Name          string   `json:"name"`
	FullName      string   `json:"full_name"`
	Description   string   `json:"description"`
	HTMLURL       string   `json:"html_url"`
	Stars         int      `json:"stargazers_count"`
	Forks         int      `json:"forks_count"`
	OpenIssues    int      `json:"open_issues_count"`
	DefaultBranch string   `json:"default_branch"`
	License       *License `json:"license"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// License represents the compact license object returned by GitHub.
type License struct {
	Name string `json:"name"`
}

// Analysis is the API response consumed by the RepoRadar frontend.
type Analysis struct {
	Name          string         `json:"name"`
	FullName      string         `json:"fullName"`
	Description   string         `json:"description"`
	URL           string         `json:"url"`
	Stars         int            `json:"stars"`
	Forks         int            `json:"forks"`
	OpenIssues    int            `json:"openIssues"`
	DefaultBranch string         `json:"defaultBranch"`
	License       string         `json:"license"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	Languages     map[string]int `json:"languages"`
	TotalBytes    int            `json:"totalBytes"`
}
