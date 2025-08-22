package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/api"
)

// Formatter handles output formatting
type Formatter struct {
	format string
}

// NewFormatter creates a new formatter with the specified format
func NewFormatter(format string) *Formatter {
	return &Formatter{
		format: strings.ToLower(format),
	}
}

// FormatSearchResult formats and outputs the search results
func (f *Formatter) FormatSearchResult(resp *api.SearchResponse) error {
	switch f.format {
	case "json":
		return f.formatJSON(resp)
	case "table", "":
		return f.formatTable(resp)
	default:
		return fmt.Errorf("지원하지 않는 출력 형식: %s (table, json 중 선택)", f.format)
	}
}

// formatJSON outputs results in JSON format
func (f *Formatter) formatJSON(resp *api.SearchResponse) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(resp)
}

// formatTable outputs results in table format
func (f *Formatter) formatTable(resp *api.SearchResponse) error {
	// Show summary
	fmt.Printf("총 %d개의 법령을 찾았습니다.\n\n", resp.TotalCount)
	
	// If no results, return early
	if len(resp.Laws) == 0 {
		fmt.Println("검색 결과가 없습니다.")
		return nil
	}
	
	// Create simple table output
	// Print header
	fmt.Printf("%-5s %-45s %-10s %-15s %-12s\n", "번호", "법령명", "법령구분", "소관부처", "시행일자")
	fmt.Println(strings.Repeat("-", 100))
	
	// Add data rows
	for i, law := range resp.Laws {
		// Format dates (YYYYMMDD -> YYYY-MM-DD)
		effectDate := formatDate(law.EffectDate)
		
		// Truncate long names for better display
		name := truncateString(law.Name, 40)
		dept := truncateString(law.Department, 13)
		
		fmt.Printf("%-5d %-45s %-10s %-15s %-12s\n",
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
		totalPages := (resp.TotalCount + len(resp.Laws) - 1) / len(resp.Laws)
		fmt.Printf("\n페이지 %d/%d (--page 옵션으로 다른 페이지 조회 가능)\n", currentPage, totalPages)
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
	if len(s) <= maxLen {
		return s
	}
	// Handle Korean characters properly by using rune slice
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}