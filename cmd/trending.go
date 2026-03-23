package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/atani/github-discover/internal/cache"
	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/ui"
	"github.com/spf13/cobra"
)

var (
	trendCount    int
	trendLanguage string
	trendSince    string
)

var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Show trending repositories",
	Long:  `Display trending GitHub repositories sorted by star count.`,
	RunE:  runTrending,
}

func init() {
	rootCmd.AddCommand(trendingCmd)
	trendingCmd.Flags().IntVarP(&trendCount, "number", "n", 20, "Number of repositories to show")
	trendingCmd.Flags().StringVarP(&trendLanguage, "language", "l", "", "Filter by programming language")
	trendingCmd.Flags().StringVar(&trendSince, "since", "weekly", "Time range: daily, weekly, monthly")
}

func runTrending(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	client := github.NewClient(os.Getenv("GITHUB_TOKEN"))
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	cacheKey := fmt.Sprintf("%strending_%s_%s", cache.TrendPrefix, trendSince, trendLanguage)

	var result *github.SearchResult

	if data, ok := c.Get(cacheKey, cache.SearchTTL); ok && !refresh {
		if err := json.Unmarshal(data, &result); err == nil {
			goto render
		}
	}

	result, err = client.GetTrendingByStars(trendLanguage, trendSince, trendCount)
	if err != nil {
		return fmt.Errorf("failed to get trending repositories: %w", err)
	}

	if data, err := json.Marshal(result); err == nil {
		_ = c.Set(cacheKey, data)
	}

render:
	rows := buildRows(result.Items, trendCount)

	sinceLabel := trendSince
	if trendLanguage != "" {
		sinceLabel += ", " + trendLanguage
	}
	title := i18n.T("trending.title", sinceLabel)

	ui.PrintRepoTable(title, rows, i18n.T("trending.tip"))
	return nil
}

func buildRows(repos []github.Repository, limit int) []ui.RepoRow {
	if len(repos) > limit {
		repos = repos[:limit]
	}

	rows := make([]ui.RepoRow, len(repos))
	for i, repo := range repos {
		rows[i] = ui.RepoRow{
			Rank:        i + 1,
			Name:        repo.FullName,
			Stars:       repo.StargazersCount,
			Language:    repo.Language,
			Description: repo.Description,
		}
	}
	return rows
}
