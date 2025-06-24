package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabasePath string `yaml:"database_path" json:"database_path"`
	LogLevel     string `yaml:"log_level" json:"log_level"`
	CacheDir     string `yaml:"cache_dir" json:"cache_dir"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".gcp-iam")

	cfg := &Config{
		DatabasePath: filepath.Join(configDir, "database.sqlite"),
		LogLevel:     "info",
		CacheDir:     filepath.Join(configDir, "cache"),
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if err := loadFromFile(cfg, configPath); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.ensureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return cfg, nil
}

func loadFromFile(cfg *Config, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

func (cfg *Config) ensureDirectories() error {
	dirs := []string{
		filepath.Dir(cfg.DatabasePath),
		cfg.CacheDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".gcp-iam", "config.yaml"), nil
}
