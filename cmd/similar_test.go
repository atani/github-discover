package cmd

import (
	"testing"
)

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		desc     string
		minCount int
		maxCount int
	}{
		{"The Go programming language", 1, 3},
		{"A simple web framework for building APIs", 1, 3},
		{"", 0, 0},
		{"a an the", 0, 0},
		{"kubernetes container orchestration platform for automating deployment", 1, 3},
	}

	for _, tt := range tests {
		got := extractKeywords(tt.desc)
		if len(got) < tt.minCount || len(got) > tt.maxCount {
			t.Errorf("extractKeywords(%q): got %d keywords %v, want between %d and %d",
				tt.desc, len(got), got, tt.minCount, tt.maxCount)
		}
	}
}

func TestExtractKeywords_NoStopWords(t *testing.T) {
	keywords := extractKeywords("The quick brown fox jumps over the lazy dog")

	stopWords := map[string]bool{"the": true, "a": true, "an": true, "is": true, "of": true}
	for _, kw := range keywords {
		if stopWords[kw] {
			t.Errorf("extractKeywords returned stop word: %q", kw)
		}
	}
}

func TestExtractKeywords_MaxThree(t *testing.T) {
	keywords := extractKeywords("kubernetes container orchestration platform for automating deployment scaling management")
	if len(keywords) > 3 {
		t.Errorf("extractKeywords returned %d keywords, want at most 3", len(keywords))
	}
}
