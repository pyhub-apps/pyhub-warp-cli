package main

import (
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/cmd"
)

// Build variables - will be set during build
var (
	version   = "0.1.0-dev"
	gitCommit = "unknown"
	buildDate = "unknown"
)

func main() {
	// Set version information
	cmd.SetVersionInfo(version, gitCommit, buildDate)
	
	// Execute the root command
	cmd.Execute()
}