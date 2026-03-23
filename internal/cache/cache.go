package cache

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SearchTTL = 1 * time.Hour
	RepoTTL   = 24 * time.Hour

	SearchPrefix = "search_"
	TrendPrefix  = "trend_"
)

type Cache struct {
	dir string
}

func New() (*Cache, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(cacheDir, "github-discover")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	return &Cache{dir: dir}, nil
}

func sanitizeKey(key string) string {
	r := strings.NewReplacer("/", "_", "\\", "_", " ", "_", ":", "_")
	return r.Replace(key)
}

func (c *Cache) Get(key string, ttl time.Duration) ([]byte, bool) {
	path := filepath.Join(c.dir, sanitizeKey(key)+".json")
	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}

	if time.Since(info.ModTime()) > ttl {
		return nil, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	return data, true
}

func (c *Cache) Set(key string, data []byte) error {
	path := filepath.Join(c.dir, sanitizeKey(key)+".json")
	return os.WriteFile(path, data, 0o600)
}
