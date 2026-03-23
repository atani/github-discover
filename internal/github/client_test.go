package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)
	client := NewTestClient(server.URL)
	return server, client
}

func TestNewClient(t *testing.T) {
	c := NewClient("my-token")
	if c.token != "my-token" {
		t.Errorf("token: got %q, want %q", c.token, "my-token")
	}
	if c.baseURL != defaultBaseURL {
		t.Errorf("baseURL: got %q, want %q", c.baseURL, defaultBaseURL)
	}
}

func TestClient_do_setsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Errorf("Accept: got %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization: got %q", got)
		}
		if got := r.Header.Get("X-GitHub-Api-Version"); got != "2022-11-28" {
			t.Errorf("X-GitHub-Api-Version: got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), token: "test-token", baseURL: server.URL}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestClient_do_noAuthWithoutToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization should be empty, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{httpClient: server.Client(), token: "", baseURL: server.URL}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = resp.Body.Close()
}

func TestClient_SearchRepositories(t *testing.T) {
	expected := SearchResult{
		TotalCount: 2,
		Items: []Repository{
			{FullName: "owner/repo1", StargazersCount: 500, Language: "Go"},
			{FullName: "owner/repo2", StargazersCount: 300, Language: "Rust"},
		},
	}

	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/repositories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("q") == "" {
			t.Error("query should not be empty")
		}
		if q.Get("sort") != "stars" {
			t.Errorf("sort: got %q, want stars", q.Get("sort"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("order: got %q, want desc", q.Get("order"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("per_page: got %q, want 10", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	result, err := client.SearchRepositories("test query", "stars", "desc", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalCount != 2 {
		t.Errorf("TotalCount: got %d, want 2", result.TotalCount)
	}
	if len(result.Items) != 2 {
		t.Fatalf("Items: got %d, want 2", len(result.Items))
	}
	if result.Items[0].FullName != "owner/repo1" {
		t.Errorf("Items[0].FullName: got %q", result.Items[0].FullName)
	}
}

func TestClient_SearchRepositories_NoSortOrder(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sort") != "" {
			t.Errorf("sort should be empty, got %q", q.Get("sort"))
		}
		if q.Get("order") != "" {
			t.Errorf("order should be empty, got %q", q.Get("order"))
		}
		_ = json.NewEncoder(w).Encode(SearchResult{})
	})
	defer server.Close()

	_, err := client.SearchRepositories("test", "", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_SearchRepositories_APIError(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"rate limit exceeded"}`))
	})
	defer server.Close()

	_, err := client.SearchRepositories("test", "", "", 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_GetRepository(t *testing.T) {
	expected := Repository{
		FullName:        "golang/go",
		Description:     "The Go programming language",
		StargazersCount: 130000,
		Language:        "Go",
		Topics:          []string{"go", "programming-language"},
	}

	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/golang/go" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	repo, err := client.GetRepository("golang", "go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.FullName != "golang/go" {
		t.Errorf("FullName: got %q", repo.FullName)
	}
	if repo.StargazersCount != 130000 {
		t.Errorf("Stars: got %d", repo.StargazersCount)
	}
}

func TestClient_GetRepository_NotFound(t *testing.T) {
	server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	})
	defer server.Close()

	_, err := client.GetRepository("no", "exist")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_GetTrending(t *testing.T) {
	cases := []struct {
		since    string
		language string
	}{
		{"daily", ""},
		{"weekly", "go"},
		{"monthly", "rust"},
	}

	for _, tc := range cases {
		t.Run(tc.since+"_"+tc.language, func(t *testing.T) {
			server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query().Get("q")
				if q == "" {
					t.Error("query should not be empty")
				}
				_ = json.NewEncoder(w).Encode(SearchResult{TotalCount: 1, Items: []Repository{{FullName: "test/repo"}}})
			})
			defer server.Close()

			result, err := client.GetTrending(tc.language, tc.since, 5)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.TotalCount != 1 {
				t.Errorf("TotalCount: got %d", result.TotalCount)
			}
		})
	}
}

func TestClient_GetTrendingByStars(t *testing.T) {
	cases := []struct {
		since    string
		language string
	}{
		{"daily", ""},
		{"weekly", "python"},
		{"monthly", ""},
	}

	for _, tc := range cases {
		t.Run(tc.since+"_"+tc.language, func(t *testing.T) {
			server, client := newTestServer(func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(SearchResult{TotalCount: 5, Items: []Repository{{FullName: "trending/repo"}}})
			})
			defer server.Close()

			result, err := client.GetTrendingByStars(tc.language, tc.since, 10)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.TotalCount != 5 {
				t.Errorf("TotalCount: got %d", result.TotalCount)
			}
		})
	}
}
