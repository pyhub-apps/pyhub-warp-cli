package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pyhub-apps/sejong-cli/internal/api"
	"github.com/pyhub-apps/sejong-cli/internal/i18n"
	"github.com/pyhub-apps/sejong-cli/internal/logger"
	outputPkg "github.com/pyhub-apps/sejong-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	lawDetailCmd *cobra.Command
	showArticles bool
)

// initLawDetailCmd initializes the law detail command
func initLawDetailCmd() {
	lawDetailCmd = &cobra.Command{
		Use:   "detail <법령ID>",
		Short: i18n.T("law.detail.short"),
		Long:  i18n.T("law.detail.long"),
		Example: `  # 법령ID로 상세 조회
  sejong law detail 001234
  
  # 조문 포함하여 조회
  sejong law detail 001234 --articles
  
  # JSON 형식으로 출력
  sejong law detail 001234 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: runLawDetailCommand,
	}

	// Flags
	lawDetailCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", i18n.T("law.flag.format"))
	lawDetailCmd.Flags().BoolVarP(&showArticles, "articles", "a", false, i18n.T("law.detail.flag.articles"))
}

// updateLawDetailCommand updates law detail command descriptions
func updateLawDetailCommand() {
	if lawDetailCmd != nil {
		lawDetailCmd.Short = i18n.T("law.detail.short")
		lawDetailCmd.Long = i18n.T("law.detail.long")

		// Update flag descriptions
		if flag := lawDetailCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = i18n.T("law.flag.format")
		}
		if flag := lawDetailCmd.Flags().Lookup("articles"); flag != nil {
			flag.Usage = i18n.T("law.detail.flag.articles")
		}
	}
}

func runLawDetailCommand(cmd *cobra.Command, args []string) error {
	// Get law ID
	lawID := strings.TrimSpace(args[0])
	if lawID == "" {
		return fmt.Errorf(i18n.T("law.detail.error.emptyID"))
	}

	logger.Info(i18n.Tf("law.detail.searching", lawID))

	// Create API client
	client, err := api.CreateDefaultClient()
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Get law detail with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detail, err := client.GetDetail(ctx, lawID)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(cmd.OutOrStdout(), err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.Error("Failed to get law detail: %v", err)
		return fmt.Errorf(i18n.T("law.detail.error.failed"), err)
	}

	// Use name if available, otherwise use ID or serial number
	nameToShow := detail.Name
	if nameToShow == "" {
		if detail.ID != "" {
			nameToShow = fmt.Sprintf("법령ID: %s", detail.ID)
		} else if detail.SerialNo != "" {
			nameToShow = fmt.Sprintf("법령일련번호: %s", detail.SerialNo)
		} else {
			nameToShow = "법령"
		}
	}
	logger.Info(i18n.Tf("law.detail.searchComplete", nameToShow))

	// Format and output results
	formatter := outputPkg.NewFormatter(outputFormat)

	// Filter articles if not requested
	if !showArticles {
		detail.Articles = nil
	}

	formattedOutput, err := formatter.FormatDetailToString(detail)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf(i18n.T("law.outputFailed"))
	}

	// Write formatted output
	fmt.Fprint(cmd.OutOrStdout(), formattedOutput)

	return nil
}
