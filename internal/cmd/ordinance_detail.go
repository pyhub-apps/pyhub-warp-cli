package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pyhub-apps/sejong-cli/internal/api"
	cliErrors "github.com/pyhub-apps/sejong-cli/internal/errors"
	"github.com/pyhub-apps/sejong-cli/internal/i18n"
	"github.com/pyhub-apps/sejong-cli/internal/logger"
	"github.com/pyhub-apps/sejong-cli/internal/onboarding"
	"github.com/pyhub-apps/sejong-cli/internal/output"
	"github.com/spf13/cobra"
)

// ordinanceDetailCmd represents the ordinance detail subcommand
var ordinanceDetailCmd *cobra.Command

// initOrdinanceDetailCmd initializes the ordinance detail subcommand
func initOrdinanceDetailCmd() {
	ordinanceDetailCmd = &cobra.Command{
		Use:   "detail <조례ID>",
		Short: i18n.T("ordinance.detail.short"),
		Long:  i18n.T("ordinance.detail.long"),
		Example: `  # 조례 상세 조회
  sejong ordinance detail ORD123456
  
  # JSON 형식으로 출력
  sejong ordinance detail ORD123456 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: runOrdinanceDetailCommand,
	}

	// Flags are inherited from parent command
}

// updateOrdinanceDetailCommand updates ordinance detail command descriptions
func updateOrdinanceDetailCommand() {
	if ordinanceDetailCmd != nil {
		ordinanceDetailCmd.Short = i18n.T("ordinance.detail.short")
		ordinanceDetailCmd.Long = i18n.T("ordinance.detail.long")
	}
}

func runOrdinanceDetailCommand(cmd *cobra.Command, args []string) error {
	// Get ordinance ID
	ordinanceID := strings.TrimSpace(args[0])
	if ordinanceID == "" {
		logger.Debug("Empty ordinance ID provided")
		return fmt.Errorf("조례 ID를 입력해주세요")
	}

	logger.Info("조례 상세 정보 조회 중... (ID: %s)", ordinanceID)

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
				logger.Error("Failed to create API client: %v", err)
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil
			}

			// Also check for direct API key error message
			if strings.Contains(err.Error(), "API 키가 설정되지 않았습니다") {
				logger.Error("Failed to create API client: %v", err)
				guide := onboarding.NewGuideWithWriter(cmd.OutOrStdout(), false)
				guide.ShowAPIKeySetup()
				return nil
			}

			logger.Error("Failed to create API client: %v", err)
			return err
		}
		client = apiClient
	}

	// Get verbose flag from parent command
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Get detail from API
	ctx := context.Background()
	detail, err := client.GetDetail(ctx, ordinanceID)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(cmd.OutOrStdout(), err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.LogError(err, verbose)
		return err
	}

	// Get format flag
	format, _ := cmd.Flags().GetString("format")

	// Create formatter with the specified format
	formatter := output.NewFormatter(format)

	// Format and output result
	outputStr, err := formatter.FormatDetailToString(detail)

	if err != nil {
		logger.LogError(err, verbose)
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), outputStr)
	return nil
}
