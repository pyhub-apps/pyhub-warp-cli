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
	case "markdown", "md":
		return f.formatMarkdown(resp)
	case "csv":
		return f.formatCSV(resp)
	case "html":
		return f.formatHTML(resp)
	default:
		return fmt.Errorf("지원하지 않는 출력 형식: %s (table, json, markdown, csv, html 중 선택)", f.format)
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
	case "markdown", "md":
		return f.formatMarkdownToString(resp)
	case "csv":
		return f.formatCSVToString(resp)
	case "html":
		return f.formatHTMLToString(resp)
	default:
		return "", fmt.Errorf("지원하지 않는 출력 형식: %s (table, json, markdown, csv, html 중 선택)", f.format)
	}
}

// FormatDetailToString formats law detail and returns as string
func (f *Formatter) FormatDetailToString(detail *api.LawDetail) (string, error) {
	return f.FormatDetailToStringWithOptions(detail, false, false, false)
}

// FormatDetailToStringWithOptions formats law detail with display options
func (f *Formatter) FormatDetailToStringWithOptions(detail *api.LawDetail, showArticles, showTables, showSupplementary bool) (string, error) {
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
		return f.formatDetailTableWithOptions(detail, showArticles, showTables, showSupplementary), nil
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

// formatTable outputs results in table format using tablewriter
func (f *Formatter) formatTable(resp *api.SearchResponse) error {
	result, err := f.formatTableToString(resp)
	if err != nil {
		return err
	}
	fmt.Print(result)
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

	// Check if we have source information (unified search)
	hasSource := false
	for _, law := range resp.Laws {
		if law.Source != "" {
			hasSource = true
			break
		}
	}

	// Prepare headers and rows
	var headers []string
	var rows [][]string

	if hasSource {
		headers = []string{"번호", "법령명", "구분", "출처", "소관부처", "시행일자"}
	} else {
		headers = []string{"번호", "법령ID", "법령명", "법령구분", "소관부처", "시행일자"}
	}

	// Add data rows
	for i, law := range resp.Laws {
		// Format dates (YYYYMMDD -> YYYY-MM-DD)
		effectDate := formatDate(law.EffectDate)
		if effectDate == "" && law.PromulDate != "" {
			effectDate = formatDate(law.PromulDate)
		}

		var row []string
		if hasSource {
			source := law.Source
			if source == "" {
				source = "-"
			}
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.Name,
				law.LawType,
				source,
				law.Department,
				effectDate,
			}
		} else {
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.ID,
				law.Name,
				law.LawType,
				law.Department,
				effectDate,
			}
		}
		rows = append(rows, row)
	}

	// Use the new table writer
	style := GetDefaultTableStyle()
	tableStr := RenderTable(headers, rows, style)
	fmt.Fprint(&buf, tableStr)

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
	return f.formatDetailTableWithOptions(detail, len(detail.Articles) > 0, false, false)
}

