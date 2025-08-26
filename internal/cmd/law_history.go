package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/api"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/i18n"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
	outputPkg "github.com/pyhub-apps/pyhub-warp-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	lawHistoryCmd *cobra.Command
	historyLimit  int
)

// initLawHistoryCmd initializes the law history command
func initLawHistoryCmd() {
	lawHistoryCmd = &cobra.Command{
		Use:   "history <법령ID>",
		Short: i18n.T("law.history.short"),
		Long:  i18n.T("law.history.long"),
		Example: `  # 법령ID로 이력 조회
  warp law history 001234
  
  # 최근 10개만 조회
  warp law history 001234 --limit 10
  
  # JSON 형식으로 출력
  warp law history 001234 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: runLawHistoryCommand,
	}

	// Flags
	lawHistoryCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", i18n.T("law.flag.format"))
	lawHistoryCmd.Flags().IntVarP(&historyLimit, "limit", "l", 0, i18n.T("law.history.flag.limit"))
}

// updateLawHistoryCommand updates law history command descriptions
func updateLawHistoryCommand() {
	if lawHistoryCmd != nil {
		lawHistoryCmd.Short = i18n.T("law.history.short")
		lawHistoryCmd.Long = i18n.T("law.history.long")

		// Update flag descriptions
		if flag := lawHistoryCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = i18n.T("law.flag.format")
		}
		if flag := lawHistoryCmd.Flags().Lookup("limit"); flag != nil {
			flag.Usage = i18n.T("law.history.flag.limit")
		}
	}
}

func runLawHistoryCommand(cmd *cobra.Command, args []string) error {
	// Get law ID
	lawID := strings.TrimSpace(args[0])
	if lawID == "" {
		return fmt.Errorf(i18n.T("law.history.error.emptyID"))
	}

	logger.Info(i18n.Tf("law.history.searching", lawID))

	// Create API client
	client, err := api.CreateDefaultClient()
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Get law history with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	history, err := client.GetHistory(ctx, lawID)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(cmd.OutOrStdout(), err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.Error("Failed to get law history: %v", err)
		return fmt.Errorf(i18n.T("law.history.error.failed"), err)
	}

	// Apply limit if specified
	if historyLimit > 0 && len(history.Histories) > historyLimit {
		history.Histories = history.Histories[:historyLimit]
	}

	logger.Info(i18n.Tf("law.history.searchComplete", len(history.Histories)))

	// Format and output results
	formatter := outputPkg.NewFormatter(outputFormat)
	formattedOutput, err := formatter.FormatHistoryToString(history)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf(i18n.T("law.outputFailed"))
	}

	// Write formatted output
	fmt.Fprint(cmd.OutOrStdout(), formattedOutput)

	return nil
}
