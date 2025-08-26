package main

import (
	"testing"
)

func TestBuildVariables(t *testing.T) {
	// Test that build variables have default values
	if version == "" {
		t.Error("version should have a default value")
	}
	if gitCommit == "" {
		t.Error("gitCommit should have a default value")
	}
	if buildDate == "" {
		t.Error("buildDate should have a default value")
	}

	// Check default values
	if version != "dev" {
		t.Errorf("version = %s, want dev", version)
	}
	if gitCommit != "unknown" {
		t.Errorf("gitCommit = %s, want unknown", gitCommit)
	}
	if buildDate != "unknown" {
		t.Errorf("buildDate = %s, want unknown", buildDate)
	}
}

// Note: We can't easily test main() as it calls cmd.Execute() which calls os.Exit
// In a production scenario, we would refactor main() to be more testable
// by having it return an error instead of calling os.Exit directly
