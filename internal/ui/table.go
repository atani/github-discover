package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type RepoRow struct {
	Rank        int
	Name        string
	Stars       int
	Language    string
	Description string
}

func FormatStars(count int) string {
	switch {
	case count >= 1000000:
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	case count >= 1000:
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	default:
		return fmt.Sprintf("%d", count)
	}
}

func PrintRepoTable(title string, rows []RepoRow, tip string) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	rankStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Width(4).Align(lipgloss.Right)
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	starStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	langStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	tipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)

	fmt.Println()
	fmt.Println(titleStyle.Render(title))
	fmt.Println(strings.Repeat("─", 80))

	for _, row := range rows {
		rank := rankStyle.Render(fmt.Sprintf("%d.", row.Rank))
		name := nameStyle.Render(fmt.Sprintf("%-40s", truncate(row.Name, 40)))
		stars := starStyle.Render(fmt.Sprintf("★ %-7s", FormatStars(row.Stars)))

		lang := ""
		if row.Language != "" {
			lang = langStyle.Render(fmt.Sprintf("[%s]", row.Language))
		}

		fmt.Printf("%s %s %s %s\n", rank, name, stars, lang)

		if row.Description != "" {
			desc := descStyle.Render("     " + truncate(row.Description, 73))
			fmt.Println(desc)
		}
	}

	fmt.Println(strings.Repeat("─", 80))
	if tip != "" {
		fmt.Println(tipStyle.Render(tip))
	}
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}
