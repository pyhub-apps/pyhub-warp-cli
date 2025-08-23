package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/api"
	cliErrors "github.com/pyhub-kr/pyhub-sejong-cli/internal/errors"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/i18n"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/logger"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/onboarding"
	outputPkg "github.com/pyhub-kr/pyhub-sejong-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	lawSearchCmd *cobra.Command
)

// initLawSearchCmd initializes the law search command
func initLawSearchCmd() {
	lawSearchCmd = &cobra.Command{
		Use:   "search <검색어>",
		Short: i18n.T("law.search.short"),
		Long:  i18n.T("law.search.long"),
		Example: `  # 기본 검색
  sejong law search "개인정보 보호법"
  
  # JSON 형식으로 출력
  sejong law search "도로교통법" --format json
  
  # 페이지네이션 옵션
  sejong law search "민법" --page 2 --size 20`,
		Args: cobra.ExactArgs(1),
		RunE: runLawSearchCommand,
	}

	// Flags
	lawSearchCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", i18n.T("law.flag.format"))
	lawSearchCmd.Flags().IntVarP(&pageNo, "page", "p", 1, i18n.T("law.flag.page"))
	lawSearchCmd.Flags().IntVarP(&pageSize, "size", "s", 10, i18n.T("law.flag.size"))
}

// updateLawSearchCommand updates law search command descriptions
func updateLawSearchCommand() {
	if lawSearchCmd != nil {
		lawSearchCmd.Short = i18n.T("law.search.short")
		lawSearchCmd.Long = i18n.T("law.search.long")

		// Update flag descriptions
		if flag := lawSearchCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = i18n.T("law.flag.format")
		}
		if flag := lawSearchCmd.Flags().Lookup("page"); flag != nil {
			flag.Usage = i18n.T("law.flag.page")
		}
		if flag := lawSearchCmd.Flags().Lookup("size"); flag != nil {
			flag.Usage = i18n.T("law.flag.size")
		}
	}
}

func runLawSearchCommand(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(args[0])
	if query == "" {
		logger.Debug("Empty query provided")
		return cliErrors.ErrEmptyQuery
	}

	logger.Debug("Starting law search for query: %s", query)

	// Use test client if available (for testing)
	var client APIClient
	if testAPIClient != nil {
		client = testAPIClient
	} else {
		// Create API client using the new factory
		apiClient, err := api.CreateDefaultClient()
		if err != nil {
			// Check if it's an API key error
			var cliErr *cliErrors.CLIError
			if errors.As(err, &cliErr) && cliErr.Code == cliErrors.ErrCodeNoAPIKey {
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil // Return nil to avoid printing the error twice
			}

			verbose, _ := cmd.Flags().GetBool("verbose")
			logger.LogError(err, verbose)
			return err
		}
		client = apiClient
	}

	// Get verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Use searchLaws for the actual search logic
	return searchLaws(client, query, outputFormat, pageNo, pageSize, cmd.OutOrStdout(), verbose)
}

// searchLaws performs the actual law search - reused from law.go
func searchLaws(client APIClient, query string, format string, page int, size int, output io.Writer, verbose bool) error {
	logger.Info(i18n.Tf("law.searching", query, page, size))

	// Create search request
	req := &api.UnifiedSearchRequest{
		Query:    query,
		Type:     "JSON",
		PageNo:   page,
		PageSize: size,
	}

	// Search with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Search(ctx, req)
	if err != nil {
		logger.LogError(err, verbose)

		// Show user-friendly error with hint
		var cliErr *cliErrors.CLIError
		if errors.As(err, &cliErr) {
			guide := onboarding.NewGuideWithWriter(output, false)
			guide.ShowError(err.Error())
			return nil // Error already displayed
		}
		return err
	}

	logger.Info(i18n.Tf("law.searchComplete", resp.TotalCount, page, size))

	// Format and output results using the formatter package
	formatter := outputPkg.NewFormatter(format)
	formattedOutput, err := formatter.FormatSearchResultToString(resp)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return cliErrors.Wrap(err, cliErrors.New(
			cliErrors.ErrCodeDataFormat,
			i18n.T("law.outputFailed"),
			i18n.T("law.checkFormat"),
		))
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}