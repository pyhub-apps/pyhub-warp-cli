package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pyhub-apps/sejong-cli/internal/api"
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

// FormatSearchResultToString formats the search results and returns as string
func (f *Formatter) FormatSearchResultToString(resp *api.SearchResponse) (string, error) {
	if resp == nil {
		return "", fmt.Errorf("검색 결과가 없습니다")
	}

	switch f.format {
	case "json":
		return f.formatJSONToString(resp)
	case "table", "":
		return f.formatTableToString(resp)
	default:
		return "", fmt.Errorf("지원하지 않는 출력 형식: %s (table, json 중 선택)", f.format)
	}
}

// FormatDetailToString formats law detail and returns as string
func (f *Formatter) FormatDetailToString(detail *api.LawDetail) (string, error) {
	if detail == nil {
		return "", fmt.Errorf("법령 상세 정보가 없습니다")
	}

	switch f.format {
	case "json":
		data, err := json.MarshalIndent(detail, "", "  ")
		if err != nil {
			return "", fmt.Errorf("JSON 변환 실패: %w", err)
		}
		return string(data) + "\n", nil
	case "table", "":
		return f.formatDetailTable(detail), nil
	default:
		return "", fmt.Errorf("지원하지 않는 출력 형식: %s (table, json 중 선택)", f.format)
	}
}

// FormatHistoryToString formats law history and returns as string
func (f *Formatter) FormatHistoryToString(history *api.LawHistory) (string, error) {
	if history == nil {
		return "", fmt.Errorf("법령 이력 정보가 없습니다")
	}

	switch f.format {
	case "json":
		data, err := json.MarshalIndent(history, "", "  ")
		if err != nil {
			return "", fmt.Errorf("JSON 변환 실패: %w", err)
		}
		return string(data) + "\n", nil
	case "table", "":
		return f.formatHistoryTable(history), nil
	default:
		return "", fmt.Errorf("지원하지 않는 출력 형식: %s (table, json 중 선택)", f.format)
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
		// Use EffectDate if available, otherwise use PromulDate (for interpretations, precedents)
		dateToShow := law.EffectDate
		if dateToShow == "" && law.PromulDate != "" {
			dateToShow = law.PromulDate
		}
		effectDate := formatDate(dateToShow)

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
		// Use a default page size of 10 if not enough items to determine
		pageSize := 10
		if len(resp.Laws) > 0 {
			pageSize = len(resp.Laws)
		}
		totalPages := (resp.TotalCount + pageSize - 1) / pageSize
		fmt.Printf("\n페이지 %d/%d (--page 옵션으로 다른 페이지 조회 가능)\n", currentPage, totalPages)
	}

	return nil
}

// formatDate converts YYYYMMDD to YYYY-MM-DD format
func formatDate(date string) string {
	// Handle YYYYMMDD format
	if len(date) == 8 {
		return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
	}
	// Handle YYYY.MM.DD format (from some APIs like legal interpretations)
	if len(date) == 10 && date[4] == '.' && date[7] == '.' {
		return strings.ReplaceAll(date, ".", "-")
	}
	// Return as-is for other formats
	return date
}

