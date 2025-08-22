package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/api"
	cliErrors "github.com/pyhub-kr/pyhub-sejong-cli/internal/errors"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/logger"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/onboarding"
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

// APIClient interface for dependency injection and testing
type APIClient interface {
	Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error)
}

func runLawCommand(cmd *cobra.Command, args []string) error {
	// Get search query
	query := strings.TrimSpace(args[0])
	if query == "" {
		logger.Debug("Empty query provided")
		return cliErrors.ErrEmptyQuery
	}
	
	logger.Debug("Starting law search for query: %s", query)
	
	// Create API client
	client, err := api.NewClient()
	if err != nil {
		// Check if it's an API key error
		var cliErr *cliErrors.CLIError
		if errors.As(err, &cliErr) && cliErr.Code == cliErrors.ErrCodeNoAPIKey {
			guide := onboarding.NewGuide()
			guide.ShowAPIKeySetup()
			return nil // Return nil to avoid printing the error twice
		}
		
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger.LogError(err, verbose)
		return err
	}
	
	// Use searchLaws for the actual search logic
	return searchLaws(client, query, outputFormat, pageNo, pageSize, cmd.OutOrStdout())
}

// searchLaws performs the actual law search - extracted for testing
func searchLaws(client APIClient, query string, format string, page int, size int, output io.Writer) error {
	logger.Info("Searching for: %s (page: %d, size: %d)", query, page, size)
	
	// Create search request
	req := &api.SearchRequest{
		Query:    query,
		Type:     api.TypeJSON,
		PageNo:   page,
		PageSize: size,
	}
	
	// Search with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	resp, err := client.Search(ctx, req)
	if err != nil {
		logger.LogError(err, false)
		
		// Show user-friendly error with hint
		var cliErr *cliErrors.CLIError
		if errors.As(err, &cliErr) {
			guide := onboarding.NewGuide()
			guide.ShowError(err.Error())
			return nil // Error already displayed
		}
		return err
	}
	
	logger.Info("Search completed: %d results found", resp.TotalCount)
	
	// Output results based on format
	if format == "json" {
		// JSON output
		encoder := json.NewEncoder(output)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(resp); err != nil {
			logger.Error("Failed to format output: %v", err)
			return cliErrors.Wrap(err, cliErrors.New(
				cliErrors.ErrCodeDataFormat,
				"출력 실패",
				"출력 형식을 확인하세요",
			))
		}
	} else {
		// Table output - write directly without formatter
		var buf strings.Builder
		
		// Show summary
		fmt.Fprintf(&buf, "총 %d개의 법령을 찾았습니다.\n\n", resp.TotalCount)
		
		// If no results, return early
		if len(resp.Laws) == 0 {
			fmt.Fprintln(&buf, "검색 결과가 없습니다.")
			fmt.Fprint(output, buf.String())
			return nil
		}
		
		// Create simple table output
		// Print header
		fmt.Fprintf(&buf, "%-5s %-45s %-10s %-15s %-12s\n", "번호", "법령명", "법령구분", "소관부처", "시행일자")
		fmt.Fprintln(&buf, strings.Repeat("-", 100))
		
		// Add data rows
		for i, law := range resp.Laws {
			// Format dates (YYYYMMDD -> YYYY-MM-DD)
			effectDate := formatDate(law.EffectDate)
			
			// Truncate long names for better display
			name := truncateString(law.Name, 40)
			dept := truncateString(law.Department, 13)
			
			fmt.Fprintf(&buf, "%-5d %-45s %-10s %-15s %-12s\n",
				i+1,
				name,
				law.LawType,
				dept,
				effectDate,
			)
		}
		
		// Show pagination info if there are more results
		if resp.TotalCount > len(resp.Laws) {
			currentPage := resp.Page
			// Use a default page size of 10 if not enough items to determine
			pageSize := 10
			if len(resp.Laws) > 0 {
				pageSize = len(resp.Laws)
			}
			totalPages := (resp.TotalCount + pageSize - 1) / pageSize
			fmt.Fprintf(&buf, "\n페이지 %d/%d (--page 옵션으로 다른 페이지 조회 가능)\n", currentPage, totalPages)
		}
		
		fmt.Fprint(output, buf.String())
	}
	
	return nil
}

// formatDate converts YYYYMMDD to YYYY-MM-DD format
func formatDate(date string) string {
	if len(date) != 8 {
		return date
	}
	return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
}

// truncateString truncates a string to maxLen and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	
	// Handle Unicode characters properly by using rune slice
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	
	// Ensure we don't underflow when adding ellipsis
	if maxLen <= 3 {
		// Return just ellipsis dots up to maxLen
		return "..."[:maxLen]
	}
	
	return string(runes[:maxLen-3]) + "..."
}