package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		expected string
	}{
		{
			name: "Error with hint",
			err: &CLIError{
				Code:    ErrCodeNetwork,
				Message: "Connection failed",
				Hint:    "Check your internet",
			},
			expected: "Connection failed\n💡 Check your internet",
		},
		{
			name: "Error without hint",
			err: &CLIError{
				Code:    ErrCodeAPIResponse,
				Message: "Bad response",
			},
			expected: "Bad response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCLIError_DetailedError(t *testing.T) {
	t.Run("With underlying error", func(t *testing.T) {
		underlying := errors.New("socket timeout")
		err := &CLIError{
			Code:       ErrCodeTimeout,
			Message:    "Request timeout",
			Hint:       "Try again later",
			Underlying: underlying,
		}

		detailed := err.DetailedError()
		
		// Check for all expected parts
		if !strings.Contains(detailed, "[NET002]") {
			t.Error("DetailedError should contain error code")
		}
		if !strings.Contains(detailed, "Request timeout") {
			t.Error("DetailedError should contain message")
		}
		if !strings.Contains(detailed, "힌트: Try again later") {
			t.Error("DetailedError should contain hint")
		}
		if !strings.Contains(detailed, "socket timeout") {
			t.Error("DetailedError should contain underlying error")
		}
	})

	t.Run("Without underlying error", func(t *testing.T) {
		err := &CLIError{
			Code:    ErrCodeInvalidInput,
			Message: "Invalid input",
			Hint:    "Check your input",
		}

		detailed := err.DetailedError()
		
		// Check that it contains expected parts
		if !strings.Contains(detailed, "[VAL001]") {
			t.Error("DetailedError should contain error code")
		}
		if !strings.Contains(detailed, "Invalid input") {
			t.Error("DetailedError should contain message")
		}
		if !strings.Contains(detailed, "힌트: Check your input") {
			t.Error("DetailedError should contain hint")
		}
		// Should NOT contain "상세" field when no underlying error
		if strings.Contains(detailed, "🔍 상세:") {
			t.Error("DetailedError should not contain 상세 field when underlying is nil")
		}
	})
}

func TestWrap(t *testing.T) {
	underlying := errors.New("connection refused")
	wrapped := Wrap(underlying, ErrNoNetwork)
	
	if wrapped.Underlying != underlying {
		t.Error("Wrap should preserve underlying error")
	}
	if wrapped.Code != ErrNoNetwork.Code {
		t.Error("Wrap should preserve error code")
	}
	if wrapped.Message != ErrNoNetwork.Message {
		t.Error("Wrap should preserve message")
	}
	
	// Test unwrap
	if errors.Unwrap(wrapped) != underlying {
		t.Error("Unwrap should return underlying error")
	}
}

func TestWithHint(t *testing.T) {
	original := &CLIError{
		Code:    ErrCodeInvalidInput,
		Message: "Invalid input",
		Hint:    "Original hint",
	}
	
	newHint := "New hint"
	modified := WithHint(original, newHint)
	
	if modified.Hint != newHint {
		t.Errorf("WithHint should update hint, got %q, want %q", modified.Hint, newHint)
	}
	if modified.Code != original.Code {
		t.Error("WithHint should preserve error code")
	}
	if modified.Message != original.Message {
		t.Error("WithHint should preserve message")
	}
}

func TestCommonErrors(t *testing.T) {
	// Test that common errors are properly defined
	commonErrors := []*CLIError{
		ErrNoNetwork,
		ErrTimeout,
		ErrNoAPIKey,
		ErrInvalidAPIKey,
		ErrAPIServerError,
		ErrRateLimit,
		ErrJSONParse,
		ErrEmptyQuery,
	}
	
	for _, err := range commonErrors {
		if err.Code == "" {
			t.Errorf("Error %v should have a code", err)
		}
		if err.Message == "" {
			t.Errorf("Error %v should have a message", err)
		}
		// Most errors should have hints
		if err.Hint == "" && err != ErrAPIServerError {
			t.Logf("Warning: Error %v has no hint", err)
		}
	}
}

func TestNew(t *testing.T) {
	err := New(ErrCodeConfigRead, "Cannot read config", "Check file permissions")
	
	if err.Code != ErrCodeConfigRead {
		t.Error("New should set error code")
	}
	if err.Message != "Cannot read config" {
		t.Error("New should set message")
	}
	if err.Hint != "Check file permissions" {
		t.Error("New should set hint")
	}
}