// formatJSONToString formats results in JSON format and returns as string
func (f *Formatter) formatJSONToString(resp *api.SearchResponse) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resp); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// formatTableToString formats results in table format and returns as string
func (f *Formatter) formatTableToString(resp *api.SearchResponse) (string, error) {
	var buf bytes.Buffer

	// Show summary
	fmt.Fprintf(&buf, "총 %d개의 법령을 찾았습니다.\n\n", resp.TotalCount)

	// If no results, return early
	if len(resp.Laws) == 0 {
		fmt.Fprintln(&buf, "검색 결과가 없습니다.")
		return buf.String(), nil
	}

	// Create simple table output
	// Check if we have source information (unified search)
	hasSource := false
	for _, law := range resp.Laws {
		if law.Source != "" {
			hasSource = true
			break
		}
	}

	// Print header
	if hasSource {
		fmt.Fprintf(&buf, "%-5s %-40s %-10s %-10s %-15s %-12s\n", "번호", "법령명", "구분", "출처", "소관부처", "시행일자")
		fmt.Fprintln(&buf, strings.Repeat("-", 107))
	} else {
		fmt.Fprintf(&buf, "%-5s %-45s %-10s %-15s %-12s\n", "번호", "법령명", "법령구분", "소관부처", "시행일자")
		fmt.Fprintln(&buf, strings.Repeat("-", 100))
	}

	// Add data rows
	for i, law := range resp.Laws {
		// Format dates (YYYYMMDD -> YYYY-MM-DD)
		effectDate := formatDate(law.EffectDate)

		// Truncate long names for better display
		if hasSource {
			name := truncateString(law.Name, 38)
			dept := truncateString(law.Department, 13)
			source := law.Source
			if source == "" {
				source = "-"
			}

			fmt.Fprintf(&buf, "%-5d %-40s %-10s %-10s %-15s %-12s\n",
				i+1,
				name,
				law.LawType,
				source,
				dept,
				effectDate,
			)
		} else {
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

	return buf.String(), nil
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

// formatDetailTable formats law detail as a table
func (f *Formatter) formatDetailTable(detail *api.LawDetail) string {
	var buf bytes.Buffer

	// Basic information
	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n")
	fmt.Fprintf(&buf, " 법령 상세 정보\n")
	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n\n")

	fmt.Fprintf(&buf, "법령ID:       %s\n", detail.ID)
	fmt.Fprintf(&buf, "법령명:       %s\n", detail.Name)
	if detail.NameAbbrev != "" {
		fmt.Fprintf(&buf, "약칭:         %s\n", detail.NameAbbrev)
	}
	fmt.Fprintf(&buf, "법령구분:     %s\n", detail.LawType)
	fmt.Fprintf(&buf, "소관부처:     %s\n", detail.Department)
	fmt.Fprintf(&buf, "공포일자:     %s\n", formatDate(detail.PromulDate))
	fmt.Fprintf(&buf, "공포번호:     %s\n", detail.PromulNo)
	fmt.Fprintf(&buf, "시행일자:     %s\n", formatDate(detail.EffectDate))
	fmt.Fprintf(&buf, "제개정구분:   %s\n", detail.Category)

	// Articles if present
	if len(detail.Articles) > 0 {
		fmt.Fprintf(&buf, "\n───────────────────────────────────────────────────────────\n")
		fmt.Fprintf(&buf, " 조문 (%d개)\n", len(detail.Articles))
		fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")

		for _, article := range detail.Articles {
			fmt.Fprintf(&buf, "%s", article.Number)
			if article.Title != "" {
				fmt.Fprintf(&buf, " (%s)", article.Title)
			}
			fmt.Fprintf(&buf, "\n")

			// Clean and format content
			content := strings.TrimSpace(article.Content)
			content = strings.ReplaceAll(content, "\r\n", "\n")
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					fmt.Fprintf(&buf, "  %s\n", line)
				}
			}
			fmt.Fprintf(&buf, "\n")
		}
	}

	// Related laws if present
	if len(detail.RelatedLaws) > 0 {
		fmt.Fprintf(&buf, "\n───────────────────────────────────────────────────────────\n")
		fmt.Fprintf(&buf, " 관련 법령\n")
		fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")
		for _, law := range detail.RelatedLaws {
			fmt.Fprintf(&buf, "  • %s\n", law)
		}
	}

	// Attachments if present
	if len(detail.Attachments) > 0 {
		fmt.Fprintf(&buf, "\n───────────────────────────────────────────────────────────\n")
		fmt.Fprintf(&buf, " 첨부 파일\n")
		fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")
		for _, file := range detail.Attachments {
			fmt.Fprintf(&buf, "  • %s\n", file)
		}
	}

	fmt.Fprintf(&buf, "\n═══════════════════════════════════════════════════════════\n")

	return buf.String()
}

// formatHistoryTable formats law history as a table
func (f *Formatter) formatHistoryTable(history *api.LawHistory) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n")
	fmt.Fprintf(&buf, " 법령 제/개정 이력\n")
	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n\n")

	if history.LawName != "" {
		fmt.Fprintf(&buf, "법령명: %s\n", history.LawName)
	}
	if history.LawID != "" {
		fmt.Fprintf(&buf, "법령ID: %s\n", history.LawID)
	}

	if len(history.Histories) == 0 {
		fmt.Fprintf(&buf, "\n이력이 없습니다.\n")
	} else {
		fmt.Fprintf(&buf, "\n총 %d개의 이력\n", len(history.Histories))
		fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")

		for i, record := range history.Histories {
			fmt.Fprintf(&buf, "[%d] %s - %s\n", i+1, formatDate(record.Date), record.Type)
			if record.PromulNo != "" {
				fmt.Fprintf(&buf, "    공포번호: %s\n", record.PromulNo)
			}
			if record.EffectDate != "" {
				fmt.Fprintf(&buf, "    시행일자: %s\n", formatDate(record.EffectDate))
			}
			if record.Reason != "" {
				fmt.Fprintf(&buf, "    개정이유: %s\n", record.Reason)
			}
			fmt.Fprintf(&buf, "\n")
		}
	}

	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n")

	return buf.String()
}
