package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kborovik/gcp-iam/internal/constants"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabasePath string `yaml:"database_path" json:"database_path"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, constants.ConfigDirName)

	cfg := &Config{
		DatabasePath: filepath.Join(configDir, "database.sqlite"),
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
	dir := filepath.Dir(cfg.DatabasePath)
	if err := os.MkdirAll(dir, constants.DefaultDirPermissions); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, constants.ConfigDirName, "config.yaml"), nil
}
