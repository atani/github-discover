package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atani/github-discover/internal/github"
)

func setupMockGitHub(t *testing.T) *httptest.Server {
	t.Helper()

	searchResult := github.SearchResult{
		TotalCount: 2,
		Items: []github.Repository{
			{
				FullName:        "test/alpha",
				Name:            "alpha",
				Owner:           github.Owner{Login: "test"},
				Description:     "Alpha project",
				HTMLURL:         "https://github.com/test/alpha",
				Language:        "Go",
				StargazersCount: 1000,
				ForksCount:      50,
				Topics:          []string{"cli", "go"},
			},
			{
				FullName:        "test/beta",
				Name:            "beta",
				Owner:           github.Owner{Login: "test"},
				Description:     "Beta project",
				HTMLURL:         "https://github.com/test/beta",
				Language:        "Rust",
				StargazersCount: 500,
				ForksCount:      20,
				Topics:          []string{"web"},
			},
		},
	}

	repoDetail := github.Repository{
		FullName:        "test/alpha",
		Name:            "alpha",
		Owner:           github.Owner{Login: "test"},
		Description:     "Alpha project for testing",
		HTMLURL:         "https://github.com/test/alpha",
		Language:        "Go",
		StargazersCount: 1000,
		ForksCount:      50,
		OpenIssuesCount: 5,
		Topics:          []string{"cli", "go", "terminal"},
		License:         &github.License{Name: "MIT License", Key: "mit"},
		Homepage:        "https://alpha.example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/search/repositories":
			_ = json.NewEncoder(w).Encode(searchResult)
		case "/repos/test/alpha":
			_ = json.NewEncoder(w).Encode(repoDetail)
		case "/repos/test/notfound":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"Not Found"}`))
		default:
			_ = json.NewEncoder(w).Encode(repoDetail)
		}
	}))

	return server
}

// withMockServer sets githubBaseURL, runs fn, then restores.
func withMockServer(t *testing.T, fn func()) {
	t.Helper()
	server := setupMockGitHub(t)
	defer server.Close()

	old := githubBaseURL
	githubBaseURL = server.URL
	defer func() { githubBaseURL = old }()

	// Use temp cache dir to avoid interference
	tmpDir, err := os.MkdirTemp("", "gh-discover-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Force refresh to avoid cached stale data
	oldRefresh := refresh
	refresh = true
	defer func() { refresh = oldRefresh }()

	fn()
}

func TestTrendingCommand_Execute(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"trending", "-n", "5"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("trending command failed: %v", err)
		}
	})
}

func TestSearchCommand_Execute(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"search", "test"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("search command failed: %v", err)
		}
	})
}

func TestRandomCommand_Execute(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"random", "-n", "2"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("random command failed: %v", err)
		}
	})
}

func TestRandomCommand_Execute_Single(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"random", "-n", "1"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("random single command failed: %v", err)
		}
	})
}

func TestBrowseCommand_Execute_Overview(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"browse"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("browse overview failed: %v", err)
		}
	})
}

func TestBrowseCommand_Execute_Category(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"browse", "cli"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("browse cli failed: %v", err)
		}
	})
}

func TestBrowseCommand_Execute_InvalidCategory(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"browse", "invalid-category"})
		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for invalid category")
		}
	})
}

func TestInfoCommand_Execute(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"info", "test/alpha"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("info command failed: %v", err)
		}
	})
}

func TestInfoCommand_Execute_InvalidFormat(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"info", "noslash"})
		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for invalid format")
		}
	})
}

func TestInfoCommand_Execute_NotFound(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"info", "test/notfound"})
		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for not found repo")
		}
	})
}

func TestSimilarCommand_Execute(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"similar", "test/alpha"})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("similar command failed: %v", err)
		}
	})
}

func TestSimilarCommand_Execute_InvalidFormat(t *testing.T) {
	withMockServer(t, func() {
		rootCmd.SetArgs([]string{"similar", "noslash"})
		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for invalid format")
		}
	})
}

func TestExecute_Version(t *testing.T) {
	SetVersion("test-version")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"--version"})
	_ = rootCmd.Execute()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("version output should not be empty")
	}
}
