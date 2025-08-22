package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/testutil"
	"github.com/spf13/viper"
)

func TestInitialize(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	// Test initialization
	err := Initialize()
	if err != nil {
		t.Errorf("Initialize() error = %v, want nil", err)
	}

	// Check if config file was created
	configFile := filepath.Join(tempDir, ConfigFileName+"."+ConfigFileType)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Check if cfg is initialized
	if cfg == nil {
		t.Error("Config struct was not initialized")
	}
}

func TestInitialize_ExistingConfig(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	// Create a config file with content
	configFile := filepath.Join(tempDir, ConfigFileName+"."+ConfigFileType)
	content := `law:
  key: "test-api-key"`
	
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Initialize with existing config
	err := Initialize()
	if err != nil {
		t.Errorf("Initialize() error = %v, want nil", err)
	}

	// Verify the API key was loaded
	if cfg.Law.Key != "test-api-key" {
		t.Errorf("API key = %q, want %q", cfg.Law.Key, "test-api-key")
	}
}

func TestInitialize_InvalidYAML(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	// Create an invalid YAML config file
	configFile := filepath.Join(tempDir, ConfigFileName+"."+ConfigFileType)
	content := `law:
  key: [invalid yaml`
	
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Initialize should fail with invalid YAML
	err := Initialize()
	if err == nil {
		t.Error("Initialize() should have returned an error for invalid YAML")
	}
}

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		expected string
	}{
		{
			name: "No config initialized",
			setup: func() {
				cfg = nil
			},
			expected: "",
		},
		{
			name: "Empty API key",
			setup: func() {
				cfg = &Config{}
			},
			expected: "",
		},
		{
			name: "Valid API key",
			setup: func() {
				cfg = &Config{}
				cfg.Law.Key = "test-key-123"
			},
			expected: "test-key-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := GetAPIKey()
			if got != tt.expected {
				t.Errorf("GetAPIKey() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSetAPIKey(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	// Initialize config first
	if err := Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Set API key
	testKey := "new-test-key-456"
	err := SetAPIKey(testKey)
	if err != nil {
		t.Errorf("SetAPIKey() error = %v, want nil", err)
	}

	// Verify in-memory config was updated
	if cfg.Law.Key != testKey {
		t.Errorf("In-memory API key = %q, want %q", cfg.Law.Key, testKey)
	}

	// Verify viper config was updated
	viperKey := viper.GetString("law.key")
	if viperKey != testKey {
		t.Errorf("Viper API key = %q, want %q", viperKey, testKey)
	}

	// Verify config was saved to file
	viper.Reset()
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(tempDir)
	
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}
	
	savedKey := viper.GetString("law.key")
	if savedKey != testKey {
		t.Errorf("Saved API key = %q, want %q", savedKey, testKey)
	}
}

func TestIsAPIKeySet(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		expected bool
	}{
		{
			name: "No config",
			setup: func() {
				cfg = nil
			},
			expected: false,
		},
		{
			name: "Empty key",
			setup: func() {
				cfg = &Config{}
				cfg.Law.Key = ""
			},
			expected: false,
		},
		{
			name: "Valid key",
			setup: func() {
				cfg = &Config{}
				cfg.Law.Key = "some-key"
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := IsAPIKeySet()
			if got != tt.expected {
				t.Errorf("IsAPIKeySet() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Setup viper with test values
	viper.Reset()
	viper.Set("test.key", "test-value")
	viper.Set("test.number", 42)
	viper.Set("test.bool", true)

	tests := []struct {
		name     string
		key      string
		expected interface{}
	}{
		{"String value", "test.key", "test-value"},
		{"Number value", "test.number", 42},
		{"Bool value", "test.bool", true},
		{"Non-existent key", "test.missing", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.key)
			if got != tt.expected {
				t.Errorf("Get(%q) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	// Setup viper with test values
	viper.Reset()
	viper.Set("test.string", "hello")
	viper.Set("test.number", 123)

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"String value", "test.string", "hello"},
		{"Number as string", "test.number", "123"},
		{"Non-existent key", "test.missing", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetString(tt.key)
			if got != tt.expected {
				t.Errorf("GetString(%q) = %q, want %q", tt.key, got, tt.expected)
			}
		})
	}
}

func TestSet(t *testing.T) {
	viper.Reset()

	// Test setting various types
	Set("test.string", "value")
	Set("test.int", 100)
	Set("test.bool", false)

	// Verify values were set
	if v := viper.GetString("test.string"); v != "value" {
		t.Errorf("Set string failed: got %q, want %q", v, "value")
	}
	if v := viper.GetInt("test.int"); v != 100 {
		t.Errorf("Set int failed: got %d, want %d", v, 100)
	}
	if v := viper.GetBool("test.bool"); v != false {
		t.Errorf("Set bool failed: got %v, want %v", v, false)
	}
}

func TestSave(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	// Initialize config
	if err := Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Modify config
	Set("test.value", "saved-value")

	// Save config
	err := Save()
	if err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// Read saved config
	viper.Reset()
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(tempDir)
	
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	// Verify saved value
	savedValue := viper.GetString("test.value")
	if savedValue != "saved-value" {
		t.Errorf("Saved value = %q, want %q", savedValue, "saved-value")
	}
}

func TestGetConfigPath(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	expected := filepath.Join(tempDir, ConfigFileName+"."+ConfigFileType)
	got := GetConfigPath()
	
	if got != expected {
		t.Errorf("GetConfigPath() = %q, want %q", got, expected)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-config-test-*")
	defer cleanup()
	
	// Reset config and set test path
	ResetConfig()
	SetTestConfigPath(tempDir)

	// Create default config
	err := createDefaultConfig()
	if err != nil {
		t.Errorf("createDefaultConfig() error = %v, want nil", err)
	}

	// Check if file exists
	configFile := filepath.Join(tempDir, ConfigFileName+"."+ConfigFileType)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Default config file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Check for expected content
	expectedStrings := []string{
		"Sejong CLI Configuration",
		"law:",
		"key:",
	}

	for _, expected := range expectedStrings {
		if !contains(string(content), expected) {
			t.Errorf("Config file should contain %q", expected)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}