package main

import (
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/cmd"
)

// Build variables - will be set during build
// Following LINE HeadVer versioning: {head}.{yearweek}.{build}
var (
	version   = "1.2534.6-dev"
	gitCommit = "unknown"
	buildDate = "unknown"
)

func main() {
	// Set version information
	cmd.SetVersionInfo(version, gitCommit, buildDate)
	
	// Execute the root command
	cmd.Execute()
}