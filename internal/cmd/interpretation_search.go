package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pyhub-apps/sejong-cli/internal/api"
	"github.com/pyhub-apps/sejong-cli/internal/logger"
	outputPkg "github.com/pyhub-apps/sejong-cli/internal/output"
	"github.com/spf13/cobra"
)

var interpretationSearchCmd *cobra.Command

// initInterpretationSearchCmd initializes the legal interpretation search subcommand
func initInterpretationSearchCmd() {
	interpretationSearchCmd = &cobra.Command{
		Use:   "search <검색어>",
		Short: "법령해석례 검색",
		Long:  "키워드로 법령해석례를 검색합니다.",
		Example: `  # 법령해석례 검색
  sejong interpretation search "근로시간"
  
  # JSON 형식으로 출력
  sejong interpretation search "휴가" --format json
  
  # 페이지 지정
  sejong interpretation search "임금" --page 2 --size 20`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"s"},
		RunE:    runInterpretationSearchCommand,
	}

	// Flags
	interpretationSearchCmd.Flags().StringVarP(&interpOutputFormat, "format", "f", "table", "출력 형식 (table, json)")
	interpretationSearchCmd.Flags().IntVarP(&interpPageNo, "page", "p", 1, "페이지 번호")
	interpretationSearchCmd.Flags().IntVarP(&interpPageSize, "size", "s", 50, "페이지 크기")
}

// updateInterpretationSearchCommand updates legal interpretation search command descriptions
func updateInterpretationSearchCommand() {
	if interpretationSearchCmd != nil {
		interpretationSearchCmd.Short = "법령해석례 검색"
		interpretationSearchCmd.Long = "키워드로 법령해석례를 검색합니다."

		// Update flag descriptions
		if flag := interpretationSearchCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json)"
		}
		if flag := interpretationSearchCmd.Flags().Lookup("page"); flag != nil {
			flag.Usage = "페이지 번호"
		}
		if flag := interpretationSearchCmd.Flags().Lookup("size"); flag != nil {
			flag.Usage = "페이지 크기"
		}
	}
}

func runInterpretationSearchCommand(cmd *cobra.Command, args []string) error {
	// Get output writer
	output := cmd.OutOrStdout()

	// Join arguments as search query
	query := strings.Join(args, " ")
	query = strings.TrimSpace(query)

	if query == "" {
		logger.Error("Search query is empty")
		return fmt.Errorf("검색어가 비어있습니다")
	}

	logger.Info("법령해석례 검색 중... (검색어: %s, 페이지: %d, 크기: %d)", query, interpPageNo, interpPageSize)

	// Create API client for legal interpretation
	client, err := api.CreateClient(api.APITypeExpc)
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Prepare search request
	req := &api.UnifiedSearchRequest{
		Query:    query,
		PageNo:   interpPageNo,
		PageSize: interpPageSize,
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
		results.TotalCount, interpPageNo, interpPageSize)

	// Format and output results
	formatter := outputPkg.NewFormatter(interpOutputFormat)
	formattedOutput, err := formatter.FormatSearchResultToString(results)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf("출력 실패")
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}
