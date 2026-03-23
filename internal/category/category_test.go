package category

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		name     string
		topics   []string
		desc     string
		expected Category
	}{
		{
			name:     "CLI tool by topic",
			topics:   []string{"cli", "terminal"},
			desc:     "A command-line tool",
			expected: CategoryCLI,
		},
		{
			name:     "Web framework by topic",
			topics:   []string{"web", "framework"},
			desc:     "A web framework",
			expected: CategoryWeb,
		},
		{
			name:     "AI project by topic",
			topics:   []string{"machine-learning", "deep-learning"},
			desc:     "Neural network library",
			expected: CategoryAI,
		},
		{
			name:     "DevOps by topic",
			topics:   []string{"kubernetes", "docker"},
			desc:     "Container orchestration",
			expected: CategoryDevOps,
		},
		{
			name:     "Security by description",
			topics:   []string{},
			desc:     "Encryption and security toolkit",
			expected: CategorySecurity,
		},
		{
			name:     "Database by topic",
			topics:   []string{"database", "sql"},
			desc:     "A database engine",
			expected: CategoryData,
		},
		{
			name:     "Mobile by topic",
			topics:   []string{"ios", "swift"},
			desc:     "Mobile development framework",
			expected: CategoryMobile,
		},
		{
			name:     "Other when no match",
			topics:   []string{},
			desc:     "Something completely unrelated to anything",
			expected: CategoryOther,
		},
		{
			name:     "Topics take priority over description",
			topics:   []string{"cli", "terminal", "tui"},
			desc:     "A web-based database security tool",
			expected: CategoryCLI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Classify(tt.topics, tt.desc)
			if got != tt.expected {
				t.Errorf("Classify(%v, %q): got %q, want %q", tt.topics, tt.desc, got, tt.expected)
			}
		})
	}
}

func TestGetEmoji(t *testing.T) {
	tests := []struct {
		cat      Category
		expected string
	}{
		{CategoryCLI, ">"},
		{CategoryWeb, "#"},
		{CategoryAI, "*"},
		{CategoryDevOps, "!"},
		{CategorySecurity, "~"},
		{CategoryData, "="},
		{CategoryMobile, "@"},
		{CategoryOther, "+"},
	}

	for _, tt := range tests {
		got := GetEmoji(tt.cat)
		if got != tt.expected {
			t.Errorf("GetEmoji(%q): got %q, want %q", tt.cat, got, tt.expected)
		}
	}
}

func TestAllCategories(t *testing.T) {
	expected := 8
	if len(AllCategories) != expected {
		t.Errorf("AllCategories length: got %d, want %d", len(AllCategories), expected)
	}
}
