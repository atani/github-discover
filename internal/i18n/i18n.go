package i18n

import (
	"fmt"
	"os"
	"strings"
)

var currentLang = "en"

var translations = map[string]map[string]string{
	"en": {
		"trending.title":      "Trending Repositories (%s)",
		"trending.tip":        "Use --language to filter by language, --since to change time range",
		"search.results":      "Search Results: \"%s\"",
		"search.count":        "%d repositories found",
		"random.title":        "Random Discovery",
		"random.title.plural": "Random Discoveries (%d)",
		"random.tip":          "Run again for different results!",
		"browse.select":       "Categories",
		"info.stars":          "Stars",
		"info.forks":          "Forks",
		"info.language":       "Language",
		"info.license":        "License",
		"info.created":        "Created",
		"info.updated":        "Last push",
		"info.homepage":       "Homepage",
		"info.topics":         "Topics",
		"info.popularity":     "Popularity",
		"info.details":        "Details",
		"info.open_issues":    "Open issues",
		"error.not_found":     "Repository not found: %s",
		"category.cli":        "CLI & Terminal Tools",
		"category.web":        "Web Development",
		"category.ai":         "AI & Machine Learning",
		"category.devops":     "DevOps & Infrastructure",
		"category.security":   "Security",
		"category.data":       "Data & Databases",
		"category.mobile":     "Mobile Development",
		"category.other":      "Other",
	},
	"ja": {
		"trending.title":      "トレンドリポジトリ (%s)",
		"trending.tip":        "--language で言語フィルタ、--since で期間変更",
		"search.results":      "検索結果: \"%s\"",
		"search.count":        "%d 件のリポジトリ",
		"random.title":        "ランダム発見",
		"random.title.plural": "ランダム発見 (%d件)",
		"random.tip":          "もう一度実行すると別の結果が出ます",
		"browse.select":       "カテゴリ",
		"info.stars":          "スター数",
		"info.forks":          "フォーク数",
		"info.language":       "言語",
		"info.license":        "ライセンス",
		"info.created":        "作成日",
		"info.updated":        "最終更新",
		"info.homepage":       "ホームページ",
		"info.topics":         "トピック",
		"info.popularity":     "人気度",
		"info.details":        "詳細",
		"info.open_issues":    "未解決Issue",
		"error.not_found":     "リポジトリが見つかりません: %s",
		"category.cli":        "CLI・ターミナルツール",
		"category.web":        "Web開発",
		"category.ai":         "AI・機械学習",
		"category.devops":     "DevOps・インフラ",
		"category.security":   "セキュリティ",
		"category.data":       "データ・DB",
		"category.mobile":     "モバイル開発",
		"category.other":      "その他",
	},
}

func SetLanguage(lang string) {
	lang = strings.ToLower(lang)
	if _, ok := translations[lang]; ok {
		currentLang = lang
	}
}

func init() {
	if envLang := os.Getenv("LANG"); envLang != "" {
		if strings.HasPrefix(envLang, "ja") {
			currentLang = "ja"
		}
	}
}

// T returns the translated string for the given key, with optional format args.
func T(key string, args ...any) string {
	dict := translations[currentLang]
	tmpl, ok := dict[key]
	if !ok {
		if fallback, ok := translations["en"][key]; ok {
			tmpl = fallback
		} else {
			return key
		}
	}

	if len(args) > 0 {
		return fmt.Sprintf(tmpl, args...)
	}
	return tmpl
}
