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
		Key  string `mapstructure:"key"`  // Legacy: NLIC API key
		NLIC struct {
			Key string `mapstructure:"key"` // National Law Information Center API key
		} `mapstructure:"nlic"`
		ELIS struct {
			Key string `mapstructure:"key"` // Local Regulations Information System API key
		} `mapstructure:"elis"`
	} `mapstructure:"law"`
}

var (
	cfg        *Config
	configPath string
)

// SetTestConfigPath sets a custom config path for testing
// This should only be used in test files
func SetTestConfigPath(path string) {
	configPath = path
}

// ResetConfig resets the configuration state for testing
// This should only be used in test files
func ResetConfig() {
	cfg = nil
	configPath = ""
	viper.Reset()
}

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

	// Create config directory if it doesn't exist (restricted permissions for security)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Configure Viper
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(configPath)

	// Set defaults
	viper.SetDefault("law.key", "")
	viper.SetDefault("law.nlic.key", "")
	viper.SetDefault("law.elis.key", "")

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

# 법령 정보 API 설정
law:
  # 기본 API 키 (NLIC와 호환)
  key: ""
  
  # 국가법령정보센터 (National Law Information Center) API
  nlic:
    # API 인증키
    # https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9 에서 발급
    key: ""
  
  # 자치법규정보시스템 (Local Regulations Information System) API
  elis:
    # API 인증키
    # https://www.elis.go.kr 에서 발급
    key: ""
`

	// Write default config
	if err := os.WriteFile(configFile, []byte(defaultConfig), 0600); err != nil {
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

// GetAPIKey returns the configured API key (backward compatibility - returns NLIC key)
func GetAPIKey() string {
	if cfg == nil {
		return ""
	}
	// First check NLIC-specific key
	if cfg.Law.NLIC.Key != "" {
		return cfg.Law.NLIC.Key
	}
	// Fall back to legacy key
	return cfg.Law.Key
}

// SetAPIKey sets the API key and saves the configuration (backward compatibility - sets NLIC key)
func SetAPIKey(key string) error {
	// Set both legacy and NLIC keys for compatibility
	Set("law.key", key)
	Set("law.nlic.key", key)
	if err := Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	// Update in-memory config
	if cfg != nil {
		cfg.Law.Key = key
		cfg.Law.NLIC.Key = key
	}
	return nil
}

// IsAPIKeySet checks if an API key is configured (backward compatibility - checks NLIC key)
func IsAPIKeySet() bool {
	key := GetAPIKey()
	return key != ""
}

// GetNLICAPIKey returns the NLIC API key
func GetNLICAPIKey() string {
	if cfg == nil {
		return ""
	}
	// First check NLIC-specific key
	if cfg.Law.NLIC.Key != "" {
		return cfg.Law.NLIC.Key
	}
	// Fall back to legacy key
	return cfg.Law.Key
}

// SetNLICAPIKey sets the NLIC API key
func SetNLICAPIKey(key string) error {
	Set("law.nlic.key", key)
	// Also set legacy key for backward compatibility if it's empty
	if GetString("law.key") == "" {
		Set("law.key", key)
	}
	if err := Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	// Update in-memory config
	if cfg != nil {
		cfg.Law.NLIC.Key = key
		if cfg.Law.Key == "" {
			cfg.Law.Key = key
		}
	}
	return nil
}

// IsNLICAPIKeySet checks if NLIC API key is configured
func IsNLICAPIKeySet() bool {
	key := GetNLICAPIKey()
	return key != ""
}

// GetELISAPIKey returns the ELIS API key
func GetELISAPIKey() string {
	if cfg == nil {
		return ""
	}
	return cfg.Law.ELIS.Key
}

// SetELISAPIKey sets the ELIS API key
func SetELISAPIKey(key string) error {
	Set("law.elis.key", key)
	if err := Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	// Update in-memory config
	if cfg != nil {
		cfg.Law.ELIS.Key = key
	}
	return nil
}

// IsELISAPIKeySet checks if ELIS API key is configured
func IsELISAPIKeySet() bool {
	key := GetELISAPIKey()
	return key != ""
}

// GetConfigPath returns the configuration file path
func GetConfigPath() string {
	return filepath.Join(configPath, ConfigFileName+"."+ConfigFileType)
}
