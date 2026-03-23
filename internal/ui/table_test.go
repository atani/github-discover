package ui

import (
	"bytes"
	"os"
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

func TestPrintRepoTable(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rows := []RepoRow{
		{Rank: 1, Name: "owner/repo1", Stars: 15000, Language: "Go", Description: "A great project"},
		{Rank: 2, Name: "owner/repo2", Stars: 500, Language: "Rust", Description: "Another project"},
		{Rank: 3, Name: "owner/repo3", Stars: 50, Language: "", Description: ""},
	}

	PrintRepoTable("Test Title", rows, "A helpful tip")

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Fatal("PrintRepoTable produced no output")
	}

	// Check that key content appears in output
	checks := []string{"Test Title", "owner/repo1", "owner/repo2", "owner/repo3", "15.0k", "A helpful tip", "A great project"}
	for _, check := range checks {
		if !bytesContains(output, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestPrintRepoTable_Empty(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintRepoTable("Empty Table", nil, "")

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !bytesContains(output, "Empty Table") {
		t.Error("output missing title for empty table")
	}
}

func TestPrintRepoTable_NoTip(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rows := []RepoRow{
		{Rank: 1, Name: "test/repo", Stars: 100, Language: "Go", Description: "Test"},
	}
	PrintRepoTable("Title", rows, "")

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Fatal("PrintRepoTable produced no output")
	}
}

func TestPrintRepoTable_NoLanguage(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rows := []RepoRow{
		{Rank: 1, Name: "test/nolang", Stars: 42, Language: "", Description: "No language set"},
	}
	PrintRepoTable("Title", rows, "")

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Should not panic or error
}

func bytesContains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
