package cmd

import (
	"github.com/spf13/cobra"
)

var (
	admruleCmd       *cobra.Command
	admrOutputFormat string
	admrPageNo       int
	admrPageSize     int
)

// initAdmruleCmd initializes the administrative rule command and its subcommands
func initAdmruleCmd() {
	admruleCmd = &cobra.Command{
		Use:   "admrule",
		Short: "행정규칙 정보 검색 및 조회",
		Long: `정부 부처 및 기관의 행정규칙(고시, 훈령, 예규 등)을 검색하고 상세 정보를 조회합니다.

예시:
  sejong admrule search "공공기관"  # 행정규칙 검색
  sejong admrule detail 12345       # 행정규칙 상세 조회`,
		Aliases: []string{"admr", "rule"},
	}

	// Initialize subcommands
	initAdmruleSearchCmd()
	initAdmruleDetailCmd()

	// Add subcommands
	admruleCmd.AddCommand(admruleSearchCmd)
	admruleCmd.AddCommand(admruleDetailCmd)
}

// updateAdmruleCommand updates administrative rule command descriptions for i18n
func updateAdmruleCommand() {
	if admruleCmd != nil {
		admruleCmd.Short = "행정규칙 정보 검색 및 조회"
		admruleCmd.Long = `정부 부처 및 기관의 행정규칙(고시, 훈령, 예규 등)을 검색하고 상세 정보를 조회합니다.

예시:
  sejong admrule search "공공기관"  # 행정규칙 검색
  sejong admrule detail 12345       # 행정규칙 상세 조회`
	}

	// Update subcommands
	updateAdmruleSearchCommand()
	updateAdmruleDetailCommand()
}
