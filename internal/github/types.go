package github

import "time"

type SearchResult struct {
	TotalCount int          `json:"total_count"`
	Items      []Repository `json:"items"`
}

type Repository struct {
	FullName        string    `json:"full_name"`
	Name            string    `json:"name"`
	Owner           Owner     `json:"owner"`
	Description     string    `json:"description"`
	HTMLURL         string    `json:"html_url"`
	Homepage        string    `json:"homepage"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	WatchersCount   int       `json:"watchers_count"`
	Topics          []string  `json:"topics"`
	License         *License  `json:"license"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Archived        bool      `json:"archived"`
	Fork            bool      `json:"fork"`
}

type Owner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
}
