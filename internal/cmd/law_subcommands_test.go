package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pyhub-apps/sejong-cli/internal/i18n"
	"github.com/spf13/cobra"
)

func TestLawSubcommands(t *testing.T) {
	// Initialize i18n for testing (Korean by default)
	if err := i18n.Init(); err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}
	i18n.SetLanguage("ko")

	// Initialize law command
	initLawCmd()

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		// Search subcommand tests
		{
			name:        "Search without query",
			args:        []string{"search"},
			wantErr:     true,
			errContains: "accepts 1 arg",
		},
		{
			name:        "Search with empty query",
			args:        []string{"search", ""},
			wantErr:     true,
			errContains: "검색어를 입력해주세요",
		},
		{
			name:       "Search with valid query (no API key)",
			args:       []string{"search", "테스트"},
			wantErr:    false, // Shows API key setup guide, not an error
			wantOutput: "API 키 설정이 필요합니다",
		},
		// Detail subcommand tests
		{
			name:        "Detail without ID",
			args:        []string{"detail"},
			wantErr:     true,
			errContains: "accepts 1 arg",
		},
		{
			name:        "Detail with empty ID",
			args:        []string{"detail", ""},
			wantErr:     true,
			errContains: "법령ID가 비어있습니다",
		},
		{
			name:        "Detail with valid ID (no API key)",
			args:        []string{"detail", "001234"},
			wantErr:     true,
			errContains: "NLIC API 키가 설정되지 않았습니다",
		},
		// History subcommand tests
		{
			name:        "History without ID",
			args:        []string{"history"},
			wantErr:     true,
			errContains: "accepts 1 arg",
		},
		{
			name:        "History with empty ID",
			args:        []string{"history", ""},
			wantErr:     true,
			errContains: "법령ID가 비어있습니다",
		},
		{
			name:        "History with valid ID (no API key)",
			args:        []string{"history", "001234"},
			wantErr:     true,
			errContains: "NLIC API 키가 설정되지 않았습니다",
		},
		// Help tests
		{
			name:       "Law help",
			args:       []string{"--help"},
			wantErr:    false,
			wantOutput: "국가법령정보센터에서 법령 정보를 검색하고 상세 정보를 조회합니다",
		},
		{
			name:       "Search help",
			args:       []string{"search", "--help"},
			wantErr:    false,
			wantOutput: "국가법령정보센터에서 법령을 검색합니다",
		},
		{
			name:       "Detail help",
			args:       []string{"detail", "--help"},
			wantErr:    false,
			wantOutput: "법령ID로 상세 정보를 조회합니다",
		},
		{
			name:       "History help",
			args:       []string{"history", "--help"},
			wantErr:    false,
			wantOutput: "법령ID로 제정 및 개정 이력을 조회합니다",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(lawCmd)

			// Set args
			args := append([]string{"law"}, tt.args...)
			cmd.SetArgs(args)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

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

			// Check output if expected
			if tt.wantOutput != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantOutput) {
					t.Errorf("Output should contain %q, got %q", tt.wantOutput, output)
				}
			}
		})
	}
}

func TestLawBackwardCompatibility(t *testing.T) {
	// Initialize i18n for testing
	if err := i18n.Init(); err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}
	i18n.SetLanguage("ko")

	// Initialize law command
	initLawCmd()

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		description string
	}{
		{
			name:        "Direct search without subcommand",
			args:        []string{"개인정보"},
			wantErr:     false, // Should work for backward compatibility
			description: "sejong law '검색어' should still work",
		},
		{
			name:        "Direct search with flags",
			args:        []string{"개인정보", "--format", "json"},
			wantErr:     false,
			description: "sejong law '검색어' --format json should still work",
		},
		{
			name:        "Direct search with pagination",
			args:        []string{"민법", "--page", "2", "--size", "20"},
			wantErr:     false,
			description: "sejong law '검색어' with pagination should still work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(lawCmd)

			// Set args
			args := append([]string{"law"}, tt.args...)
			cmd.SetArgs(args)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Execute command
			err := cmd.Execute()

			// For backward compatibility tests, we expect to see API key error
			// (not a command structure error)
			if err != nil {
				// Should fail with API key error, not command structure error
				if !strings.Contains(buf.String(), "API 키") &&
					!strings.Contains(err.Error(), "API 키") {
					t.Errorf("%s: Expected API key error, got %v", tt.description, err)
				}
			}
		})
	}
}