// formatDetailTableWithOptions formats law detail with display options
func (f *Formatter) formatDetailTableWithOptions(detail *api.LawDetail, showArticles, showTables, showSupplementary bool) string {
	var buf bytes.Buffer

	// Basic information
	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n")
	fmt.Fprintf(&buf, " 법령 상세 정보\n")
	fmt.Fprintf(&buf, "═══════════════════════════════════════════════════════════\n\n")

	// Display available information, use "N/A" for missing fields
	if detail.ID != "" {
		fmt.Fprintf(&buf, "법령ID:       %s\n", detail.ID)
	} else if detail.SerialNo != "" {
		fmt.Fprintf(&buf, "법령일련번호: %s\n", detail.SerialNo)
	}
	
	if detail.Name != "" {
		fmt.Fprintf(&buf, "법령명:       %s\n", detail.Name)
	} else {
		fmt.Fprintf(&buf, "법령명:       (정보 없음)\n")
	}
	
	if detail.NameAbbrev != "" {
		fmt.Fprintf(&buf, "약칭:         %s\n", detail.NameAbbrev)
	}
	
	if detail.LawType != "" {
		fmt.Fprintf(&buf, "법령구분:     %s\n", detail.LawType)
	}
	
	if detail.Department != "" {
		fmt.Fprintf(&buf, "소관부처:     %s\n", detail.Department)
	}
	
	if detail.PromulDate != "" {
		fmt.Fprintf(&buf, "공포일자:     %s\n", formatDate(detail.PromulDate))
	}
	
	if detail.PromulNo != "" {
		fmt.Fprintf(&buf, "공포번호:     %s\n", detail.PromulNo)
	}
	
	if detail.EffectDate != "" {
		fmt.Fprintf(&buf, "시행일자:     %s\n", formatDate(detail.EffectDate))
	}
	
	if detail.Category != "" {
		fmt.Fprintf(&buf, "제개정구분:   %s\n", detail.Category)
	}

	// Show summary of contents
	fmt.Fprintf(&buf, "\n───────────────────────────────────────────────────────────\n")
	fmt.Fprintf(&buf, " 내용 요약\n")
	fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")

	// Show counts
	if len(detail.Articles) > 0 {
		fmt.Fprintf(&buf, "조문:         %d개\n", len(detail.Articles))
	}
	if len(detail.Tables) > 0 {
		fmt.Fprintf(&buf, "별표:         %d개\n", len(detail.Tables))
	}
	if len(detail.SupplementaryProvisions) > 0 {
		fmt.Fprintf(&buf, "부칙:         %d개\n", len(detail.SupplementaryProvisions))
	}
	if detail.HasRevisionText {
		fmt.Fprintf(&buf, "개정문:       있음\n")
	}

	// Show hints for additional content
	if len(detail.Articles) > 0 && !showArticles {
		fmt.Fprintf(&buf, "\n※ 조문 상세 내용은 --articles 옵션을 사용하세요\n")
	}
	if len(detail.Tables) > 0 && !showTables {
		fmt.Fprintf(&buf, "※ 별표 내용은 --tables 옵션을 사용하세요\n")
	}
	if len(detail.SupplementaryProvisions) > 0 && !showSupplementary {
		fmt.Fprintf(&buf, "※ 부칙 내용은 --addendum 옵션을 사용하세요\n")
	}

	// Articles if present and requested
	if showArticles && len(detail.Articles) > 0 {
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

	// Tables if present and requested
	if showTables && len(detail.Tables) > 0 {
		fmt.Fprintf(&buf, "\n───────────────────────────────────────────────────────────\n")
		fmt.Fprintf(&buf, " 별표 (%d개)\n", len(detail.Tables))
		fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")

		for _, table := range detail.Tables {
			fmt.Fprintf(&buf, "%s", table.Number)
			if table.Title != "" {
				fmt.Fprintf(&buf, " - %s", table.Title)
			}
			fmt.Fprintf(&buf, "\n")

			// Clean and format content
			if table.Content != "" {
				content := strings.TrimSpace(table.Content)
				content = strings.ReplaceAll(content, "\r\n", "\n")
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						fmt.Fprintf(&buf, "  %s\n", line)
					}
				}
			}
			fmt.Fprintf(&buf, "\n")
		}
	}

	// Supplementary provisions if present and requested
	if showSupplementary && len(detail.SupplementaryProvisions) > 0 {
		fmt.Fprintf(&buf, "\n───────────────────────────────────────────────────────────\n")
		fmt.Fprintf(&buf, " 부칙 (%d개)\n", len(detail.SupplementaryProvisions))
		fmt.Fprintf(&buf, "───────────────────────────────────────────────────────────\n\n")

		for _, supp := range detail.SupplementaryProvisions {
			if supp.PromulgationDate != "" || supp.PromulgationNo != "" {
				fmt.Fprintf(&buf, "부칙")
				if supp.PromulgationNo != "" {
					fmt.Fprintf(&buf, " <%s>", supp.PromulgationNo)
				}
				if supp.PromulgationDate != "" {
					fmt.Fprintf(&buf, " (%s)", formatDate(supp.PromulgationDate))
				}
				fmt.Fprintf(&buf, "\n")
			} else if supp.Number != "" {
				fmt.Fprintf(&buf, "%s\n", supp.Number)
			}

			// Clean and format content
			if supp.Content != "" {
				content := strings.TrimSpace(supp.Content)
				content = strings.ReplaceAll(content, "\r\n", "\n")
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						fmt.Fprintf(&buf, "  %s\n", line)
					}
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

// formatMarkdown outputs results in markdown format
func (f *Formatter) formatMarkdown(resp *api.SearchResponse) error {
	result, err := f.formatMarkdownToString(resp)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

// formatMarkdownToString formats results as markdown and returns as string
func (f *Formatter) formatMarkdownToString(resp *api.SearchResponse) (string, error) {
	var buf bytes.Buffer

	// Show summary
	fmt.Fprintf(&buf, "## 검색 결과\n\n")
	fmt.Fprintf(&buf, "총 **%d**개의 법령을 찾았습니다.\n\n", resp.TotalCount)

	// If no results, return early
	if len(resp.Laws) == 0 {
		fmt.Fprintln(&buf, "_검색 결과가 없습니다._")
		return buf.String(), nil
	}

	// Check if we have source information
	hasSource := false
	for _, law := range resp.Laws {
		if law.Source != "" {
			hasSource = true
			break
		}
	}

	// Prepare headers and rows
	var headers []string
	var rows [][]string

	if hasSource {
		headers = []string{"번호", "법령명", "구분", "출처", "소관부처", "시행일자"}
	} else {
		headers = []string{"번호", "법령ID", "법령명", "법령구분", "소관부처", "시행일자"}
	}

	// Add data rows
	for i, law := range resp.Laws {
		effectDate := formatDate(law.EffectDate)
		if effectDate == "" && law.PromulDate != "" {
			effectDate = formatDate(law.PromulDate)
		}

		var row []string
		if hasSource {
			source := law.Source
			if source == "" {
				source = "-"
			}
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.Name,
				law.LawType,
				source,
				law.Department,
				effectDate,
			}
		} else {
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.ID,
				law.Name,
				law.LawType,
				law.Department,
				effectDate,
			}
		}
		rows = append(rows, row)
	}

	// Render markdown table
	tableStr := RenderMarkdownTable(headers, rows)
	fmt.Fprint(&buf, tableStr)

	// Show pagination info
	if resp.TotalCount > len(resp.Laws) {
		currentPage := resp.Page
		pageSize := 10
		if len(resp.Laws) > 0 {
			pageSize = len(resp.Laws)
		}
		totalPages := (resp.TotalCount + pageSize - 1) / pageSize
		fmt.Fprintf(&buf, "\n> 페이지 %d/%d (--page 옵션으로 다른 페이지 조회 가능)\n", currentPage, totalPages)
	}

	return buf.String(), nil
}

// formatCSV outputs results in CSV format
func (f *Formatter) formatCSV(resp *api.SearchResponse) error {
	result, err := f.formatCSVToString(resp)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

// formatCSVToString formats results as CSV and returns as string
func (f *Formatter) formatCSVToString(resp *api.SearchResponse) (string, error) {
	if len(resp.Laws) == 0 {
		return "", nil
	}

	// Check if we have source information
	hasSource := false
	for _, law := range resp.Laws {
		if law.Source != "" {
			hasSource = true
			break
		}
	}

	// Prepare headers and rows
	var headers []string
	var rows [][]string

	if hasSource {
		headers = []string{"번호", "법령명", "구분", "출처", "소관부처", "시행일자"}
	} else {
		headers = []string{"번호", "법령ID", "법령명", "법령구분", "소관부처", "시행일자"}
	}

	// Add data rows
	for i, law := range resp.Laws {
		effectDate := formatDate(law.EffectDate)
		if effectDate == "" && law.PromulDate != "" {
			effectDate = formatDate(law.PromulDate)
		}

		var row []string
		if hasSource {
			source := law.Source
			if source == "" {
				source = "-"
			}
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.Name,
				law.LawType,
				source,
				law.Department,
				effectDate,
			}
		} else {
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.ID,
				law.Name,
				law.LawType,
				law.Department,
				effectDate,
			}
		}
		rows = append(rows, row)
	}

	// Render CSV with BOM for Excel compatibility
	return RenderCSV(headers, rows, true)
}

