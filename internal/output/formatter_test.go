package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/api"
)

// captureStdout captures stdout during function execution
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"table", "table"},
		{"TABLE", "table"},
		{"json", "json"},
		{"JSON", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := NewFormatter(tt.input)
			if f.format != tt.expected {
				t.Errorf("NewFormatter(%q).format = %q, want %q", tt.input, f.format, tt.expected)
			}
		})
	}
}

func TestFormatSearchResult_JSON(t *testing.T) {
	// Test data
	resp := &api.SearchResponse{
		TotalCount: 2,
		Page:       1,
		Laws: []api.LawInfo{
			{
				ID:         "001234",
				Name:       "개인정보 보호법",
				LawType:    "법률",
				Department: "개인정보보호위원회",
				EffectDate: "20200805",
			},
			{
				ID:         "001235",
				Name:       "개인정보 보호법 시행령",
				LawType:    "대통령령",
				Department: "개인정보보호위원회",
				EffectDate: "20210202",
			},
		},
	}

	// Capture output
	var err error
	output := captureStdout(func() {
		f := NewFormatter("json")
		err = f.FormatSearchResult(resp)
	})
	
	if err != nil {
		t.Fatalf("FormatSearchResult failed: %v", err)
	}

	// Parse JSON output
	var result api.SearchResponse
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify
	if result.TotalCount != resp.TotalCount {
		t.Errorf("TotalCount = %d, want %d", result.TotalCount, resp.TotalCount)
	}
	if len(result.Laws) != len(resp.Laws) {
		t.Errorf("Laws count = %d, want %d", len(result.Laws), len(resp.Laws))
	}
}

func TestFormatSearchResult_Table(t *testing.T) {
	// Test data
	resp := &api.SearchResponse{
		TotalCount: 2,
		Page:       1,
		Laws: []api.LawInfo{
			{
				Name:       "테스트 법령",
				LawType:    "법률",
				Department: "테스트부",
				EffectDate: "20231201",
			},
		},
	}

	// Capture output
	var err error
	output := captureStdout(func() {
		f := NewFormatter("table")
		err = f.FormatSearchResult(resp)
	})
	
	if err != nil {
		t.Fatalf("FormatSearchResult failed: %v", err)
	}

	// Verify output contains expected elements
	if !strings.Contains(output, "총 2개의 법령을 찾았습니다") {
		t.Error("Output should contain total count message")
	}
	if !strings.Contains(output, "테스트 법령") {
		t.Error("Output should contain law name")
	}
	if !strings.Contains(output, "2023-12-01") {
		t.Error("Output should contain formatted date")
	}
}

func TestFormatSearchResult_EmptyResult(t *testing.T) {
	resp := &api.SearchResponse{
		TotalCount: 0,
		Page:       1,
		Laws:       []api.LawInfo{},
	}

	// Capture output
	var err error
	output := captureStdout(func() {
		f := NewFormatter("table")
		err = f.FormatSearchResult(resp)
	})
	
	if err != nil {
		t.Fatalf("FormatSearchResult failed: %v", err)
	}

	// Verify
	if !strings.Contains(output, "검색 결과가 없습니다") {
		t.Error("Output should contain 'no results' message")
	}
}

