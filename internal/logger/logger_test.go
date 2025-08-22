package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel, &buf, false)
	
	// Debug should not be logged at INFO level
	logger.Debug("debug message")
	if buf.String() != "" {
		t.Error("Debug message should not be logged at INFO level")
	}
	
	// Info should be logged
	buf.Reset()
	logger.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message should be logged at INFO level")
	}
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Error("Info message should contain [INFO] tag")
	}
	
	// Warn should be logged
	buf.Reset()
	logger.Warn("warning message")
	if !strings.Contains(buf.String(), "warning message") {
		t.Error("Warning message should be logged at INFO level")
	}
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Error("Warning message should contain [WARN] tag")
	}
	
	// Error should be logged
	buf.Reset()
	logger.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error message should be logged at INFO level")
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Error("Error message should contain [ERROR] tag")
	}
}

func TestLogger_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DebugLevel, &buf, false)
	
	// All messages should be logged at DEBUG level
	logger.Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should be logged at DEBUG level")
	}
	
	buf.Reset()
	logger.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message should be logged at DEBUG level")
	}
}

func TestLogger_ErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(ErrorLevel, &buf, false)
	
	// Debug and Info should not be logged
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	if buf.String() != "" {
		t.Error("Debug/Info/Warn messages should not be logged at ERROR level")
	}
	
	// Error should be logged
	buf.Reset()
	logger.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error message should be logged at ERROR level")
	}
}

func TestLogger_Formatting(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel, &buf, false)
	
	logger.Info("test %s %d", "string", 42)
	output := buf.String()
	
	if !strings.Contains(output, "test string 42") {
		t.Error("Logger should support printf-style formatting")
	}
	
	// Check timestamp format (HH:MM:SS)
	if !strings.Contains(output, ":") {
		t.Error("Logger output should contain timestamp")
	}
}

func TestLogger_ColorOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel, &buf, true)
	
	logger.Info("colored message")
	output := buf.String()
	
	// When color is enabled, ANSI codes should be present
	// The actual presence of ANSI codes depends on the terminal detection
	// but we can check that the message is still there
	if !strings.Contains(output, "colored message") {
		t.Error("Colored logger should still output the message")
	}
}

func TestSetVerbose(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(InfoLevel)
	
	// Debug should not be logged initially
	Debug("debug message")
	if buf.String() != "" {
		t.Error("Debug should not be logged at INFO level")
	}
	
	// Enable verbose mode
	buf.Reset()
	SetVerbose(true)
	Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug should be logged when verbose is enabled")
	}
	
	// Disable verbose mode
	buf.Reset()
	SetVerbose(false)
	Debug("debug message")
	if buf.String() != "" {
		t.Error("Debug should not be logged when verbose is disabled")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"info", InfoLevel},
		{"INFO", InfoLevel},
		{"warn", WarnLevel},
		{"WARN", WarnLevel},
		{"warning", WarnLevel},
		{"WARNING", WarnLevel},
		{"error", ErrorLevel},
		{"ERROR", ErrorLevel},
		{"fatal", FatalLevel},
		{"FATAL", FatalLevel},
		{"unknown", InfoLevel}, // Default to INFO
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLogError(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(InfoLevel)
	
	// Test with nil error
	LogError(nil, false)
	if buf.String() != "" {
		t.Error("LogError should not log nil errors")
	}
	
	// Test with non-verbose mode
	buf.Reset()
	err := &testError{msg: "test error"}
	LogError(err, false)
	if !strings.Contains(buf.String(), "test error") {
		t.Error("LogError should log error message")
	}
	
	// Test with verbose mode
	buf.Reset()
	LogError(err, true)
	output := buf.String()
	if !strings.Contains(output, "Error occurred:") {
		t.Error("LogError in verbose mode should include 'Error occurred:' prefix")
	}
	if !strings.Contains(output, "test error") {
		t.Error("LogError in verbose mode should include error message")
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}