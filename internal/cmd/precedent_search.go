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

var precedentSearchCmd *cobra.Command

// initPrecedentSearchCmd initializes the precedent search subcommand
func initPrecedentSearchCmd() {
	precedentSearchCmd = &cobra.Command{
		Use:   "search <검색어>",
		Short: "판례 검색",
		Long:  "키워드로 판례를 검색합니다.",
		Example: `  # 판례 검색
  sejong precedent search "계약 해지"
  
  # JSON 형식으로 출력
  sejong precedent search "손해배상" --format json
  
  # 페이지 지정
  sejong precedent search "부당이득" --page 2 --size 20`,
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"s"},
		RunE:    runPrecedentSearchCommand,
	}

	// Flags
	precedentSearchCmd.Flags().StringVarP(&precOutputFormat, "format", "f", "table", "출력 형식 (table, json)")
	precedentSearchCmd.Flags().IntVarP(&precPageNo, "page", "p", 1, "페이지 번호")
	precedentSearchCmd.Flags().IntVarP(&precPageSize, "size", "s", 10, "페이지 크기")
}

// updatePrecedentSearchCommand updates precedent search command descriptions
func updatePrecedentSearchCommand() {
	if precedentSearchCmd != nil {
		precedentSearchCmd.Short = "판례 검색"
		precedentSearchCmd.Long = "키워드로 판례를 검색합니다."

		// Update flag descriptions
		if flag := precedentSearchCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json)"
		}
		if flag := precedentSearchCmd.Flags().Lookup("page"); flag != nil {
			flag.Usage = "페이지 번호"
		}
		if flag := precedentSearchCmd.Flags().Lookup("size"); flag != nil {
			flag.Usage = "페이지 크기"
		}
	}
}

func runPrecedentSearchCommand(cmd *cobra.Command, args []string) error {
	// Get output writer
	output := cmd.OutOrStdout()

	// Join arguments as search query
	query := strings.Join(args, " ")
	query = strings.TrimSpace(query)

	if query == "" {
		logger.Error("Search query is empty")
		return fmt.Errorf("검색어가 비어있습니다")
	}

	logger.Info("판례 검색 중... (검색어: %s, 페이지: %d, 크기: %d)", query, precPageNo, precPageSize)

	// Create API client for precedent
	client, err := api.CreateClient(api.APITypePrec)
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Prepare search request
	req := &api.UnifiedSearchRequest{
		Query:    query,
		PageNo:   precPageNo,
		PageSize: precPageSize,
		Type:     "JSON",
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
		results.TotalCount, precPageNo, precPageSize)

	// Format and output results
	formatter := outputPkg.NewFormatter(precOutputFormat)
	formattedOutput, err := formatter.FormatSearchResultToString(results)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf("출력 실패")
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}
