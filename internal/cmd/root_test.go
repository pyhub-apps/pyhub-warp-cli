package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantOutput  string
		wantErr     bool
		checkOutput bool
	}{
		{
			name:        "No arguments shows help",
			args:        []string{},
			wantOutput:  "Sejong CLI는 국가법령정보센터",
			checkOutput: true,
		},
		{
			name:        "Help flag",
			args:        []string{"--help"},
			wantOutput:  "Sejong CLI는 국가법령정보센터",
			checkOutput: true,
		},
		{
			name:        "Version flag",
			args:        []string{"--version"},
			wantOutput:  "version",
			checkOutput: true,
		},
		{
			name:    "Unknown subcommand",
			args:    []string{"unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := rootCmd

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

			// Check output if needed
			if tt.checkOutput && !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf("Output should contain %q, got %q", tt.wantOutput, buf.String())
			}
		})
	}
}

func TestRootCommandVerboseFlag(t *testing.T) {
	// Test that verbose flag is properly registered
	if rootCmd.PersistentFlags().Lookup("verbose") == nil {
		t.Error("verbose flag not registered")
	}

	// Test short flag
	if rootCmd.PersistentFlags().ShorthandLookup("v") == nil {
		t.Error("verbose short flag 'v' not registered")
	}

	// Test default value
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag.DefValue != "false" {
		t.Errorf("verbose flag default = %s, want false", verboseFlag.DefValue)
	}
}

func TestSetVersionInfo(t *testing.T) {
	// Save original values
	origVersion := Version
	origCommit := GitCommit
	origDate := BuildDate
	defer func() {
		Version = origVersion
		GitCommit = origCommit
		BuildDate = origDate
	}()

	// Test SetVersionInfo
	testVersion := "1.0.0"
	testCommit := "abc123"
	testDate := "2024-01-01"

	SetVersionInfo(testVersion, testCommit, testDate)

	if Version != testVersion {
		t.Errorf("Version = %s, want %s", Version, testVersion)
	}
	if GitCommit != testCommit {
		t.Errorf("GitCommit = %s, want %s", GitCommit, testCommit)
	}
	if BuildDate != testDate {
		t.Errorf("BuildDate = %s, want %s", BuildDate, testDate)
	}

	// Check that rootCmd.Version is updated
	expectedVersion := "1.0.0 (built 2024-01-01, commit abc123)"
	if rootCmd.Version != expectedVersion {
		t.Errorf("rootCmd.Version = %s, want %s", rootCmd.Version, expectedVersion)
	}
}

func TestInitConfig(t *testing.T) {
	// Test that initConfig doesn't panic
	// This is a basic smoke test
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("initConfig() panicked: %v", r)
		}
	}()

	initConfig()
}

func TestExecuteFunction(t *testing.T) {
	// This test is tricky because Execute() calls os.Exit on error
	// We'll test the success case

	// Save original rootCmd
	origCmd := rootCmd
	defer func() {
		rootCmd = origCmd
	}()

	// Create a test command that doesn't exit
	testCmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd = testCmd

	// Capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{})

	// This should not panic or exit
	Execute()
}
