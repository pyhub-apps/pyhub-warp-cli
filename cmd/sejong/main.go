package main

import (
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/cmd"
)

// Build variables - will be set during build
// Following HeadVer versioning: YYYY.0M.0D
var (
	version   = "2025.08.21-dev"
	gitCommit = "unknown"
	buildDate = "unknown"
)

func main() {
	// Set version information
	cmd.SetVersionInfo(version, gitCommit, buildDate)
	
	// Execute the root command
	cmd.Execute()
}