package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	SnippetsDir string
}

func GetSnippetDir() (string, error) {
	if envDir := os.Getenv("YANKR_SNIPPETS_DIR"); envDir != "" {
		return envDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "yankr", "snippets"), nil
}

func Load() (*Config, error) {
	snippetsDir, err := GetSnippetDir()
	if err != nil {
		return nil, err
	}

	config := Config{
		SnippetsDir: snippetsDir,
	}

	return &config, nil
}
