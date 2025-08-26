package onboarding

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/fatih/color"
)

// Guide provides user onboarding assistance
type Guide struct {
	writer   io.Writer
	useColor bool
}

// NewGuide creates a new onboarding guide
func NewGuide() *Guide {
	return &Guide{
		writer:   os.Stderr,
		useColor: isTerminal() && !isColorDisabled(),
	}
}

// NewGuideWithWriter creates a guide with custom writer (for testing)
func NewGuideWithWriter(w io.Writer, useColor bool) *Guide {
	return &Guide{
		writer:   w,
		useColor: useColor,
	}
}

// ShowAPIKeySetup displays the API key setup guide
func (g *Guide) ShowAPIKeySetup() {
	if g.useColor {
		g.showColoredAPIKeySetup()
	} else {
		g.showPlainAPIKeySetup()
	}
}

func (g *Guide) showColoredAPIKeySetup() {
	// Force color output even when not a terminal (for testing)
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)
	bold := color.New(color.Bold)

	// Ensure color output is used
	if g.useColor {
		color.NoColor = false
	}

	// Header
	red.Fprintln(g.writer, "🔐 API 설정이 필요합니다")
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "국가법령정보센터 오픈 API를 사용하려면 이메일 인증이 필요합니다.")
	fmt.Fprintln(g.writer)

	// Steps
	bold.Fprintln(g.writer, "📋 설정 방법:")
	fmt.Fprintln(g.writer)

	// Step 1
	yellow.Fprint(g.writer, "1️⃣  Open API 신청하기")
	fmt.Fprintln(g.writer)
	fmt.Fprint(g.writer, "   → ")
	cyan.Fprintln(g.writer, "https://open.law.go.kr")
	fmt.Fprintln(g.writer, "   • 회원가입 및 로그인")
	fmt.Fprintln(g.writer, "   • [OPEN API] → [OPEN API 신청] 메뉴")
	fmt.Fprintln(g.writer, "   • 필요한 법령 종류 체크 (법령, 판례, 행정규칙 등)")
	fmt.Fprint(g.writer, "   ")
	red.Fprintln(g.writer, "⚠️  중요: 도메인 주소는 반드시 \"도메인 없음\"으로 설정")
	fmt.Fprintln(g.writer)

	// Step 2
	yellow.Fprint(g.writer, "2️⃣  이메일 ID 설정하기")
	fmt.Fprintln(g.writer)
	fmt.Fprint(g.writer, "   → ")
	green.Fprintln(g.writer, "warp config set law.key <이메일ID>")
	fmt.Fprintln(g.writer, "   예: example@gmail.com → example")
	fmt.Fprintln(g.writer)

	// Tip
	fmt.Fprint(g.writer, "💡 ")
	bold.Fprintln(g.writer, "팁: 위 명령어를 복사하여 사용하세요!")

	// Platform-specific copy hint
	g.showCopyHint()
}

func (g *Guide) showPlainAPIKeySetup() {
	fmt.Fprintln(g.writer, "❌ API 설정이 필요합니다")
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "국가법령정보센터 오픈 API를 사용하려면 이메일 인증이 필요합니다.")
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "📋 설정 방법:")
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "1. Open API 신청하기")
	fmt.Fprintln(g.writer, "   → https://open.law.go.kr")
	fmt.Fprintln(g.writer, "   • 회원가입 및 로그인")
	fmt.Fprintln(g.writer, "   • [OPEN API] → [OPEN API 신청] 메뉴")
	fmt.Fprintln(g.writer, "   • 필요한 법령 종류 체크 (법령, 판례, 행정규칙 등)")
	fmt.Fprintln(g.writer, "   ⚠️  중요: 도메인 주소는 반드시 \"도메인 없음\"으로 설정")
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "2. 이메일 ID 설정하기")
	fmt.Fprintln(g.writer, "   → warp config set law.key <이메일ID>")
	fmt.Fprintln(g.writer, "   예: example@gmail.com → example")
	fmt.Fprintln(g.writer)
	fmt.Fprintln(g.writer, "💡 팁: 위 명령어를 복사하여 사용하세요!")

	g.showCopyHint()
}

func (g *Guide) showCopyHint() {
	switch runtime.GOOS {
	case "darwin":
		fmt.Fprintln(g.writer, "   (Mac: Cmd+C로 복사)")
	case "windows":
		fmt.Fprintln(g.writer, "   (Windows: Ctrl+C로 복사 또는 마우스 우클릭)")
	default:
		fmt.Fprintln(g.writer, "   (Linux: Ctrl+Shift+C로 복사)")
	}
}

// ShowSearchProgress displays a search in progress message
func (g *Guide) ShowSearchProgress(query string) {
	if g.useColor {
		spinner := color.New(color.FgCyan)
		spinner.Fprintf(g.writer, "🔍 검색 중... (%s)\n", query)
	} else {
		fmt.Fprintf(g.writer, "검색 중... (%s)\n", query)
	}
}

// ShowSuccess displays a success message
func (g *Guide) ShowSuccess(message string) {
	if g.useColor {
		green := color.New(color.FgGreen, color.Bold)
		green.Fprintf(g.writer, "✅ %s\n", message)
	} else {
		fmt.Fprintf(g.writer, "✅ %s\n", message)
	}
}

// ShowError displays an error message
func (g *Guide) ShowError(message string) {
	if g.useColor {
		red := color.New(color.FgRed, color.Bold)
		red.Fprintf(g.writer, "❌ %s\n", message)
	} else {
		fmt.Fprintf(g.writer, "❌ %s\n", message)
	}
}

// ShowWarning displays a warning message
func (g *Guide) ShowWarning(message string) {
	if g.useColor {
		yellow := color.New(color.FgYellow)
		yellow.Fprintf(g.writer, "⚠️  %s\n", message)
	} else {
		fmt.Fprintf(g.writer, "! %s\n", message)
	}
}

// isTerminal checks if output is a terminal
func isTerminal() bool {
	fileInfo, err := os.Stderr.Stat()
	if err != nil {
		// If we can't stat stderr, assume it's not a terminal
		// This is safe because we'll just disable colors
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// isColorDisabled checks if color output is disabled via environment
func isColorDisabled() bool {
	// Check common environment variables that disable color
	if os.Getenv("NO_COLOR") != "" {
		return true
	}
	if os.Getenv("TERM") == "dumb" {
		return true
	}
	return false
}