func TestFormatSearchResult_InvalidFormat(t *testing.T) {
	resp := &api.SearchResponse{}
	f := NewFormatter("invalid")
	err := f.FormatSearchResult(resp)
	if err == nil {
		t.Error("Expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "지원하지 않는 출력 형식") {
		t.Errorf("Error message should mention unsupported format, got: %v", err)
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"20231201", "2023-12-01"},
		{"20200805", "2020-08-05"},
		{"2023", "2023"},        // Invalid format, return as-is
		{"", ""},                 // Empty string
		{"not-a-date", "not-a-date"}, // Invalid input
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"한글테스트입니다", 5, "한글..."},
		{"", 10, ""},
		{"abcdef", 3, "..."},
		{"abcdef", 2, ".."},
		{"abcdef", 1, "."},
		{"abcdef", 0, ""},
		{"test", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestFormatSearchResultToString(t *testing.T) {
	resp := &api.SearchResponse{
		TotalCount: 2,
		Page:       1,
		Laws: []api.LawInfo{
			{
				ID:         "001",
				Name:       "개인정보 보호법",
				NameAbbrev: "개인정보보호법",
				SerialNo:   "12345",
				PromulDate: "20110329",
				PromulNo:   "제10465호",
				Category:   "제정",
				Department: "개인정보보호위원회",
				EffectDate: "20110930",
				LawType:    "법률",
			},
			{
				ID:         "002",
				Name:       "정보통신망 이용촉진 및 정보보호 등에 관한 법률 시행령",
				NameAbbrev: "정보통신망법",
				SerialNo:   "67890",
				PromulDate: "20200610",
				PromulNo:   "제17344호",
				Category:   "일부개정",
				Department: "과학기술정보통신부",
				EffectDate: "20201210",
				LawType:    "법률",
			},
		},
	}

	tests := []struct {
		name       string
		format     string
		resp       *api.SearchResponse
		wantErr    bool
		contains   []string
	}{
		{
			name:    "Table format",
			format:  "table",
			resp:    resp,
			wantErr: false,
			contains: []string{
				"총 2개의 법령을 찾았습니다",
				"개인정보 보호법",
				"법률",
				"2011-09-30",
			},
		},
		{
			name:    "JSON format",
			format:  "json",
			resp:    resp,
			wantErr: false,
			contains: []string{
				`"totalCnt": 2`,
				`"법령명한글": "개인정보 보호법"`,
			},
		},
		{
			name:    "Invalid format",
			format:  "xml",
			resp:    resp,
			wantErr: true,
		},
		{
			name:   "Empty results",
			format: "table",
			resp: &api.SearchResponse{
				TotalCount: 0,
				Page:       1,
				Laws:       []api.LawInfo{},
			},
			wantErr: false,
			contains: []string{
				"총 0개의 법령을 찾았습니다",
				"검색 결과가 없습니다",
			},
		},
		{
			name:   "Pagination info",
			format: "table",
			resp: &api.SearchResponse{
				TotalCount: 100,
				Page:       2,
				Laws: []api.LawInfo{
					{
						ID:         "001",
						Name:       "테스트 법령",
						Department: "테스트부",
						EffectDate: "20240101",
						LawType:    "법률",
					},
				},
			},
			wantErr: false,
			contains: []string{
				"페이지 2/100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(tt.format)
			result, err := f.FormatSearchResultToString(tt.resp)

			if (err != nil) != tt.wantErr {
				t.Errorf("FormatSearchResultToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, substr := range tt.contains {
					if !strings.Contains(result, substr) {
						t.Errorf("Result should contain %q, got:\n%s", substr, result)
					}
				}

				// For JSON format, validate it's proper JSON
				if tt.format == "json" {
					var js map[string]interface{}
					if err := json.Unmarshal([]byte(result), &js); err != nil {
						t.Errorf("JSON output is invalid: %v", err)
					}
				}
			}
		})
	}
}

func TestFormatJSONToString(t *testing.T) {
	resp := &api.SearchResponse{
		TotalCount: 1,
		Page:       1,
		Laws: []api.LawInfo{
			{
				ID:   "001",
				Name: "테스트 법령",
			},
		},
	}

	f := NewFormatter("json")
	result, err := f.formatJSONToString(resp)
	if err != nil {
		t.Fatalf("formatJSONToString() error = %v", err)
	}

	// Check it's valid JSON
	var parsed api.SearchResponse
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Errorf("JSON output is invalid: %v", err)
	}

	// Check indentation (should have 2 spaces)
	if !strings.Contains(result, "  ") {
		t.Error("JSON should be indented with 2 spaces")
	}
}

func TestFormatTableToString(t *testing.T) {
	tests := []struct {
		name     string
		resp     *api.SearchResponse
		contains []string
	}{
		{
			name: "With results",
			resp: &api.SearchResponse{
				TotalCount: 1,
				Page:       1,
				Laws: []api.LawInfo{
					{
						Name:       "테스트 법령",
						Department: "테스트부",
						EffectDate: "20240101",
						LawType:    "법률",
					},
				},
			},
			contains: []string{
				"총 1개의 법령을 찾았습니다",
				"테스트 법령",
				"테스트부",
				"2024-01-01",
				"법률",
			},
		},
		{
			name: "Long names truncated",
			resp: &api.SearchResponse{
				TotalCount: 1,
				Page:       1,
				Laws: []api.LawInfo{
					{
						Name:       "매우 긴 법령 이름입니다. 이것은 정말로 너무 길어서 잘려야 하는 법령 이름입니다.",
						Department: "매우 긴 부처명입니다. 이것도 잘려야 합니다.",
						EffectDate: "20240101",
						LawType:    "법률",
					},
				},
			},
			contains: []string{
				"...", // Should have ellipsis for truncated text
			},
		},
		{
			name: "No results",
			resp: &api.SearchResponse{
				TotalCount: 0,
				Page:       1,
				Laws:       []api.LawInfo{},
			},
			contains: []string{
				"총 0개의 법령을 찾았습니다",
				"검색 결과가 없습니다",
			},
		},
		{
			name: "Pagination",
			resp: &api.SearchResponse{
				TotalCount: 50,
				Page:       3,
				Laws:       make([]api.LawInfo, 10),
			},
			contains: []string{
				"페이지 3/5",
				"--page 옵션으로 다른 페이지 조회 가능",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter("table")
			result, err := f.formatTableToString(tt.resp)
			if err != nil {
				t.Fatalf("formatTableToString() error = %v", err)
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Result should contain %q, got:\n%s", substr, result)
				}
			}
		})
	}
}