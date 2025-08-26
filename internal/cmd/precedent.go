package cmd

import (
	"github.com/spf13/cobra"
)

var (
	precedentCmd     *cobra.Command
	precOutputFormat string
	precPageNo       int
	precPageSize     int
)

// initPrecedentCmd initializes the precedent command and its subcommands
func initPrecedentCmd() {
	precedentCmd = &cobra.Command{
		Use:   "precedent",
		Short: "판례 정보 검색 및 조회",
		Long: `대법원 및 각급 법원의 판례를 검색하고 상세 정보를 조회합니다.

예시:
  warp precedent search "계약 해지"  # 판례 검색
  warp precedent detail 12345        # 판례 상세 조회`,
		Aliases: []string{"prec"},
	}

	// Initialize subcommands
	initPrecedentSearchCmd()
	initPrecedentDetailCmd()

	// Add subcommands
	precedentCmd.AddCommand(precedentSearchCmd)
	precedentCmd.AddCommand(precedentDetailCmd)
}

// updatePrecedentCommand updates precedent command descriptions for i18n
func updatePrecedentCommand() {
	if precedentCmd != nil {
		precedentCmd.Short = "판례 정보 검색 및 조회"
		precedentCmd.Long = `대법원 및 각급 법원의 판례를 검색하고 상세 정보를 조회합니다.

예시:
  warp precedent search "계약 해지"  # 판례 검색
  warp precedent detail 12345        # 판례 상세 조회`
	}

	// Update subcommands
	updatePrecedentSearchCommand()
	updatePrecedentDetailCommand()
}
