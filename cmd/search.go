package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/atani/github-discover/internal/cache"
	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/ui"
	"github.com/spf13/cobra"
)

var (
	searchLanguage string
	searchLimit    int
	searchSort     string
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search repositories with enhanced results",
	Long:  `Search GitHub repositories by keyword, sorted by stars or relevance.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&searchLanguage, "language", "l", "", "Filter by programming language")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 20, "Maximum number of results")
	searchCmd.Flags().StringVar(&searchSort, "sort", "stars", "Sort by: stars, forks, updated, best-match")
}

func runSearch(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	query := args[0]
	if searchLanguage != "" {
		query += " language:" + searchLanguage
	}

	client := newGitHubClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	cacheKey := fmt.Sprintf("%s%s_%s_%s_%d", cache.SearchPrefix, args[0], searchLanguage, searchSort, searchLimit)

	var result *github.SearchResult

	if data, ok := c.Get(cacheKey, cache.SearchTTL); ok && !refresh {
		if err := json.Unmarshal(data, &result); err == nil {
			goto render
		}
	}

	{
		sort := searchSort
		if sort == "best-match" {
			sort = ""
		}
		result, err = client.SearchRepositories(query, sort, "desc", searchLimit)
		if err != nil {
			return fmt.Errorf("failed to search repositories: %w", err)
		}
	}

	if data, err := json.Marshal(result); err == nil {
		_ = c.Set(cacheKey, data)
	}

render:
	rows := buildRows(result.Items, searchLimit)
	title := i18n.T("search.results", args[0])
	tip := i18n.T("search.count", result.TotalCount)

	ui.PrintRepoTable(title, rows, tip)
	return nil
}
