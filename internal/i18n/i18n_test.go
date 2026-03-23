package i18n

import (
	"testing"
)

func TestSetLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"en", "en"},
		{"ja", "ja"},
		{"EN", "en"},
		{"JA", "ja"},
	}

	for _, tt := range tests {
		SetLanguage(tt.input)
		if currentLang != tt.expected {
			t.Errorf("SetLanguage(%q): got %q, want %q", tt.input, currentLang, tt.expected)
		}
	}
}

func TestSetLanguage_Invalid(t *testing.T) {
	SetLanguage("en")
	SetLanguage("unknown")
	if currentLang != "en" {
		t.Errorf("SetLanguage(unknown) should not change language, got %q", currentLang)
	}
}

func TestT_English(t *testing.T) {
	SetLanguage("en")

	tests := []struct {
		key      string
		expected string
	}{
		{"info.stars", "Stars"},
		{"info.forks", "Forks"},
		{"info.language", "Language"},
		{"category.cli", "CLI & Terminal Tools"},
		{"category.ai", "AI & Machine Learning"},
	}

	for _, tt := range tests {
		if got := T(tt.key); got != tt.expected {
			t.Errorf("T(%q): got %q, want %q", tt.key, got, tt.expected)
		}
	}
}

func TestT_Japanese(t *testing.T) {
	SetLanguage("ja")

	tests := []struct {
		key      string
		expected string
	}{
		{"info.stars", "スター数"},
		{"info.forks", "フォーク数"},
		{"category.cli", "CLI・ターミナルツール"},
		{"category.ai", "AI・機械学習"},
	}

	for _, tt := range tests {
		if got := T(tt.key); got != tt.expected {
			t.Errorf("T(%q): got %q, want %q", tt.key, got, tt.expected)
		}
	}
}

func TestT_WithArgs(t *testing.T) {
	SetLanguage("en")

	got := T("trending.title", "weekly")
	expected := "Trending Repositories (weekly)"

	if got != expected {
		t.Errorf("T with args: got %q, want %q", got, expected)
	}
}

func TestT_UnknownKey(t *testing.T) {
	SetLanguage("en")

	key := "unknown.key"
	if got := T(key); got != key {
		t.Errorf("T unknown key: got %q, want %q", got, key)
	}
}
