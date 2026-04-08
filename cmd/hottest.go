package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/atani/github-discover/internal/cache"
	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/trending"
	"github.com/atani/github-discover/internal/ui"
	"github.com/spf13/cobra"
)

var (
	hottestCount    int
	hottestLanguage string
	hottestMode     string
)

var hottestCmd = &cobra.Command{
	Use:   "hottest",
	Short: "Detect repos with anomalous star growth",
	Long: `Find repositories with unusually high star acquisition velocity
over the last 7 days. Repos are scored by stars-per-day and those
significantly above the mean are flagged as anomalies.

Modes:
  anomaly  - Show only repos with statistically anomalous growth (default)
  velocity - Rank all recent repos by star velocity`,
	RunE: runHottest,
}

func init() {
	rootCmd.AddCommand(hottestCmd)
	hottestCmd.Flags().IntVarP(&hottestCount, "number", "n", 20, "Number of candidate repositories to analyze")
	hottestCmd.Flags().StringVarP(&hottestLanguage, "language", "l", "", "Filter by programming language")
	hottestCmd.Flags().StringVarP(&hottestMode, "mode", "m", "anomaly", "Detection mode: anomaly or velocity")
}

func runHottest(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	client := newGitHubClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	cacheKey := fmt.Sprintf("%shottest_%s_%s_%d", cache.TrendPrefix, hottestLanguage, stars, hottestCount)

	var result *github.SearchResult

	if data, ok := c.Get(cacheKey, cache.SearchTTL); ok && !refresh {
		if err := json.Unmarshal(data, &result); err == nil {
			goto analyze
		}
	}

	result, err = client.GetTrendingByStars(hottestLanguage, "weekly", hottestCount, starsQuery())
	if err != nil {
		return fmt.Errorf("failed to fetch repositories: %w", err)
	}

	if data, err := json.Marshal(result); err == nil {
		_ = c.Set(cacheKey, data)
	}

analyze:
	var scored []trending.ScoredRepo

	switch hottestMode {
	case "velocity":
		scored = trending.RankByVelocity(result.Items)
	default:
		scored = trending.DetectAnomalies(result.Items)
		if len(scored) == 0 {
			// Fall back to velocity ranking when no anomalies detected
			scored = trending.RankByVelocity(result.Items)
		}
	}

	rows := buildHottestRows(scored, hottestCount)
	title := i18n.T("hottest.title", hottestMode)
	tip := i18n.T("hottest.tip")
	ui.PrintRepoTable(title, rows, tip)
	return nil
}

func buildHottestRows(scored []trending.ScoredRepo, limit int) []ui.RepoRow {
	if len(scored) > limit {
		scored = scored[:limit]
	}

	rows := make([]ui.RepoRow, len(scored))
	for i, sr := range scored {
		desc := sr.Repo.Description
		velocityTag := fmt.Sprintf("[%.1f stars/day]", sr.Velocity)
		if desc != "" {
			desc = velocityTag + " " + desc
		} else {
			desc = velocityTag
		}

		rows[i] = ui.RepoRow{
			Rank:        i + 1,
			Name:        sr.Repo.FullName,
			Stars:       sr.Repo.StargazersCount,
			Language:    sr.Repo.Language,
			Description: desc,
		}
	}
	return rows
}
