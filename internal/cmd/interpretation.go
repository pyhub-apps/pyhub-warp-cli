package cmd

import (
	"github.com/spf13/cobra"
)

var (
	interpretationCmd  *cobra.Command
	interpOutputFormat string
	interpPageNo       int
	interpPageSize     int
)

// initInterpretationCmd initializes the legal interpretation command and its subcommands
func initInterpretationCmd() {
	interpretationCmd = &cobra.Command{
		Use:   "interpretation",
		Short: "법령해석례 정보 검색 및 조회",
		Long: `정부 부처의 법령해석례를 검색하고 상세 정보를 조회합니다.
법령해석례는 법령의 적용과 해석에 대한 정부의 공식 견해입니다.

예시:
  sejong interpretation search "근로시간"  # 법령해석례 검색
  sejong interpretation detail 12345       # 법령해석례 상세 조회`,
		Aliases: []string{"interp", "expc"},
	}

	// Initialize subcommands
	initInterpretationSearchCmd()
	initInterpretationDetailCmd()

	// Add subcommands
	interpretationCmd.AddCommand(interpretationSearchCmd)
	interpretationCmd.AddCommand(interpretationDetailCmd)
}

// updateInterpretationCommand updates legal interpretation command descriptions for i18n
func updateInterpretationCommand() {
	if interpretationCmd != nil {
		interpretationCmd.Short = "법령해석례 정보 검색 및 조회"
		interpretationCmd.Long = `정부 부처의 법령해석례를 검색하고 상세 정보를 조회합니다.
법령해석례는 법령의 적용과 해석에 대한 정부의 공식 견해입니다.

예시:
  sejong interpretation search "근로시간"  # 법령해석례 검색
  sejong interpretation detail 12345       # 법령해석례 상세 조회`
	}

	// Update subcommands
	updateInterpretationSearchCommand()
	updateInterpretationDetailCommand()
}
