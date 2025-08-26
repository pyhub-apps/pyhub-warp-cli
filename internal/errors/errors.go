package errors

import (
	"fmt"
)

// ErrorCode represents a specific error category
type ErrorCode string

const (
	// Network errors
	ErrCodeNetwork ErrorCode = "NET001"
	ErrCodeTimeout ErrorCode = "NET002"
	ErrCodeDNS     ErrorCode = "NET003"

	// Authentication errors
	ErrCodeNoAPIKey      ErrorCode = "AUTH001"
	ErrCodeInvalidAPIKey ErrorCode = "AUTH002"
	ErrCodeExpiredAPIKey ErrorCode = "AUTH003"

	// API errors
	ErrCodeAPIResponse ErrorCode = "API001"
	ErrCodeRateLimit   ErrorCode = "API002"
	ErrCodeServerError ErrorCode = "API003"

	// Parsing errors
	ErrCodeJSONParse  ErrorCode = "PARSE001"
	ErrCodeXMLParse   ErrorCode = "PARSE002"
	ErrCodeDataFormat ErrorCode = "PARSE003"

	// Configuration errors
	ErrCodeConfigRead   ErrorCode = "CFG001"
	ErrCodeConfigWrite  ErrorCode = "CFG002"
	ErrCodeConfigFormat ErrorCode = "CFG003"

	// Validation errors
	ErrCodeInvalidInput ErrorCode = "VAL001"
	ErrCodeMissingParam ErrorCode = "VAL002"
)

// CLIError represents a structured error with user-friendly information
type CLIError struct {
	Code       ErrorCode
	Message    string
	Hint       string
	Underlying error
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s\n💡 %s", e.Message, e.Hint)
	}
	return e.Message
}

// DetailedError returns a detailed error message including the code
func (e *CLIError) DetailedError() string {
	msg := fmt.Sprintf("[%s] %s", e.Code, e.Message)
	if e.Hint != "" {
		msg += fmt.Sprintf("\n💡 힌트: %s", e.Hint)
	}
	if e.Underlying != nil {
		msg += fmt.Sprintf("\n🔍 상세: %v", e.Underlying)
	}
	return msg
}

// Unwrap returns the underlying error for errors.Is and errors.As support
func (e *CLIError) Unwrap() error {
	return e.Underlying
}

// Common error definitions
var (
	// Network errors
	ErrNoNetwork = &CLIError{
		Code:    ErrCodeNetwork,
		Message: "서버에 연결할 수 없습니다",
		Hint:    "인터넷 연결을 확인하세요",
	}

	ErrTimeout = &CLIError{
		Code:    ErrCodeTimeout,
		Message: "요청 시간이 초과되었습니다",
		Hint:    "네트워크 상태를 확인하거나 잠시 후 다시 시도하세요",
	}

	// Authentication errors
	ErrNoAPIKey = &CLIError{
		Code:    ErrCodeNoAPIKey,
		Message: "API 키가 설정되지 않았습니다",
		Hint:    "warp config set law.key <YOUR_KEY> 명령으로 API 키를 설정하세요",
	}

	ErrInvalidAPIKey = &CLIError{
		Code:    ErrCodeInvalidAPIKey,
		Message: "API 인증에 실패했습니다",
		Hint:    "API 키가 올바른지 확인하세요: warp config get law.key",
	}

	// API errors
	ErrAPIServerError = &CLIError{
		Code:    ErrCodeServerError,
		Message: "API 서버에서 오류가 발생했습니다",
		Hint:    "잠시 후 다시 시도하거나 서비스 상태를 확인하세요",
	}

	ErrRateLimit = &CLIError{
		Code:    ErrCodeRateLimit,
		Message: "API 요청 한도를 초과했습니다",
		Hint:    "잠시 후 다시 시도하세요",
	}

	// Parsing errors
	ErrJSONParse = &CLIError{
		Code:    ErrCodeJSONParse,
		Message: "응답 데이터를 파싱할 수 없습니다",
		Hint:    "API 응답 형식이 변경되었을 수 있습니다. 최신 버전으로 업데이트하세요",
	}

	// Validation errors
	ErrEmptyQuery = &CLIError{
		Code:    ErrCodeInvalidInput,
		Message: "검색어를 입력해주세요",
		Hint:    "예: warp law \"개인정보 보호법\"",
	}
)

// New creates a new CLIError with the given parameters
func New(code ErrorCode, message, hint string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Hint:    hint,
	}
}

// Wrap wraps an existing error with CLI error context
func Wrap(err error, cliErr *CLIError) *CLIError {
	return &CLIError{
		Code:       cliErr.Code,
		Message:    cliErr.Message,
		Hint:       cliErr.Hint,
		Underlying: err,
	}
}

// WithHint adds or updates the hint for an error
func WithHint(err *CLIError, hint string) *CLIError {
	return &CLIError{
		Code:       err.Code,
		Message:    err.Message,
		Hint:       hint,
		Underlying: err.Underlying,
	}
}
