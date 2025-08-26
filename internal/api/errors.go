package api

import (
	"errors"
	"strings"
)

var (
	// ErrNoAPIKey indicates that no API key is configured
	ErrNoAPIKey = errors.New("API 키가 설정되지 않았습니다")

	// ErrNotImplemented indicates that a feature is not implemented
	ErrNotImplemented = errors.New("기능이 구현되지 않았습니다")

	// ErrInvalidAPIType indicates an invalid API type was specified
	ErrInvalidAPIType = errors.New("잘못된 API 타입입니다")
)

// ParseHTMLError extracts meaningful error message from HTML error page
func ParseHTMLError(html string) string {
	// Check for specific error messages in the HTML response
	if strings.Contains(html, "미신청된 목록/본문에 대한 접근입니다") {
		return "API 사용 권한이 없습니다. https://open.law.go.kr 에서 로그인 후 [OPEN API] -> [OPEN API 신청]에서 필요한 법령 종류를 체크해주세요. (도메인은 반드시 '도메인 없음'으로 설정)"
	}

	if strings.Contains(html, "페이지 접속에 실패하였습니다") {
		return "API 접속 실패: 이메일 ID를 확인하거나 서비스 상태를 점검해주세요"
	}

	htmlLower := strings.ToLower(html)

	// Check for authentication/key related issues
	if strings.Contains(htmlLower, "인증") || strings.Contains(htmlLower, "auth") ||
		strings.Contains(htmlLower, "key") || strings.Contains(htmlLower, "키") {
		return "API 인증 실패: 이메일 ID가 올바르지 않습니다. 'warp config set law.key YOUR_EMAIL_ID' 명령으로 이메일 @ 앞부분을 설정하세요"
	}

	// Check for rate limit
	if strings.Contains(htmlLower, "limit") || strings.Contains(htmlLower, "제한") {
		return "API 호출 제한 초과: 일일 호출 한도를 초과했습니다. 잠시 후 다시 시도하세요"
	}

	// Check for server errors
	if strings.Contains(htmlLower, "500") || strings.Contains(htmlLower, "server error") ||
		strings.Contains(htmlLower, "서버 오류") {
		return "서버 오류: 국가법령정보센터 서버에 문제가 발생했습니다. 잠시 후 다시 시도하세요"
	}

	// Check for maintenance
	if strings.Contains(htmlLower, "maintenance") || strings.Contains(htmlLower, "점검") {
		return "서비스 점검 중: 국가법령정보센터가 점검 중입니다. 나중에 다시 시도하세요"
	}

	// Default error message
	return "국가법령정보센터 API에서 오류가 발생했습니다. 이메일 ID를 확인하거나 잠시 후 다시 시도해주세요"
}
