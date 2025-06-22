package config

import (
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

	configContent := `database_path: /custom/path/db.sqlite
log_level: debug
cache_dir: /custom/cache
`

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

	if cfg.DatabasePath != "/custom/path/db.sqlite" {
		t.Errorf("Expected DatabasePath '/custom/path/db.sqlite', got '%s'", cfg.DatabasePath)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}

	if cfg.CacheDir != "/custom/cache" {
		t.Errorf("Expected CacheDir '/custom/cache', got '%s'", cfg.CacheDir)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	cfg := &Config{
		DatabasePath: "/test/path/db.sqlite",
		LogLevel:     "debug",
		CacheDir:     "/test/cache",
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	configPath := filepath.Join(tmpDir, ".gcp-iam", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedCfg.DatabasePath != cfg.DatabasePath {
		t.Errorf("Expected DatabasePath '%s', got '%s'", cfg.DatabasePath, loadedCfg.DatabasePath)
	}

	if loadedCfg.LogLevel != cfg.LogLevel {
		t.Errorf("Expected LogLevel '%s', got '%s'", cfg.LogLevel, loadedCfg.LogLevel)
	}

	if loadedCfg.CacheDir != cfg.CacheDir {
		t.Errorf("Expected CacheDir '%s', got '%s'", cfg.CacheDir, loadedCfg.CacheDir)
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