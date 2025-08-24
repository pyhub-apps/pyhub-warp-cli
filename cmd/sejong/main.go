package main

import (
	"github.com/pyhub-apps/sejong-cli/internal/cmd"
)

// Build variables - will be set during build
// Following LINE HeadVer versioning: {head}.{yearweek}.{build}
var (
	// NOTE: real value should be injected via -ldflags; keep a neutral default.
	version   = "dev"
	gitCommit = "unknown"
	buildDate = "unknown"
)

func main() {
	// Set version information
	cmd.SetVersionInfo(version, gitCommit, buildDate)

	// Execute the root command
	cmd.Execute()
}
