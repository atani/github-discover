package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/atani/github-discover/internal/cache"
	"github.com/atani/github-discover/internal/category"
	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var browseLimit int

var browseCmd = &cobra.Command{
	Use:   "browse [category]",
	Short: "Browse repositories by category",
	Long: `Browse GitHub repositories organized by category.
Available categories: cli, web, ai, devops, security, data, mobile

Without a category argument, shows all categories.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBrowse,
}

func init() {
	rootCmd.AddCommand(browseCmd)
	browseCmd.Flags().IntVarP(&browseLimit, "limit", "n", 10, "Number of repositories per category")
}

func runBrowse(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	if len(args) == 0 {
		printCategoryOverview()
		return nil
	}

	cat := category.Category(args[0])
	valid := false
	for _, c := range category.AllCategories {
		if c == cat {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unknown category: %s", args[0])
	}

	return browseCategoryRepos(cat)
}

func printCategoryOverview() {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	catStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	usageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Println()
	fmt.Println(titleStyle.Render(i18n.T("browse.select")))
	fmt.Println()

	for _, cat := range category.AllCategories {
		emoji := category.GetSymbol(cat)
		name := i18n.T("category." + string(cat))
		fmt.Printf("  %s %s  %s\n",
			emoji,
			catStyle.Render(fmt.Sprintf("%-12s", cat)),
			descStyle.Render(name))
	}

	fmt.Println()
	fmt.Println(usageStyle.Render("Usage: github-discover browse <category>"))
	fmt.Println()
}

var categoryQueries = map[category.Category]string{
	category.CategoryCLI:      "topic:cli stars:>50",
	category.CategoryWeb:      "topic:web-framework stars:>50",
	category.CategoryAI:       "topic:machine-learning stars:>50",
	category.CategoryDevOps:   "topic:devops stars:>50",
	category.CategorySecurity: "topic:security stars:>50",
	category.CategoryData:     "topic:database stars:>50",
	category.CategoryMobile:   "topic:mobile stars:>50",
	category.CategoryOther:    "stars:>100",
}

func browseCategoryRepos(cat category.Category) error {
	client := newGitHubClient()
	c, err := cache.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	cacheKey := fmt.Sprintf("%sbrowse_%s_%d", cache.SearchPrefix, cat, browseLimit)

	var result *github.SearchResult

	if data, ok := c.Get(cacheKey, cache.SearchTTL); ok && !refresh {
		if err := json.Unmarshal(data, &result); err == nil {
			goto render
		}
	}

	{
		query := categoryQueries[cat]
		result, err = client.SearchRepositories(query, "stars", "desc", browseLimit)
		if err != nil {
			return fmt.Errorf("failed to browse category: %w", err)
		}
	}

	if data, err := json.Marshal(result); err == nil {
		_ = c.Set(cacheKey, data)
	}

render:
	rows := buildRows(result.Items, browseLimit)
	title := category.GetSymbol(cat) + " " + i18n.T("category."+string(cat))
	ui.PrintRepoTable(title, rows, "")

	return nil
}
