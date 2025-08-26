package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/api"
	cliErrors "github.com/pyhub-apps/pyhub-warp-cli/internal/errors"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/onboarding"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	searchOutputFormat string
	searchPageNo       int
	searchPageSize     int
	searchSource       string // "all", "law", "ordinance"
	searchRegion       string
	searchSort         string

	// testSearchClient allows injecting a mock client for testing
	testSearchClient api.ClientInterface
)

// searchCmd represents the unified search command
var searchCmd *cobra.Command

// initSearchCmd initializes the search command
func initSearchCmd() {
	searchCmd = &cobra.Command{
		Use:   "search <검색어>",
		Short: "법령 및 자치법규 통합 검색",
		Long: `법령(국가법령)과 자치법규를 동시에 검색합니다.

UnifiedClient를 사용하여 병렬로 검색을 수행하며,
결과를 통합하여 표시합니다.`,
		Example: `  # 모든 소스에서 검색
  warp search "개인정보"
  
  # 국가법령만 검색
  warp search "개인정보" --source law
  
  # 자치법규만 검색
  warp search "주차" --source ordinance
  
  # 서울 지역 포함 검색
  warp search "주차" --region 서울
  
  # JSON 형식으로 출력
  warp search "도로교통법" --format json`,
		Args: cobra.MinimumNArgs(1),
		RunE: runSearchCommand,
	}

	// Add flags
	searchCmd.Flags().StringVarP(&searchOutputFormat, "format", "f", "table", "출력 형식 (table, json, markdown, csv, html, html-simple)")
	searchCmd.Flags().IntVarP(&searchPageNo, "page", "p", 1, "페이지 번호")
	searchCmd.Flags().IntVarP(&searchPageSize, "size", "s", 50, "페이지 크기")
	searchCmd.Flags().StringVar(&searchSource, "source", "all", "검색 대상 (all, law, ordinance)")
	searchCmd.Flags().StringVarP(&searchRegion, "region", "r", "", "지역 필터 (자치법규용)")
	searchCmd.Flags().StringVar(&searchSort, "sort", "date", "정렬 순서 (date: 날짜순, name: 이름순)")
}

// updateSearchCommand updates search command descriptions
func updateSearchCommand() {
	if searchCmd != nil {
		searchCmd.Short = "법령 및 자치법규 통합 검색"
		searchCmd.Long = `법령(국가법령)과 자치법규를 동시에 검색합니다.

UnifiedClient를 사용하여 병렬로 검색을 수행하며,
결과를 통합하여 표시합니다.`

		// Update flag descriptions
		if flag := searchCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json, markdown, csv, html, html-simple)"
		}
		if flag := searchCmd.Flags().Lookup("page"); flag != nil {
			flag.Usage = "페이지 번호"
		}
		if flag := searchCmd.Flags().Lookup("size"); flag != nil {
			flag.Usage = "페이지 크기"
		}
		if flag := searchCmd.Flags().Lookup("source"); flag != nil {
			flag.Usage = "검색 대상 (all, law, ordinance)"
		}
		if flag := searchCmd.Flags().Lookup("region"); flag != nil {
			flag.Usage = "지역 필터 (자치법규용)"
		}
		if flag := searchCmd.Flags().Lookup("sort"); flag != nil {
			flag.Usage = "정렬 순서 (date: 날짜순, name: 이름순)"
		}
	}
}

// runSearchCommand handles the unified search command
func runSearchCommand(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(strings.Join(args, " "))
	if query == "" {
		logger.Debug("Empty query provided")
		return cliErrors.ErrEmptyQuery
	}

	logger.Debug("Starting unified search for query: %s", query)

	// Get verbose flag from root command
	verbose, _ := cmd.Root().Flags().GetBool("verbose")

	// Use test client if available (for testing)
	var client api.ClientInterface
	if testSearchClient != nil {
		client = testSearchClient
	} else {
		// Determine which client to create based on source flag
		var apiType api.APIType
		switch searchSource {
		case "law":
			apiType = api.APITypeNLIC
		case "ordinance":
			apiType = api.APITypeELIS
		case "all":
			apiType = api.APITypeAll
		default:
			return fmt.Errorf("잘못된 검색 대상: %s (all, law, ordinance 중 선택)", searchSource)
		}

		// Create appropriate API client
		apiClient, err := api.CreateClient(apiType)
		if err != nil {
			// Check if it's an API key error
			var cliErr *cliErrors.CLIError
			if errors.As(err, &cliErr) && cliErr.Code == cliErrors.ErrCodeNoAPIKey {
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil
			}

			// Also check for direct API key error message
			if strings.Contains(err.Error(), "API 키가 설정되지 않았습니다") {
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil
			}

			logger.LogError(err, verbose)
			return err
		}
		client = apiClient
	}

	// Log search parameters
	if searchSource == "all" {
		logger.Info("통합 검색 중... (검색어: %s, 페이지: %d, 크기: %d)", query, searchPageNo, searchPageSize)
	} else {
		logger.Info("검색 중... (검색어: %s, 대상: %s, 페이지: %d, 크기: %d)", query, searchSource, searchPageNo, searchPageSize)
	}
	
	if searchRegion != "" {
		logger.Info("지역 필터 적용: %s", searchRegion)
	}

	// Create search request
	req := &api.UnifiedSearchRequest{
		Query:    query,
		PageNo:   searchPageNo,
		PageSize: searchPageSize,
		Region:   searchRegion,
		Sort:     searchSort,
		Type:     "JSON", // Use JSON for unified search
	}

	// Search
	ctx := context.Background()
	response, err := client.Search(ctx, req)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// If API returns an API key error, show setup guide
			guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
			guide.ShowAPIKeySetup()
			return nil
		}

		logger.LogError(err, verbose)
		return fmt.Errorf("검색 실패: %w", err)
	}

	// Log completion
	logger.Info("검색 완료: %d개의 결과 (페이지: %d, 크기: %d)", response.TotalCount, searchPageNo, searchPageSize)

	// Output results
	return outputSearchResults(response, query, searchOutputFormat, searchPageNo, searchPageSize, cmd.OutOrStdout())
}

// outputSearchResults outputs search results in the specified format
func outputSearchResults(response *api.SearchResponse, query, format string, pageNo, pageSize int, writer io.Writer) error {
	if writer == nil {
		writer = os.Stdout
	}

	// Print summary
	fmt.Fprintf(writer, "총 %d개의 법령을 찾았습니다.\n\n", response.TotalCount)

	if response.TotalCount == 0 {
		fmt.Fprintln(writer, "검색 결과가 없습니다.")
		return nil
	}

	// Create formatter
	formatter := output.NewFormatter(format)

	// Format and output
	formattedOutput, err := formatter.FormatSearchResultToString(response)
	if err != nil {
		return fmt.Errorf("출력 형식 생성 실패: %w", err)
	}

	fmt.Fprint(writer, formattedOutput)

	// Print pagination info for table format
	if format == "table" && response.TotalCount > pageSize {
		totalPages := (response.TotalCount + pageSize - 1) / pageSize
		fmt.Fprintf(writer, "\n페이지 %d/%d (--page 옵션으로 다른 페이지 조회 가능)\n", pageNo, totalPages)
	}

	return nil
}

func init() {
	// Search command will be initialized and added in Execute()
}