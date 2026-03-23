package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/atani/github-discover/internal/cache"
	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	randomCount    int
	randomLanguage string
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get random repository recommendations",
	Long:  `Discover new repositories through random recommendations from popular GitHub projects.`,
	RunE:  runRandom,
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().IntVarP(&randomCount, "number", "n", 1, "Number of random repositories to show")
	randomCmd.Flags().StringVarP(&randomLanguage, "language", "l", "", "Filter by programming language")
}

func runRandom(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	client := github.NewClient(os.Getenv("GITHUB_TOKEN"))
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Fetch a pool of popular repos to pick from
	cacheKey := fmt.Sprintf("%srandom_pool_%s", cache.SearchPrefix, randomLanguage)

	var result *github.SearchResult

	if data, ok := c.Get(cacheKey, cache.SearchTTL); ok && !refresh {
		if err := json.Unmarshal(data, &result); err == nil {
			goto pick
		}
	}

	{
		query := "stars:>100"
		if randomLanguage != "" {
			query += " language:" + randomLanguage
		}
		result, err = client.SearchRepositories(query, "stars", "desc", 100)
		if err != nil {
			return fmt.Errorf("failed to get repositories: %w", err)
		}
	}

	if data, err := json.Marshal(result); err == nil {
		_ = c.Set(cacheKey, data)
	}

pick:
	picked := pickRandomRepos(result.Items, randomCount)

	if randomCount == 1 && len(picked) == 1 {
		printRepoDetail(picked[0])
	} else {
		rows := make([]ui.RepoRow, len(picked))
		for i, repo := range picked {
			rows[i] = ui.RepoRow{
				Rank:        i + 1,
				Name:        repo.FullName,
				Stars:       repo.StargazersCount,
				Language:    repo.Language,
				Description: repo.Description,
			}
		}

		title := i18n.T("random.title.plural", len(picked))
		ui.PrintRepoTable(title, rows, i18n.T("random.tip"))
	}

	return nil
}

func pickRandomRepos(repos []github.Repository, count int) []github.Repository {
	if count >= len(repos) {
		return repos
	}

	shuffled := make([]github.Repository, len(repos))
	copy(shuffled, repos)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}

func printRepoDetail(repo github.Repository) {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2)

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	var content strings.Builder
	content.WriteString(nameStyle.Render(repo.FullName) + "\n")
	content.WriteString(strings.Repeat("─", 50) + "\n")

	if repo.Description != "" {
		content.WriteString(repo.Description + "\n\n")
	}

	content.WriteString(labelStyle.Render(i18n.T("info.stars")+": ") + ui.FormatStars(repo.StargazersCount) + "\n")
	content.WriteString(labelStyle.Render(i18n.T("info.forks")+": ") + fmt.Sprintf("%d", repo.ForksCount) + "\n")

	if repo.Language != "" {
		content.WriteString(labelStyle.Render(i18n.T("info.language")+": ") + repo.Language + "\n")
	}

	if repo.License != nil {
		content.WriteString(labelStyle.Render(i18n.T("info.license")+": ") + repo.License.Name + "\n")
	}

	content.WriteString(labelStyle.Render("URL: ") + repo.HTMLURL + "\n")

	fmt.Println()
	fmt.Println(boxStyle.Render(content.String()))
	fmt.Println()
}
