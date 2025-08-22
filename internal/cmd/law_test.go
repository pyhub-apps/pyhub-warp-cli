package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestLawCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "No arguments",
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg(s), received 0",
		},
		{
			name:        "Empty search query",
			args:        []string{""},
			wantErr:     true,
			errContains: "검색어를 입력해주세요",
		},
		{
			name:        "Multiple arguments",
			args:        []string{"arg1", "arg2"},
			wantErr:     true,
			errContains: "accepts 1 arg(s), received 2",
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
		})
	}
}

func TestLawCommandFlags(t *testing.T) {
	// Test that flags are properly registered
	if lawCmd.Flag("format") == nil {
		t.Error("format flag not registered")
	}
	if lawCmd.Flag("page") == nil {
		t.Error("page flag not registered")
	}
	if lawCmd.Flag("size") == nil {
		t.Error("size flag not registered")
	}

	// Test default values
	formatFlag := lawCmd.Flag("format")
	if formatFlag.DefValue != "table" {
		t.Errorf("format flag default = %s, want table", formatFlag.DefValue)
	}

	pageFlag := lawCmd.Flag("page")
	if pageFlag.DefValue != "1" {
		t.Errorf("page flag default = %s, want 1", pageFlag.DefValue)
	}

	sizeFlag := lawCmd.Flag("size")
	if sizeFlag.DefValue != "10" {
		t.Errorf("size flag default = %s, want 10", sizeFlag.DefValue)
	}
}