// formatHTML outputs results in HTML format
func (f *Formatter) formatHTML(resp *api.SearchResponse) error {
	result, err := f.formatHTMLToString(resp)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

// formatHTMLToString formats results as HTML and returns as string
func (f *Formatter) formatHTMLToString(resp *api.SearchResponse) (string, error) {
	var buf bytes.Buffer

	// HTML header
	fmt.Fprintln(&buf, `<!DOCTYPE html>`)
	fmt.Fprintln(&buf, `<html lang="ko">`)
	fmt.Fprintln(&buf, `<head>`)
	fmt.Fprintln(&buf, `  <meta charset="UTF-8">`)
	fmt.Fprintln(&buf, `  <title>법령 검색 결과</title>`)
	fmt.Fprintln(&buf, `</head>`)
	fmt.Fprintln(&buf, `<body>`)

	// Summary
	fmt.Fprintf(&buf, `  <h2>검색 결과</h2>%s`, "\n")
	fmt.Fprintf(&buf, `  <p>총 <strong>%d</strong>개의 법령을 찾았습니다.</p>%s`, resp.TotalCount, "\n")

	// If no results, return early
	if len(resp.Laws) == 0 {
		fmt.Fprintln(&buf, `  <p><em>검색 결과가 없습니다.</em></p>`)
		fmt.Fprintln(&buf, `</body>`)
		fmt.Fprintln(&buf, `</html>`)
		return buf.String(), nil
	}

	// Check if we have source information
	hasSource := false
	for _, law := range resp.Laws {
		if law.Source != "" {
			hasSource = true
			break
		}
	}

	// Prepare headers and rows
	var headers []string
	var rows [][]string

	if hasSource {
		headers = []string{"번호", "법령명", "구분", "출처", "소관부처", "시행일자"}
	} else {
		headers = []string{"번호", "법령ID", "법령명", "법령구분", "소관부처", "시행일자"}
	}

	// Add data rows
	for i, law := range resp.Laws {
		effectDate := formatDate(law.EffectDate)
		if effectDate == "" && law.PromulDate != "" {
			effectDate = formatDate(law.PromulDate)
		}

		var row []string
		if hasSource {
			source := law.Source
			if source == "" {
				source = "-"
			}
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.Name,
				law.LawType,
				source,
				law.Department,
				effectDate,
			}
		} else {
			row = []string{
				fmt.Sprintf("%d", i+1),
				law.ID,
				law.Name,
				law.LawType,
				law.Department,
				effectDate,
			}
		}
		rows = append(rows, row)
	}

	// Render HTML table
	tableStr := RenderHTMLTable(headers, rows)
	fmt.Fprintln(&buf, tableStr)

	// Pagination info
	if resp.TotalCount > len(resp.Laws) {
		currentPage := resp.Page
		pageSize := 10
		if len(resp.Laws) > 0 {
			pageSize = len(resp.Laws)
		}
		totalPages := (resp.TotalCount + pageSize - 1) / pageSize
		fmt.Fprintf(&buf, `  <p style="margin-top: 20px; color: #666;">페이지 %d/%d</p>%s`, currentPage, totalPages, "\n")
	}

	fmt.Fprintln(&buf, `</body>`)
	fmt.Fprintln(&buf, `</html>`)

	return buf.String(), nil
}
