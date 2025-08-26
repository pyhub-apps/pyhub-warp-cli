package cmd

import (
	"strings"

	cliErrors "github.com/pyhub-apps/pyhub-warp-cli/internal/errors"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/i18n"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
	"github.com/spf13/cobra"
)

// ordinanceSearchCmd represents the ordinance search subcommand
var ordinanceSearchCmd *cobra.Command

// initOrdinanceSearchCmd initializes the ordinance search subcommand
func initOrdinanceSearchCmd() {
	ordinanceSearchCmd = &cobra.Command{
		Use:   "search <검색어>",
		Short: i18n.T("ordinance.search.short"),
		Long:  i18n.T("ordinance.search.long"),
		Example: `  # 기본 검색
  warp ordinance search "주차 조례"
  
  # 지역별 검색
  warp ordinance search "환경" --region 서울
  warp ordinance search "도시계획" --region 부산
  
  # JSON 형식으로 출력
  warp ordinance search "건축 조례" --format json
  
  # 페이지네이션 옵션
  warp ordinance search "교통" --page 2 --size 20`,
		Args: cobra.ExactArgs(1),
		RunE: runOrdinanceSearchCommand,
	}

	// Flags are inherited from parent command
}

// updateOrdinanceSearchCommand updates ordinance search command descriptions
func updateOrdinanceSearchCommand() {
	if ordinanceSearchCmd != nil {
		ordinanceSearchCmd.Short = i18n.T("ordinance.search.short")
		ordinanceSearchCmd.Long = i18n.T("ordinance.search.long")
	}
}

func runOrdinanceSearchCommand(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(args[0])
	if query == "" {
		logger.Debug("Empty query provided")
		return cliErrors.ErrEmptyQuery
	}

	// Call the parent command's search function
	return runOrdinanceSearch(cmd, []string{query})
}
