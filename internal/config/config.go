package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// ConfigDirName is the name of the config directory
	ConfigDirName = ".sejong"
	// ConfigFileName is the name of the config file
	ConfigFileName = "config"
	// ConfigFileType is the type of the config file
	ConfigFileType = "yaml"
)

// Config holds the application configuration
type Config struct {
	Law struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"law"`
}

var (
	cfg        *Config
	configPath string
)

// Initialize sets up the configuration system
func Initialize() error {
	// Only set config path if it's not already set (allows for testing)
	if configPath == "" {
		// Get home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Set config path
		configPath = filepath.Join(homeDir, ConfigDirName)
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Configure Viper
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(configPath)

	// Set defaults
	viper.SetDefault("law.key", "")

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, create it
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(); err != nil {
				return fmt.Errorf("failed to create default config: %w", err)
			}
			// Read the newly created config
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read config after creation: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal config
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig() error {
	configFile := filepath.Join(configPath, ConfigFileName+"."+ConfigFileType)
	
	// Default configuration content
	defaultConfig := `# Sejong CLI Configuration
# 세종 CLI 설정 파일

# 국가법령정보센터 API 설정
law:
  # API 인증키
  # https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9 에서 발급
  key: ""
`

	// Write default config
	if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}

// Get returns a configuration value by key
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString returns a string configuration value by key
func GetString(key string) string {
	return viper.GetString(key)
}

// Set sets a configuration value
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

// Save writes the current configuration to file
func Save() error {
	return viper.WriteConfig()
}

// GetAPIKey returns the configured API key
func GetAPIKey() string {
	if cfg == nil {
		return ""
	}
	return cfg.Law.Key
}

// SetAPIKey sets the API key and saves the configuration
func SetAPIKey(key string) error {
	Set("law.key", key)
	if err := Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	// Update in-memory config
	if cfg != nil {
		cfg.Law.Key = key
	}
	return nil
}

// IsAPIKeySet checks if an API key is configured
func IsAPIKeySet() bool {
	key := GetAPIKey()
	return key != ""
}

// GetConfigPath returns the configuration file path
func GetConfigPath() string {
	return filepath.Join(configPath, ConfigFileName+"."+ConfigFileType)
}