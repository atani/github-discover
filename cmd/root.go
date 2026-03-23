package cmd

import (
	"fmt"
	"os"

	"github.com/atani/github-discover/internal/github"
	"github.com/spf13/cobra"
)

var (
	lang    string
	refresh bool

	// githubBaseURL allows overriding the GitHub API base URL for testing.
	githubBaseURL string
)

var rootCmd = &cobra.Command{
	Use:   "github-discover",
	Short: "Discover trending and interesting GitHub repositories",
	Long: `github-discover helps you discover new and popular GitHub repositories.

Features:
  - View trending repositories
  - Browse repositories by category
  - Get random repository recommendations
  - Search with enhanced results
  - View detailed repository information`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func newGitHubClient() *github.Client {
	if githubBaseURL != "" {
		return github.NewTestClient(githubBaseURL)
	}
	return github.NewClient(os.Getenv("GITHUB_TOKEN"))
}

func init() {
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "", "Language (en, ja)")
	rootCmd.PersistentFlags().BoolVar(&refresh, "refresh", false, "Refresh cache")
}
