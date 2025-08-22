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