package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	AppName        = "yochat"
	ConfigDirEnv   = "XDG_CONFIG_HOME" // For Linux/Unix
	ConfigFileName = "config.json"
)

type Config struct {
	APIKey string `json:"api_key"`
}

// getConfigPath determines the cross-platform configuration directory.
// It follows XDG Base Directory Specification on Linux/Unix,
// uses ~/Library/Application Support on macOS, and %APPDATA% on Windows.
func GetConfigPath() (string, error) {
	var configPath string

	switch runtime.GOOS {
	case "windows":
		// On Windows, use %APPDATA%
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		configPath = filepath.Join(appData, AppName)
	case "darwin":
		// On macOS, use ~/Library/Application Support
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, "Library", "Application Support", AppName)
	case "linux", "freebsd", "netbsd", "openbsd":
		// On Linux/Unix, follow XDG Base Directory Specification
		// Prefer XDG_CONFIG_HOME, otherwise ~/.config
		xdgConfigHome := os.Getenv(ConfigDirEnv)
		if xdgConfigHome != "" {
			configPath = filepath.Join(xdgConfigHome, AppName)
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not get user home directory: %w", err)
			}
			configPath = filepath.Join(homeDir, ".config", AppName)
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return configPath, nil
}

// ensureConfigDir ensures the configuration directory exists.
func ensureConfigDir(path string) error {
	// 0700 means read, write, and execute permissions for the owner only.
	// This is a good default for sensitive configuration files.
	return os.MkdirAll(path, 0700)
}

// saveConfig saves the Config struct to the configuration file.
func SaveConfig(cfg *Config) error {
	configDirPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	if err := ensureConfigDir(configDirPath); err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	configFilePath := filepath.Join(configDirPath, ConfigFileName)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// 0600 means read and write permissions for the owner only.
	// This is important for security, especially for API keys.
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// loadConfig loads the Config struct from the configuration file.
func LoadConfig() (*Config, error) {
	configDirPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	configFilePath := filepath.Join(configDirPath, ConfigFileName)

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		// If the file doesn't exist, return an empty config and no error,
		// or handle it as a specific "not found" error if preferred.
		if os.IsNotExist(err) {
			return &Config{}, nil // Return an empty config if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from JSON: %w", err)
	}

	return &cfg, nil
}
