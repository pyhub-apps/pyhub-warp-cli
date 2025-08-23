package testutil

import (
	"os"
	"testing"
)

// CreateTempDir creates a temporary directory for testing
// and returns the temp directory path and a cleanup function
func CreateTempDir(t *testing.T, pattern string) (string, func()) {
	t.Helper()

	// Create temporary directory for test config
	tempDir, err := os.MkdirTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}
