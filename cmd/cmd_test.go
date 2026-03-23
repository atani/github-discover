package cmd

import (
	"testing"

	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/ui"
)

func TestBuildRows(t *testing.T) {
	repos := []github.Repository{
		{FullName: "a/one", StargazersCount: 1000, Language: "Go", Description: "First"},
		{FullName: "b/two", StargazersCount: 500, Language: "Rust", Description: "Second"},
		{FullName: "c/three", StargazersCount: 100, Language: "", Description: "Third"},
	}

	// Test normal case
	rows := buildRows(repos, 10)
	if len(rows) != 3 {
		t.Fatalf("buildRows: got %d rows, want 3", len(rows))
	}
	if rows[0].Rank != 1 || rows[0].Name != "a/one" || rows[0].Stars != 1000 {
		t.Errorf("rows[0]: got %+v", rows[0])
	}
	if rows[2].Language != "" {
		t.Errorf("rows[2].Language: got %q, want empty", rows[2].Language)
	}

	// Test limit
	rows = buildRows(repos, 2)
	if len(rows) != 2 {
		t.Fatalf("buildRows with limit 2: got %d rows", len(rows))
	}
	if rows[1].Name != "b/two" {
		t.Errorf("rows[1].Name: got %q", rows[1].Name)
	}
}

func TestBuildRows_Empty(t *testing.T) {
	rows := buildRows(nil, 10)
	if len(rows) != 0 {
		t.Errorf("buildRows(nil): got %d rows, want 0", len(rows))
	}
}

func TestPickRandomRepos(t *testing.T) {
	repos := make([]github.Repository, 50)
	for i := range repos {
		repos[i] = github.Repository{FullName: "test/repo"}
	}

	// Pick fewer than available
	picked := pickRandomRepos(repos, 5)
	if len(picked) != 5 {
		t.Errorf("pickRandomRepos(50, 5): got %d, want 5", len(picked))
	}

	// Pick more than available
	picked = pickRandomRepos(repos, 100)
	if len(picked) != 50 {
		t.Errorf("pickRandomRepos(50, 100): got %d, want 50", len(picked))
	}
}

func TestBuildSimilarQuery(t *testing.T) {
	tests := []struct {
		name     string
		repo     *github.Repository
		contains []string
	}{
		{
			name: "with topics and language",
			repo: &github.Repository{
				Topics:   []string{"cli", "terminal", "go"},
				Language: "Go",
			},
			contains: []string{"topic:cli", "topic:terminal", "topic:go", "language:Go", "stars:>10"},
		},
		{
			name: "no topics, uses description",
			repo: &github.Repository{
				Description: "A fast web framework",
				Language:    "Rust",
			},
			contains: []string{"fast", "language:Rust", "stars:>10"},
		},
		{
			name: "many topics, limited to 3",
			repo: &github.Repository{
				Topics:   []string{"a", "b", "c", "d", "e"},
				Language: "Python",
			},
			contains: []string{"topic:a", "topic:b", "topic:c", "language:Python"},
		},
		{
			name: "no topics, no description",
			repo: &github.Repository{
				Language: "Java",
			},
			contains: []string{"language:Java", "stars:>10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := buildSimilarQuery(tt.repo)
			for _, want := range tt.contains {
				if !containsString(query, want) {
					t.Errorf("buildSimilarQuery: query %q should contain %q", query, want)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSetVersion(t *testing.T) {
	SetVersion("1.2.3")
	if rootCmd.Version != "1.2.3" {
		t.Errorf("Version: got %q, want 1.2.3", rootCmd.Version)
	}
}

func TestPrintRepoDetail_NoFail(t *testing.T) {
	// Ensure printRepoDetail doesn't panic with various inputs
	repos := []github.Repository{
		{FullName: "test/full", Description: "A test", StargazersCount: 100, ForksCount: 10, Language: "Go", License: &github.License{Name: "MIT"}, HTMLURL: "https://github.com/test/full"},
		{FullName: "test/minimal", HTMLURL: "https://github.com/test/minimal"},
		{FullName: "test/nolicense", Description: "No license", StargazersCount: 50, Language: "Python", HTMLURL: "https://github.com/test/nolicense"},
	}

	for _, repo := range repos {
		printRepoDetail(repo) // Should not panic
	}
}

func TestPrintRepoTable_NoFail(t *testing.T) {
	rows := []ui.RepoRow{
		{Rank: 1, Name: "test/repo", Stars: 1000, Language: "Go", Description: "A test repo"},
		{Rank: 2, Name: "test/repo2", Stars: 500, Language: "", Description: ""},
	}
	ui.PrintRepoTable("Test Title", rows, "A tip")
	ui.PrintRepoTable("Empty", nil, "")
}
