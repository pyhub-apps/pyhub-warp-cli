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

var precedentDetailCmd *cobra.Command

// initPrecedentDetailCmd initializes the precedent detail subcommand
func initPrecedentDetailCmd() {
	precedentDetailCmd = &cobra.Command{
		Use:   "detail <판례ID>",
		Short: "판례 상세 조회",
		Long:  "판례ID로 판례의 상세 정보를 조회합니다.",
		Example: `  # 판례 상세 조회
  sejong precedent detail 12345
  
  # JSON 형식으로 출력
  sejong precedent detail 12345 --format json`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"d", "info"},
		RunE:    runPrecedentDetailCommand,
	}

	// Flags (inherit format from parent)
	precedentDetailCmd.Flags().StringVarP(&precOutputFormat, "format", "f", "table", "출력 형식 (table, json)")
}

// updatePrecedentDetailCommand updates precedent detail command descriptions
func updatePrecedentDetailCommand() {
	if precedentDetailCmd != nil {
		precedentDetailCmd.Short = "판례 상세 조회"
		precedentDetailCmd.Long = "판례ID로 판례의 상세 정보를 조회합니다."

		// Update flag descriptions
		if flag := precedentDetailCmd.Flags().Lookup("format"); flag != nil {
			flag.Usage = "출력 형식 (table, json)"
		}
	}
}

func runPrecedentDetailCommand(cmd *cobra.Command, args []string) error {
	// Get output writer
	output := cmd.OutOrStdout()

	// Get precedent ID
	precID := strings.TrimSpace(args[0])
	if precID == "" {
		logger.Debug("Empty precedent ID provided")
		return fmt.Errorf("판례 ID를 입력해주세요")
	}

	logger.Info("판례 상세 정보 조회 중... (ID: %s)", precID)

	// Create API client for precedent
	client, err := api.CreateClient(api.APITypePrec)
	if err != nil {
		logger.Error("Failed to create API client: %v", err)
		return err
	}

	// Get detail with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detail, err := client.GetDetail(ctx, precID)
	if err != nil {
		// Check if it's an API key error
		var apiKeyErr *api.APIKeyError
		if errors.As(err, &apiKeyErr) {
			// Print error message without help
			fmt.Fprintln(output, err.Error())
			// Return nil to suppress both error message and help
			return nil
		}

		logger.Error("Failed to get precedent detail: %v", err)
		return fmt.Errorf("판례 상세 조회 실패: %v", err)
	}

	logger.Info("판례 상세 정보 조회 완료: %s", detail.Name)

	// Format and output results
	formatter := outputPkg.NewFormatter(precOutputFormat)
	formattedOutput, err := formatter.FormatDetailToString(detail)
	if err != nil {
		logger.Error("Failed to format output: %v", err)
		return fmt.Errorf("출력 실패")
	}

	// Write formatted output
	fmt.Fprint(output, formattedOutput)

	return nil
}
