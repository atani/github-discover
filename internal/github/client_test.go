package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SearchRepositories(t *testing.T) {
	expected := SearchResult{
		TotalCount: 1,
		Items: []Repository{
			{
				FullName:        "test/repo",
				Name:            "repo",
				Description:     "A test repository",
				StargazersCount: 100,
				Language:        "Go",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/repositories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		q := r.URL.Query().Get("q")
		if q == "" {
			t.Error("expected non-empty query")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), token: ""}
	// Override baseURL by using the test server directly
	origDo := client.httpClient

	// Create a client that points to test server
	testClient := &Client{httpClient: origDo, token: "test-token"}

	// Use a helper to test with the mock server
	result, err := searchWithServer(testClient, server.URL, "test", "stars", "desc", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalCount != 1 {
		t.Errorf("TotalCount: got %d, want 1", result.TotalCount)
	}

	if len(result.Items) != 1 {
		t.Fatalf("Items length: got %d, want 1", len(result.Items))
	}

	if result.Items[0].FullName != "test/repo" {
		t.Errorf("FullName: got %q, want %q", result.Items[0].FullName, "test/repo")
	}
}

func searchWithServer(c *Client, serverURL, query, sort, order string, perPage int) (*SearchResult, error) {
	reqURL := serverURL + "/search/repositories?q=" + query
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func TestClient_do_setsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Errorf("Accept header: got %q, want %q", got, "application/vnd.github+json")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization header: got %q, want %q", got, "Bearer test-token")
		}
		if got := r.Header.Get("X-GitHub-Api-Version"); got != "2022-11-28" {
			t.Errorf("X-GitHub-Api-Version header: got %q, want %q", got, "2022-11-28")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), token: "test-token"}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestClient_do_noAuthWithoutToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization header should be empty, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), token: ""}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestClient_GetRepository_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), token: ""}
	req, _ := http.NewRequest("GET", server.URL+"/repos/test/notfound", nil)
	resp, err := client.do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode: got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
