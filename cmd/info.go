package cmd

import (
	"fmt"
	"strings"

	"github.com/atani/github-discover/internal/github"
	"github.com/atani/github-discover/internal/i18n"
	"github.com/atani/github-discover/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <owner/repo>",
	Short: "Show detailed repository information",
	Long:  `Display detailed information about a GitHub repository.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	if lang != "" {
		i18n.SetLanguage(lang)
	}

	parts := strings.SplitN(args[0], "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("please specify repository as owner/repo (e.g. golang/go)")
	}

	client := newGitHubClient()

	repo, err := client.GetRepository(parts[0], parts[1])
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error.not_found", args[0]), err)
	}

	printRepoInfo(repo)
	return nil
}

func printRepoInfo(repo *github.Repository) {
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	divider := strings.Repeat("─", 60)

	fmt.Println()
	fmt.Println(nameStyle.Render(repo.FullName))
	fmt.Println()
	if repo.Description != "" {
		fmt.Println(repo.Description)
		fmt.Println()
	}
	fmt.Println(divider)

	// Popularity
	fmt.Println(sectionStyle.Render(i18n.T("info.popularity")))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.stars")+":"), valueStyle.Render(ui.FormatStars(repo.StargazersCount)))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.forks")+":"), valueStyle.Render(fmt.Sprintf("%d", repo.ForksCount)))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.open_issues")+":"), valueStyle.Render(fmt.Sprintf("%d", repo.OpenIssuesCount)))
	fmt.Println()

	// Details
	fmt.Println(sectionStyle.Render(i18n.T("info.details")))
	if repo.Language != "" {
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.language")+":"), valueStyle.Render(repo.Language))
	}
	if repo.License != nil {
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.license")+":"), valueStyle.Render(repo.License.Name))
	}
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.created")+":"), valueStyle.Render(repo.CreatedAt.Format("2006-01-02")))
	fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.updated")+":"), valueStyle.Render(repo.PushedAt.Format("2006-01-02")))

	if repo.Homepage != "" {
		fmt.Printf("   %s %s\n", labelStyle.Render(i18n.T("info.homepage")+":"), valueStyle.Render(repo.Homepage))
	}
	fmt.Println()

	// Topics
	if len(repo.Topics) > 0 {
		fmt.Println(sectionStyle.Render(i18n.T("info.topics")))
		fmt.Printf("   %s\n", valueStyle.Render(strings.Join(repo.Topics, ", ")))
		fmt.Println()
	}

	fmt.Printf("   URL: %s\n", repo.HTMLURL)
	fmt.Println(divider)
	fmt.Println()
}
