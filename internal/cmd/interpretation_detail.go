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

var interpretationDetailCmd *cobra.Command

// initInterpretationDetailCmd initializes the legal interpretation detail subcommand
func initInterpretationDetailCmd() {
	interpretationDetailCmd = &cobra.Command{
		Use:   "detail <법령해석례ID>",
		Short: "법령해석례 상세 조회",
		Long:  "법령해석례ID로 법령해석례의 상세 정보를 조회합니다.",
		Example: `  # 법령해석례 상세 조회
  sejong interpretation detail 12345
  
  # JSON 형식으로 출력
  sejong interpretation detail 12345 --format json`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"d", "info"},
		RunE:    runInterpretationDetailCommand,
	}

	// Flags (inherit format from parent)
	interpretationDetailCmd.Flags().StringVarP(&interpOutputFormat, "format", "f", "table", "출력 형식 (table, json)")
}

// updateInterpretationDetailCommand updates legal interpretation detail command descriptions
func updateInterpretationDetailCommand() {
	if interpretationDetailCmd != nil {
		interpretationDetailCmd.Short = "법령해석례 상세 조회"
		interpretationDetailCmd.Long = "법령해석례ID로 법령해석례의 상세 정보를 조회합니다."

		// Update flag descriptions
		if flag := interpretationDetailCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json)"
		}
	}
}

func runInterpretationDetailCommand(cmd *cobra.Command, args []string) error {
	// Get output writer
	output := cmd.OutOrStdout()

	// Get legal interpretation ID
	expcID := strings.TrimSpace(args[0])
	if expcID == "" {
		logger.Debug("Empty legal interpretation ID provided")
		return fmt.Errorf("법령해석례 ID를 입력해주세요")
	}

	logger.Info("법령해석례 상세 정보 조회 중... (ID: %s)", expcID)

	// Create API client for legal interpretation
	client, err := api.CreateClient(api.APITypeExpc)
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Get detail with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detail, err := client.GetDetail(ctx, expcID)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(output, err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.Error("Failed to get legal interpretation detail: %v", err)
		return fmt.Errorf("법령해석례 상세 조회 실패: %v", err)
	}

	logger.Info("법령해석례 상세 정보 조회 완료: %s", detail.Name)

	// Format and output results
	formatter := outputPkg.NewFormatter(interpOutputFormat)
	formattedOutput, err := formatter.FormatDetailToString(detail)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf("출력 실패")
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}
