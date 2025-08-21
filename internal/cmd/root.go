package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information - will be set during build
	Version   = "0.1.0-dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sejong",
	Short: "한국 법령 정보 검색 CLI 도구",
	Long: `Sejong CLI는 국가법령정보센터 오픈 API를 활용하여
터미널에서 빠르고 쉽게 한국 법령 정보를 검색할 수 있는 도구입니다.

개발자, 연구원, 법률 전문가들이 효율적으로 법령 정보에
접근할 수 있도록 설계되었습니다.`,
	Example: `  # 법령 검색
  sejong law "개인정보 보호법"
  
  # JSON 형식으로 출력
  sejong law "도로교통법" --format json
  
  # API 키 설정
  sejong config set law.key YOUR_API_KEY`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "상세 로그 출력")
	
	// Version flag
	rootCmd.Version = fmt.Sprintf("%s (built %s, commit %s)", Version, BuildDate, GitCommit)
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`)
}

// SetVersionInfo sets the version information for the CLI
func SetVersionInfo(version, commit, date string) {
	Version = version
	GitCommit = commit
	BuildDate = date
	rootCmd.Version = fmt.Sprintf("%s (built %s, commit %s)", version, date, commit)
}