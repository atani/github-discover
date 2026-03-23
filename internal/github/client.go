package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		baseURL:    defaultBaseURL,
	}
}

// NewTestClient creates a client pointing to a custom base URL (for testing).
func NewTestClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.httpClient.Do(req)
}

// SearchRepositories searches GitHub repositories with the given query and options.
func (c *Client) SearchRepositories(query string, sort string, order string, perPage int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	if sort != "" {
		params.Set("sort", sort)
	}
	if order != "" {
		params.Set("order", order)
	}
	if perPage > 0 {
		params.Set("per_page", fmt.Sprintf("%d", perPage))
	}

	reqURL := fmt.Sprintf("%s/search/repositories?%s", c.baseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(body))
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRepository fetches detailed information about a specific repository.
func (c *Client) GetRepository(owner, repo string) (*Repository, error) {
	reqURL := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(body))
	}

	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}

	return &repository, nil
}

// GetTrending fetches trending repositories by looking at recently created repos with high star counts.
func (c *Client) GetTrending(language string, since string, perPage int) (*SearchResult, error) {
	var dateRange string
	now := time.Now()

	switch since {
	case "weekly":
		dateRange = now.AddDate(0, 0, -7).Format("2006-01-02")
	case "monthly":
		dateRange = now.AddDate(0, -1, 0).Format("2006-01-02")
	default: // daily
		dateRange = now.AddDate(0, 0, -1).Format("2006-01-02")
	}

	var queryParts []string
	queryParts = append(queryParts, fmt.Sprintf("created:>%s", dateRange))
	if language != "" {
		queryParts = append(queryParts, fmt.Sprintf("language:%s", language))
	}

	query := strings.Join(queryParts, " ")
	return c.SearchRepositories(query, "stars", "desc", perPage)
}

// GetTrendingByStars fetches repos that gained the most stars recently.
func (c *Client) GetTrendingByStars(language string, since string, perPage int) (*SearchResult, error) {
	var dateRange string
	now := time.Now()

	switch since {
	case "weekly":
		dateRange = now.AddDate(0, 0, -7).Format("2006-01-02")
	case "monthly":
		dateRange = now.AddDate(0, -1, 0).Format("2006-01-02")
	default: // daily
		dateRange = now.AddDate(0, 0, -1).Format("2006-01-02")
	}

	var queryParts []string
	queryParts = append(queryParts, fmt.Sprintf("pushed:>%s", dateRange))
	queryParts = append(queryParts, "stars:>10")
	if language != "" {
		queryParts = append(queryParts, fmt.Sprintf("language:%s", language))
	}

	query := strings.Join(queryParts, " ")
	return c.SearchRepositories(query, "stars", "desc", perPage)
}
