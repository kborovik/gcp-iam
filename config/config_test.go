package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	if cfg.DatabasePath == "" {
		t.Error("DatabasePath should not be empty")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.LogLevel)
	}

	if cfg.CacheDir == "" {
		t.Error("CacheDir should not be empty")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".gcp-iam")
	configFile := filepath.Join(configDir, "config.yaml")

	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configContent := fmt.Sprintf(`database_path: %s/db.sqlite
log_level: debug
cache_dir: %s/cache
`, configDir, configDir)

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	expectedDBPath := filepath.Join(configDir, "db.sqlite")
	if cfg.DatabasePath != expectedDBPath {
		t.Errorf("Expected DatabasePath '%s', got '%s'", expectedDBPath, cfg.DatabasePath)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}

	expectedCacheDir := filepath.Join(configDir, "cache")
	if cfg.CacheDir != expectedCacheDir {
		t.Errorf("Expected CacheDir '%s', got '%s'", expectedCacheDir, cfg.CacheDir)
	}
}

func TestEnsureDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		DatabasePath: filepath.Join(tmpDir, "data", "db.sqlite"),
		CacheDir:     filepath.Join(tmpDir, "cache"),
	}

	err := cfg.ensureDirectories()
	if err != nil {
		t.Fatalf("Failed to ensure directories: %v", err)
	}

	dbDir := filepath.Dir(cfg.DatabasePath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		t.Errorf("Database directory was not created: %s", dbDir)
	}

	if _, err := os.Stat(cfg.CacheDir); os.IsNotExist(err) {
		t.Errorf("Cache directory was not created: %s", cfg.CacheDir)
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	path, err := GetDefaultConfigPath()
	if err != nil {
		t.Fatalf("Failed to get default config path: %v", err)
	}

	if path == "" {
		t.Error("Config path should not be empty")
	}

	if !filepath.IsAbs(path) {
		t.Error("Config path should be absolute")
	}
}
