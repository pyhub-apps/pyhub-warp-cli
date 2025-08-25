package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/pyhub-apps/sejong-cli/internal/api"
	cliErrors "github.com/pyhub-apps/sejong-cli/internal/errors"
	"github.com/pyhub-apps/sejong-cli/internal/i18n"
	"github.com/pyhub-apps/sejong-cli/internal/logger"
	"github.com/pyhub-apps/sejong-cli/internal/onboarding"
	"github.com/pyhub-apps/sejong-cli/internal/output"
	"github.com/spf13/cobra"
)

// ordinanceCmd represents the ordinance command
var ordinanceCmd *cobra.Command

// Command flags
var (
	ordinanceOutputFormat string
	ordinancePageNo       int
	ordinancePageSize     int
	ordinanceRegion       string
	ordinanceSort         string
)

// Test helper - allows injection of mock client
var testOrdinanceClient api.ClientInterface

// initOrdinanceCmd initializes the ordinance command
func initOrdinanceCmd() {
	ordinanceCmd = &cobra.Command{
		Use:   "ordinance",
		Short: i18n.T("ordinance.short"),
		Long:  i18n.T("ordinance.long"),
		Example: `  # 조례 검색
  sejong ordinance search "주차 조례"
  sejong ordinance "환경 보호"  # search는 생략 가능
  
  # 지역별 조례 검색
  sejong ordinance search "도시계획" --region 서울
  
  # 조례 상세 조회
  sejong ordinance detail ORD123456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If args are provided without subcommand, run search
			if len(args) > 0 {
				return runOrdinanceSearch(cmd, args)
			}
			// Otherwise show help
			return cmd.Help()
		},
	}

	// Initialize subcommands
	initOrdinanceSearchCmd()
	initOrdinanceDetailCmd()

	// Add subcommands
	ordinanceCmd.AddCommand(ordinanceSearchCmd)
	ordinanceCmd.AddCommand(ordinanceDetailCmd)

	// Flags
	ordinanceCmd.PersistentFlags().StringVarP(&ordinanceOutputFormat, "format", "f", "table", i18n.T("ordinance.flag.format"))
	ordinanceCmd.PersistentFlags().IntVarP(&ordinancePageNo, "page", "p", 1, i18n.T("ordinance.flag.page"))
	ordinanceCmd.PersistentFlags().IntVarP(&ordinancePageSize, "size", "s", 50, i18n.T("ordinance.flag.size"))
	ordinanceCmd.PersistentFlags().StringVarP(&ordinanceRegion, "region", "r", "", i18n.T("ordinance.flag.region"))
	ordinanceCmd.PersistentFlags().StringVar(&ordinanceSort, "sort", "date", i18n.T("ordinance.flag.sort"))
}

// updateOrdinanceCommand updates ordinance command descriptions
func updateOrdinanceCommand() {
	if ordinanceCmd != nil {
		ordinanceCmd.Short = i18n.T("ordinance.short")
		ordinanceCmd.Long = i18n.T("ordinance.long")

		// Update flag descriptions
		if flag := ordinanceCmd.PersistentFlags().Lookup("format"); flag != nil {
			flag.Usage = i18n.T("ordinance.flag.format")
		}
		if flag := ordinanceCmd.PersistentFlags().Lookup("page"); flag != nil {
			flag.Usage = i18n.T("ordinance.flag.page")
		}
		if flag := ordinanceCmd.PersistentFlags().Lookup("size"); flag != nil {
			flag.Usage = i18n.T("ordinance.flag.size")
		}
		if flag := ordinanceCmd.PersistentFlags().Lookup("region"); flag != nil {
			flag.Usage = i18n.T("ordinance.flag.region")
		}
		if flag := ordinanceCmd.PersistentFlags().Lookup("sort"); flag != nil {
			flag.Usage = i18n.T("ordinance.flag.sort")
		}
	}

	// Update subcommands
	updateOrdinanceSearchCommand()
	updateOrdinanceDetailCommand()
}

// runOrdinanceSearch handles the ordinance search command
func runOrdinanceSearch(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(strings.Join(args, " "))
	if query == "" {
		logger.Debug("Empty query provided")
		return cliErrors.ErrEmptyQuery
	}

	logger.Debug("Starting ordinance search for query: %s", query)

	// Use test client if available (for testing)
	var client api.ClientInterface
	if testOrdinanceClient != nil {
		client = testOrdinanceClient
	} else {
		// Create ELIS API client
		apiClient, err := api.CreateClient(api.APITypeELIS)
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

			verbose, _ := cmd.Flags().GetBool("verbose")
			logger.LogError(err, verbose)
			return err
		}
		client = apiClient
	}

	// Get verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Search ordinances
	return searchOrdinances(client, query, ordinanceRegion, ordinanceOutputFormat, ordinancePageNo, ordinancePageSize, ordinanceSort, cmd.OutOrStdout(), verbose)
}

// searchOrdinances performs the actual ordinance search
func searchOrdinances(client api.ClientInterface, query string, region string, format string, pageNo int, pageSize int, sort string, writer io.Writer, verbose bool) error {
	// Log search parameters
	if region != "" {
		logger.Info("조례 검색 중... (검색어: %s, 지역: %s, 페이지: %d, 크기: %d)", query, region, pageNo, pageSize)
	} else {
		logger.Info("조례 검색 중... (검색어: %s, 페이지: %d, 크기: %d)", query, pageNo, pageSize)
	}

	// Create search request
	searchReq := &api.UnifiedSearchRequest{
		Query:    query,
		Region:   region,
		PageNo:   pageNo,
		PageSize: pageSize,
		Sort:     sort,
		Type:     "json",
	}

	// Perform search
	ctx := context.Background()
	result, err := client.Search(ctx, searchReq)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(writer, err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.LogError(err, verbose)
		return err
	}

	// Log search results
	logger.Info("검색 완료: %d개의 결과 (페이지: %d, 크기: %d)", result.TotalCount, pageNo, pageSize)

	// Create formatter with the specified format
	formatter := output.NewFormatter(format)

	// Format and output results
	outputStr, err := formatter.FormatSearchResultToString(result)

	if err != nil {
		logger.LogError(err, verbose)
		return err
	}

	fmt.Fprint(writer, outputStr)
	return nil
}
