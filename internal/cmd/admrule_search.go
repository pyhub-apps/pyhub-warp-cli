package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/api"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
	outputPkg "github.com/pyhub-apps/pyhub-warp-cli/internal/output"
	"github.com/spf13/cobra"
)

var admruleSearchCmd *cobra.Command

// initAdmruleSearchCmd initializes the administrative rule search subcommand
func initAdmruleSearchCmd() {
	admruleSearchCmd = &cobra.Command{
		Use:   "search <검색어>",
		Short: "행정규칙 검색",
		Long:  "키워드로 행정규칙(고시, 훈령, 예규 등)을 검색합니다.",
		Example: `  # 행정규칙 검색
  warp admrule search "공공기관"
  
  # JSON 형식으로 출력
  warp admrule search "개인정보" --format json
  
  # 페이지 지정
  warp admrule search "행정처분" --page 2 --size 20`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"s"},
		RunE:    runAdmruleSearchCommand,
	}

	// Flags
	admruleSearchCmd.Flags().StringVarP(&admrOutputFormat, "format", "f", "table", "출력 형식 (table, json)")
	admruleSearchCmd.Flags().IntVarP(&admrPageNo, "page", "p", 1, "페이지 번호")
	admruleSearchCmd.Flags().IntVarP(&admrPageSize, "size", "s", 50, "페이지 크기")
}

// updateAdmruleSearchCommand updates administrative rule search command descriptions
func updateAdmruleSearchCommand() {
	if admruleSearchCmd != nil {
		admruleSearchCmd.Short = "행정규칙 검색"
		admruleSearchCmd.Long = "키워드로 행정규칙(고시, 훈령, 예규 등)을 검색합니다."

		// Update flag descriptions
		if flag := admruleSearchCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json)"
		}
		if flag := admruleSearchCmd.Flags().Lookup("page"); flag != nil {
			flag.Usage = "페이지 번호"
		}
		if flag := admruleSearchCmd.Flags().Lookup("size"); flag != nil {
			flag.Usage = "페이지 크기"
		}
	}
}

func runAdmruleSearchCommand(cmd *cobra.Command, args []string) error {
	// Get output writer
	output := cmd.OutOrStdout()

	// Join arguments as search query
	query := strings.Join(args, " ")
	query = strings.TrimSpace(query)

	if query == "" {
		logger.Error("Search query is empty")
		return fmt.Errorf("검색어가 비어있습니다")
	}

	logger.Info("행정규칙 검색 중... (검색어: %s, 페이지: %d, 크기: %d)", query, admrPageNo, admrPageSize)

	// Create API client for administrative rule
	client, err := api.CreateClient(api.APITypeAdmrul)
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Prepare search request
	req := &api.UnifiedSearchRequest{
		Query:    query,
		PageNo:   admrPageNo,
		PageSize: admrPageSize,
		Type:     "XML",
	}

	// Search with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := client.Search(ctx, req)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(output, err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.Error("Search failed: %v", err)
		return fmt.Errorf("검색 실패: %v", err)
	}

	logger.Info("검색 완료: %d개의 결과 (페이지: %d, 크기: %d)",
		results.TotalCount, admrPageNo, admrPageSize)

	// Format and output results
	formatter := outputPkg.NewFormatter(admrOutputFormat)
	formattedOutput, err := formatter.FormatSearchResultToString(results)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf("출력 실패")
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}
