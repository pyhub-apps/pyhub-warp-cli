package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "버전 정보 표시",
	Long:  `Sejong CLI의 버전 정보와 빌드 세부사항을 표시합니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sejong CLI\n")
		fmt.Printf("  Version:   %s\n", Version)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
		fmt.Printf("  Built:     %s\n", BuildDate)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
