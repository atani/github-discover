package trending

import (
	"math"
	"testing"
	"time"

	"github.com/atani/github-discover/internal/github"
)

func TestCalcVelocity(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		repo      github.Repository
		expected  float64
		tolerance float64
	}{
		{
			name: "repo created 7 days ago with 700 stars",
			repo: github.Repository{
				StargazersCount: 700,
				CreatedAt:       now.AddDate(0, 0, -7),
			},
			expected:  100.0, // 700 / 7 = 100 stars/day
			tolerance: 1.0,
		},
		{
			name: "repo created 1 day ago with 50 stars",
			repo: github.Repository{
				StargazersCount: 50,
				CreatedAt:       now.AddDate(0, 0, -1),
			},
			expected:  50.0,
			tolerance: 1.0,
		},
		{
			name: "repo created 30 days ago with 300 stars",
			repo: github.Repository{
				StargazersCount: 300,
				CreatedAt:       now.AddDate(0, 0, -30),
			},
			expected:  10.0,
			tolerance: 0.5,
		},
		{
			name: "repo created today (zero days) treated as 1 day",
			repo: github.Repository{
				StargazersCount: 100,
				CreatedAt:       now,
			},
			expected:  100.0,
			tolerance: 10.0,
		},
		{
			name: "repo with zero stars",
			repo: github.Repository{
				StargazersCount: 0,
				CreatedAt:       now.AddDate(0, 0, -5),
			},
			expected:  0.0,
			tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcVelocity(tt.repo)
			if math.Abs(got-tt.expected) > tt.tolerance {
				t.Errorf("CalcVelocity(): got %.2f, want %.2f (tolerance %.2f)", got, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestDetectAnomalies(t *testing.T) {
	now := time.Now()

	// Create a set of repos with known velocities
	repos := []github.Repository{
		// Normal repos: ~10 stars/day
		{FullName: "normal/one", StargazersCount: 70, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "normal/two", StargazersCount: 80, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "normal/three", StargazersCount: 60, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "normal/four", StargazersCount: 75, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "normal/five", StargazersCount: 65, CreatedAt: now.AddDate(0, 0, -7)},
		// Anomaly: ~200 stars/day (far above the rest)
		{FullName: "hot/repo", StargazersCount: 1400, CreatedAt: now.AddDate(0, 0, -7)},
	}

	results := DetectAnomalies(repos)

	if len(results) == 0 {
		t.Fatal("DetectAnomalies: expected at least one anomaly")
	}

	// The hot repo should be detected
	found := false
	for _, r := range results {
		if r.Repo.FullName == "hot/repo" {
			found = true
			if r.Velocity < 100 {
				t.Errorf("hot/repo velocity: got %.2f, want >100", r.Velocity)
			}
			if r.ZScore < 1.0 {
				t.Errorf("hot/repo z-score: got %.2f, want >1.0", r.ZScore)
			}
		}
	}
	if !found {
		t.Error("DetectAnomalies: hot/repo not found in results")
	}
}

func TestDetectAnomalies_AllSimilar(t *testing.T) {
	now := time.Now()

	// All repos have similar velocity, so no anomalies expected
	repos := []github.Repository{
		{FullName: "a/one", StargazersCount: 70, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "a/two", StargazersCount: 72, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "a/three", StargazersCount: 68, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "a/four", StargazersCount: 71, CreatedAt: now.AddDate(0, 0, -7)},
	}

	results := DetectAnomalies(repos)
	if len(results) != 0 {
		t.Errorf("DetectAnomalies: expected no anomalies for similar repos, got %d", len(results))
	}
}

func TestDetectAnomalies_Empty(t *testing.T) {
	results := DetectAnomalies(nil)
	if len(results) != 0 {
		t.Errorf("DetectAnomalies(nil): got %d, want 0", len(results))
	}
}

func TestDetectAnomalies_SingleRepo(t *testing.T) {
	now := time.Now()
	repos := []github.Repository{
		{FullName: "solo/repo", StargazersCount: 500, CreatedAt: now.AddDate(0, 0, -3)},
	}

	results := DetectAnomalies(repos)
	// A single repo cannot be anomalous relative to nothing
	if len(results) != 0 {
		t.Errorf("DetectAnomalies single repo: got %d, want 0", len(results))
	}
}

func TestRankByVelocity(t *testing.T) {
	now := time.Now()

	repos := []github.Repository{
		{FullName: "slow/repo", StargazersCount: 10, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "fast/repo", StargazersCount: 700, CreatedAt: now.AddDate(0, 0, -7)},
		{FullName: "mid/repo", StargazersCount: 140, CreatedAt: now.AddDate(0, 0, -7)},
	}

	ranked := RankByVelocity(repos)

	if len(ranked) != 3 {
		t.Fatalf("RankByVelocity: got %d items, want 3", len(ranked))
	}

	// Should be sorted descending by velocity
	if ranked[0].Repo.FullName != "fast/repo" {
		t.Errorf("ranked[0]: got %q, want fast/repo", ranked[0].Repo.FullName)
	}
	if ranked[1].Repo.FullName != "mid/repo" {
		t.Errorf("ranked[1]: got %q, want mid/repo", ranked[1].Repo.FullName)
	}
	if ranked[2].Repo.FullName != "slow/repo" {
		t.Errorf("ranked[2]: got %q, want slow/repo", ranked[2].Repo.FullName)
	}

	// Velocity values should be correct
	if math.Abs(ranked[0].Velocity-100.0) > 1.0 {
		t.Errorf("ranked[0].Velocity: got %.2f, want ~100.0", ranked[0].Velocity)
	}
}

func TestRankByVelocity_Empty(t *testing.T) {
	ranked := RankByVelocity(nil)
	if len(ranked) != 0 {
		t.Errorf("RankByVelocity(nil): got %d, want 0", len(ranked))
	}
}

func TestScoredRepoFields(t *testing.T) {
	now := time.Now()
	repo := github.Repository{
		FullName:        "test/fields",
		StargazersCount: 350,
		CreatedAt:       now.AddDate(0, 0, -7),
	}

	ranked := RankByVelocity([]github.Repository{repo})
	if len(ranked) != 1 {
		t.Fatal("expected 1 result")
	}

	sr := ranked[0]
	if sr.Repo.FullName != "test/fields" {
		t.Errorf("FullName: got %q", sr.Repo.FullName)
	}
	if sr.Velocity < 49.0 || sr.Velocity > 51.0 {
		t.Errorf("Velocity: got %.2f, want ~50.0", sr.Velocity)
	}
}
