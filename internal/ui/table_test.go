package ui

import (
	"testing"
)

func TestFormatStars(t *testing.T) {
	tests := []struct {
		count    int
		expected string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1.0k"},
		{1500, "1.5k"},
		{10000, "10.0k"},
		{100000, "100.0k"},
		{1000000, "1.0M"},
		{1500000, "1.5M"},
	}

	for _, tt := range tests {
		got := FormatStars(tt.count)
		if got != tt.expected {
			t.Errorf("FormatStars(%d): got %q, want %q", tt.count, got, tt.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is..."},
		{"", 5, ""},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.expected {
			t.Errorf("truncate(%q, %d): got %q, want %q", tt.input, tt.maxLen, got, tt.expected)
		}
	}
}
