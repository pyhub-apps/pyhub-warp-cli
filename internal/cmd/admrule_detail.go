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

var admruleDetailCmd *cobra.Command

// initAdmruleDetailCmd initializes the administrative rule detail subcommand
func initAdmruleDetailCmd() {
	admruleDetailCmd = &cobra.Command{
		Use:   "detail <행정규칙ID>",
		Short: "행정규칙 상세 조회",
		Long:  "행정규칙ID로 행정규칙의 상세 정보를 조회합니다.",
		Example: `  # 행정규칙 상세 조회
  warp admrule detail 12345
  
  # JSON 형식으로 출력
  warp admrule detail 12345 --format json`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"d", "info"},
		RunE:    runAdmruleDetailCommand,
	}

	// Flags (inherit format from parent)
	admruleDetailCmd.Flags().StringVarP(&admrOutputFormat, "format", "f", "table", "출력 형식 (table, json)")
}

// updateAdmruleDetailCommand updates administrative rule detail command descriptions
func updateAdmruleDetailCommand() {
	if admruleDetailCmd != nil {
		admruleDetailCmd.Short = "행정규칙 상세 조회"
		admruleDetailCmd.Long = "행정규칙ID로 행정규칙의 상세 정보를 조회합니다."

		// Update flag descriptions
		if flag := admruleDetailCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json)"
		}
	}
}

func runAdmruleDetailCommand(cmd *cobra.Command, args []string) error {
	// Get output writer
	output := cmd.OutOrStdout()

	// Get administrative rule ID
	admrulID := strings.TrimSpace(args[0])
	if admrulID == "" {
		logger.Debug("Empty administrative rule ID provided")
		return fmt.Errorf("행정규칙 ID를 입력해주세요")
	}

	logger.Info("행정규칙 상세 정보 조회 중... (ID: %s)", admrulID)

	// Create API client for administrative rule
	client, err := api.CreateClient(api.APITypeAdmrul)
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Get detail with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detail, err := client.GetDetail(ctx, admrulID)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(output, err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.Error("Failed to get administrative rule detail: %v", err)
		return fmt.Errorf("행정규칙 상세 조회 실패: %v", err)
	}

	logger.Info("행정규칙 상세 정보 조회 완료: %s", detail.Name)

	// Format and output results
	formatter := outputPkg.NewFormatter(admrOutputFormat)
	formattedOutput, err := formatter.FormatDetailToString(detail)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf("출력 실패")
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}
