package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/api"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	pageNo       int
	pageSize     int
)

// lawCmd represents the law command
var lawCmd = &cobra.Command{
	Use:   "law <검색어>",
	Short: "법령 정보 검색",
	Long: `국가법령정보센터에서 법령 정보를 검색합니다.

검색어를 입력하면 관련 법령들의 목록을 확인할 수 있습니다.
기본적으로 테이블 형식으로 출력되며, --format 옵션으로 JSON 형식도 지원합니다.`,
	Example: `  # 기본 검색
  sejong law "개인정보 보호법"
  
  # JSON 형식으로 출력
  sejong law "도로교통법" --format json
  
  # 페이지네이션 옵션
  sejong law "민법" --page 2 --size 20`,
	Args: cobra.ExactArgs(1),
	RunE: runLawCommand,
}

func init() {
	rootCmd.AddCommand(lawCmd)
	
	// Flags
	lawCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "출력 형식 (table, json)")
	lawCmd.Flags().IntVarP(&pageNo, "page", "p", 1, "페이지 번호")
	lawCmd.Flags().IntVarP(&pageSize, "size", "s", 10, "페이지 크기")
}

func runLawCommand(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(args[0])
	if query == "" {
		return fmt.Errorf("검색어를 입력해주세요")
	}
	
	// Create API client
	client, err := api.NewClient()
	if err != nil {
		// Check if it's an API key error
		if strings.Contains(err.Error(), "API 키") {
			fmt.Fprintln(os.Stderr, "❌ API 키가 설정되지 않았습니다.")
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "국가법령정보센터 오픈 API 이용을 위해 인증키가 필요합니다.")
			fmt.Fprintln(os.Stderr, "1. 인증키 발급: https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9")
			fmt.Fprintln(os.Stderr, "2. 키 설정: sejong config set law.key <발급받은_인증키>")
			return nil // Return nil to avoid printing the error twice
		}
		return fmt.Errorf("API 클라이언트 생성 실패: %w", err)
	}
	
	// Show searching message
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		fmt.Fprintf(os.Stderr, "검색 중... (%s)\n", query)
	}
	
	// Create search request
	req := &api.SearchRequest{
		Query:    query,
		Type:     api.TypeJSON,
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	
	// Search with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	resp, err := client.Search(ctx, req)
	if err != nil {
		return fmt.Errorf("검색 실패: %w", err)
	}
	
	// Output results
	formatter := output.NewFormatter(outputFormat)
	if err := formatter.FormatSearchResult(resp); err != nil {
		return fmt.Errorf("출력 실패: %w", err)
	}
	
	return nil
}