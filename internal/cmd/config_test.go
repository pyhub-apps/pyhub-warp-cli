package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/config"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/testutil"
	"github.com/spf13/cobra"
)

func TestConfigCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantOutput  string
		checkOutput bool
	}{
		{
			name:        "No subcommand shows help",
			args:        []string{"config"},
			wantOutput:  "Sejong CLI의 설정을 관리합니다",
			checkOutput: true,
		},
		{
			name:        "Help flag",
			args:        []string{"config", "--help"},
			wantOutput:  "Sejong CLI의 설정을 관리합니다",
			checkOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			// Check output if needed
			if tt.checkOutput && !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("Output should contain %q, got %q", tt.wantOutput, buf.String())
			}
		})
	}
}

func TestConfigSetCommand(t *testing.T) {
	// Setup test config
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-cmd-test-*")
	defer cleanup()

	// Reset config and set test path
	config.ResetConfig()
	config.SetTestConfigPath(tempDir)

	// Initialize config
	if err := config.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:        "No arguments",
			args:        []string{"config", "set"},
			wantErr:     true,
			errContains: "accepts 2 arg(s), received 0",
		},
		{
			name:        "One argument",
			args:        []string{"config", "set", "law.key"},
			wantErr:     true,
			errContains: "accepts 2 arg(s), received 1",
		},
		{
			name:        "Invalid key",
			args:        []string{"config", "set", "invalid.key", "value"},
			wantErr:     true,
			errContains: "잘못된 설정 키 형식",
		},
		{
			name:        "Empty value",
			args:        []string{"config", "set", "law.key", ""},
			wantErr:     true,
			errContains: "설정값이 비어있습니다",
		},
		{
			name:       "Valid API key",
			args:       []string{"config", "set", "law.key", "test-api-key-123"},
			wantErr:    false,
			wantOutput: "API 키가 성공적으로 설정되었습니다",
		},
		{
			name:       "Valid API key with spaces",
			args:       []string{"config", "set", "law.key", "  test-key-with-spaces  "},
			wantErr:    false,
			wantOutput: "API 키가 성공적으로 설정되었습니다",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error should contain %q, got %q", tt.errContains, err.Error())
				}
			}

			if !tt.wantErr && tt.wantOutput != "" {
				if !strings.Contains(buf.String(), tt.wantOutput) {
					t.Errorf("Output should contain %q, got %q", tt.wantOutput, buf.String())
				}
			}
		})
	}
}

func TestConfigGetCommand(t *testing.T) {
	// Setup test config
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-cmd-test-*")
	defer cleanup()

	// Reset config and set test path
	config.ResetConfig()
	config.SetTestConfigPath(tempDir)

	// Initialize config
	if err := config.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	tests := []struct {
		name        string
		setup       func()
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:        "No arguments",
			args:        []string{"config", "get"},
			wantErr:     true,
			errContains: "accepts 1 arg(s), received 0",
		},
		{
			name:        "Invalid key",
			args:        []string{"config", "get", "invalid.key"},
			wantErr:     true,
			errContains: "잘못된 설정 키 형식",
		},
		{
			name:       "API key not set",
			args:       []string{"config", "get", "law.key"},
			wantErr:    false,
			wantOutput: "API 키 설정이 필요합니다",
		},
		{
			name: "API key set",
			setup: func() {
				config.SetAPIKey("test-api-key-12345")
			},
			args:       []string{"config", "get", "law.key"},
			wantErr:    false,
			wantOutput: "law.key: test-api-k",
		},
		{
			name: "Short API key",
			setup: func() {
				config.SetAPIKey("short")
			},
			args:       []string{"config", "get", "law.key"},
			wantErr:    false,
			wantOutput: "law.key: short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if provided
			if tt.setup != nil {
				tt.setup()
			}

			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(configCmd)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error should contain %q, got %q", tt.errContains, err.Error())
				}
			}

			if !tt.wantErr && tt.wantOutput != "" {
				if !strings.Contains(buf.String(), tt.wantOutput) {
					t.Errorf("Output should contain %q, got %q", tt.wantOutput, buf.String())
				}
			}

			// Reset for next test
			config.SetAPIKey("")
		})
	}
}

func TestConfigPathCommand(t *testing.T) {
	// Setup test config
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-cmd-test-*")
	defer cleanup()

	// Reset config and set test path
	config.ResetConfig()
	config.SetTestConfigPath(tempDir)

	// Initialize config
	if err := config.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Create a new root command for testing
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(configCmd)

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Set args
	cmd.SetArgs([]string{"config", "path"})

	// Execute command
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
		return
	}

	// Check output contains path
	if !strings.Contains(buf.String(), "설정 파일 경로:") {
		t.Errorf("Output should contain '설정 파일 경로:', got %q", buf.String())
	}

	// Check that a valid path is shown
	if !strings.Contains(buf.String(), filepath.Join("config.yaml")) && !strings.Contains(buf.String(), filepath.Join("config", "yaml")) {
		t.Errorf("Output should contain valid config path, got %q", buf.String())
	}
}

func TestIsValidConfigKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"law.key", true},
		{"law.key.extra", true}, // Nested under valid key
		{"invalid", false},
		{"invalid.key", false},
		{"law", false},
		{"", false},
		{"law.", false},
		{".key", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := isValidConfigKey(tt.key)
			if got != tt.valid {
				t.Errorf("isValidConfigKey(%q) = %v, want %v", tt.key, got, tt.valid)
			}
		})
	}
}
