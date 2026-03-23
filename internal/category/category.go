package category

import "strings"

type Category string

const (
	CategoryCLI       Category = "cli"
	CategoryWeb       Category = "web"
	CategoryAI        Category = "ai"
	CategoryDevOps    Category = "devops"
	CategorySecurity  Category = "security"
	CategoryData      Category = "data"
	CategoryMobile    Category = "mobile"
	CategoryOther     Category = "other"
)

var AllCategories = []Category{
	CategoryCLI,
	CategoryWeb,
	CategoryAI,
	CategoryDevOps,
	CategorySecurity,
	CategoryData,
	CategoryMobile,
	CategoryOther,
}

var categoryKeywords = map[Category][]string{
	CategoryCLI:      {"cli", "terminal", "command-line", "shell", "console", "tui"},
	CategoryWeb:      {"web", "frontend", "backend", "http", "api", "rest", "graphql", "react", "vue", "next", "framework"},
	CategoryAI:       {"ai", "ml", "machine-learning", "deep-learning", "llm", "gpt", "neural", "nlp", "transformer", "diffusion"},
	CategoryDevOps:   {"devops", "docker", "kubernetes", "k8s", "ci", "cd", "infrastructure", "terraform", "ansible", "monitoring"},
	CategorySecurity: {"security", "crypto", "encryption", "auth", "vulnerability", "pentest", "firewall"},
	CategoryData:     {"database", "data", "sql", "nosql", "analytics", "etl", "pipeline", "streaming"},
	CategoryMobile:   {"mobile", "ios", "android", "flutter", "react-native", "swift", "kotlin"},
}

// Classify determines the category of a repository based on its topics and description.
func Classify(topics []string, description string) Category {
	descLower := strings.ToLower(description)
	topicSet := make(map[string]bool)
	for _, t := range topics {
		topicSet[strings.ToLower(t)] = true
	}

	bestCategory := CategoryOther
	bestScore := 0

	for cat, keywords := range categoryKeywords {
		score := 0
		for _, kw := range keywords {
			if topicSet[kw] {
				score += 2
			}
			if strings.Contains(descLower, kw) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	return bestCategory
}

func GetSymbol(cat Category) string {
	switch cat {
	case CategoryCLI:
		return ">"
	case CategoryWeb:
		return "#"
	case CategoryAI:
		return "*"
	case CategoryDevOps:
		return "!"
	case CategorySecurity:
		return "~"
	case CategoryData:
		return "="
	case CategoryMobile:
		return "@"
	default:
		return "+"
	}
}
