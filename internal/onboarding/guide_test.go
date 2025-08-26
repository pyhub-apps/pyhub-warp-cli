package onboarding

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewGuide(t *testing.T) {
	guide := NewGuide()
	if guide == nil {
		t.Error("Expected guide to be created")
	}
	if guide.writer == nil {
		t.Error("Expected writer to be set")
	}
}

func TestGuide_ShowAPIKeySetup_Plain(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, false)

	guide.ShowAPIKeySetup()

	output := buf.String()

	// Check for required elements
	expectedStrings := []string{
		"API 설정이 필요합니다",
		"국가법령정보센터 오픈 API",
		"설정 방법:",
		"Open API 신청하기",
		"https://open.law.go.kr",
		"도메인 없음",
		"이메일 ID 설정하기",
		"warp config set law.key",
		"팁: 위 명령어를 복사하여 사용하세요",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output should contain %q", expected)
		}
	}
}

func TestGuide_ShowAPIKeySetup_Colored(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, true)

	guide.ShowAPIKeySetup()

	output := buf.String()

	// Check for required content (color codes will be present but we check the text)
	expectedStrings := []string{
		"API 설정이 필요합니다",
		"설정 방법:",
		"Open API 신청하기",
		"도메인 없음",
		"이메일 ID 설정하기",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output should contain %q", expected)
		}
	}

	// Check that ANSI color codes are present
	if !strings.Contains(output, "\x1b[") {
		t.Error("Colored output should contain ANSI escape codes")
	}
}

func TestGuide_ShowSearchProgress(t *testing.T) {
	tests := []struct {
		name     string
		useColor bool
		query    string
		expected string
	}{
		{
			name:     "Plain text",
			useColor: false,
			query:    "test query",
			expected: "검색 중... (test query)",
		},
		{
			name:     "Colored text",
			useColor: true,
			query:    "test query",
			expected: "검색 중... (test query)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			guide := NewGuideWithWriter(&buf, tt.useColor)

			guide.ShowSearchProgress(tt.query)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Output should contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestGuide_ShowSuccess(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, false)

	guide.ShowSuccess("Operation completed")

	output := buf.String()
	if !strings.Contains(output, "Operation completed") {
		t.Errorf("Output should contain success message")
	}
	if !strings.Contains(output, "✅") {
		t.Errorf("Output should contain success indicator")
	}
}

func TestGuide_ShowError(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, false)

	guide.ShowError("Operation failed")

	output := buf.String()
	if !strings.Contains(output, "Operation failed") {
		t.Errorf("Output should contain error message")
	}
	if !strings.Contains(output, "❌") {
		t.Errorf("Output should contain error indicator")
	}
}

func TestGuide_ShowError_Colored(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, true)

	guide.ShowError("Operation failed")

	output := buf.String()
	if !strings.Contains(output, "Operation failed") {
		t.Errorf("Output should contain error message")
	}
	if !strings.Contains(output, "❌") {
		t.Errorf("Output should contain error indicator")
	}
	// Check that ANSI color codes are present
	if !strings.Contains(output, "\x1b[") {
		t.Error("Colored output should contain ANSI escape codes")
	}
}

func TestGuide_ShowWarning(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, false)

	guide.ShowWarning("This is a warning")

	output := buf.String()
	if !strings.Contains(output, "This is a warning") {
		t.Errorf("Output should contain warning message")
	}
	if !strings.Contains(output, "!") {
		t.Errorf("Output should contain warning indicator")
	}
}

func TestGuide_PlatformSpecificCopyHint(t *testing.T) {
	var buf bytes.Buffer
	guide := NewGuideWithWriter(&buf, false)

	guide.ShowAPIKeySetup()

	output := buf.String()

	// Should contain at least one of the platform-specific copy hints
	hasCopyHint := strings.Contains(output, "Cmd+C") ||
		strings.Contains(output, "Ctrl+C") ||
		strings.Contains(output, "Ctrl+Shift+C")

	if !hasCopyHint {
		t.Error("Output should contain platform-specific copy hint")
	}
}
