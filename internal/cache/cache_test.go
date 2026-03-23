package cache

import (
	"os"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	c := &Cache{dir: tmpDir}

	testData := []byte(`{"test":"data"}`)
	key := "test_key"

	err = c.Set(key, testData)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	result, ok := c.Get(key, time.Hour)
	if !ok {
		t.Fatal("Get returned false for existing cache")
	}

	if string(result) != string(testData) {
		t.Errorf("got %s, want %s", result, testData)
	}
}

func TestCache_Get_Expired(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	c := &Cache{dir: tmpDir}

	testData := []byte(`{"test":"data"}`)
	key := "test_key"

	err = c.Set(key, testData)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	_, ok := c.Get(key, 0)
	if ok {
		t.Fatal("Get returned true for expired cache")
	}
}

func TestCache_Get_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	c := &Cache{dir: tmpDir}

	_, ok := c.Get("nonexistent", time.Hour)
	if ok {
		t.Fatal("Get returned true for non-existent cache")
	}
}

func TestNew(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if c.dir == "" {
		t.Error("cache dir should not be empty")
	}
}
