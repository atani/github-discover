package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atani/github-discover/internal/cache"
	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/ui"
	"github.com/spf13/cobra"
)

var similarLimit int

var similarCmd = &cobra.Command{
	Use:   "similar <owner/repo>",
	Short: "Find similar repositories",
	Long:  `Find repositories similar to a given repository based on topics, language, and description.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSimilar,
}

func init() {
	rootCmd.AddCommand(similarCmd)
	similarCmd.Flags().IntVarP(&similarLimit, "number", "n", 10, "Number of results to show")
}

func runSimilar(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	parts := strings.SplitN(args[0], "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("please specify repository as owner/repo (e.g. golang/go)")
	}

	client := newGitHubClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Fetch the source repository to get its topics and language
	repo, err := client.GetRepository(parts[0], parts[1])
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error.not_found", args[0]), err)
	}

	// Build a search query from the repo's topics and language
	query := buildSimilarQuery(repo)

	cacheKey := fmt.Sprintf("%ssimilar_%s_%s_%d", cache.SearchPrefix, parts[0], parts[1], similarLimit)

	var result *github.SearchResult

	if data, ok := c.Get(cacheKey, cache.SearchTTL); ok && !refresh {
		if err := json.Unmarshal(data, &result); err == nil {
			goto render
		}
	}

	result, err = client.SearchRepositories(query, "stars", "desc", similarLimit+5)
	if err != nil {
		return fmt.Errorf("failed to search similar repositories: %w", err)
	}

	if data, err := json.Marshal(result); err == nil {
		_ = c.Set(cacheKey, data)
	}

render:
	// Filter out the source repo itself
	filtered := make([]github.Repository, 0, len(result.Items))
	for _, r := range result.Items {
		if r.FullName != repo.FullName {
			filtered = append(filtered, r)
		}
	}

	rows := buildRows(filtered, similarLimit)
	title := fmt.Sprintf("Similar to %s", repo.FullName)
	ui.PrintRepoTable(title, rows, "")

	return nil
}

func buildSimilarQuery(repo *github.Repository) string {
	var parts []string

	// Use up to 3 topics for the query
	topicCount := 0
	for _, topic := range repo.Topics {
		if topicCount >= 3 {
			break
		}
		parts = append(parts, fmt.Sprintf("topic:%s", topic))
		topicCount++
	}

	// If no topics, fall back to keywords from description
	if len(parts) == 0 && repo.Description != "" {
		words := extractKeywords(repo.Description)
		if len(words) > 0 {
			parts = append(parts, strings.Join(words, " "))
		}
	}

	// Add language filter
	if repo.Language != "" {
		parts = append(parts, fmt.Sprintf("language:%s", repo.Language))
	}

	// Minimum star threshold
	parts = append(parts, "stars:>10")

	return strings.Join(parts, " ")
}

func extractKeywords(description string) []string {
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"in": true, "on": true, "at": true, "to": true, "for": true,
		"of": true, "with": true, "by": true, "from": true, "as": true,
		"it": true, "its": true, "that": true, "this": true, "not": true,
		"but": true, "if": true, "no": true, "do": true, "can": true,
		"has": true, "have": true, "had": true, "will": true, "would": true,
		"your": true, "you": true, "we": true, "they": true, "their": true,
	}

	words := strings.Fields(strings.ToLower(description))
	var keywords []string
	for _, w := range words {
		w = strings.Trim(w, ".,!?;:()[]{}\"'")
		if len(w) >= 3 && !stopWords[w] {
			keywords = append(keywords, w)
		}
		if len(keywords) >= 3 {
			break
		}
	}
	return keywords
}
