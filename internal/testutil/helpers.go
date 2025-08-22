package testutil

import (
	"strings"
	"testing"
)

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, str, substr string, msgAndArgs ...interface{}) {
	t.Helper()
	if !strings.Contains(str, substr) {
		if len(msgAndArgs) > 0 {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			t.Errorf("Expected string to contain %q, got %q", substr, str)
		}
	}
}

// AssertNotContains checks if a string does not contain a substring
func AssertNotContains(t *testing.T, str, substr string, msgAndArgs ...interface{}) {
	t.Helper()
	if strings.Contains(str, substr) {
		if len(msgAndArgs) > 0 {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			t.Errorf("Expected string to not contain %q, got %q", substr, str)
		}
	}
}

// AssertError checks if an error occurred when expected
func AssertError(t *testing.T, err error, expectError bool, msgAndArgs ...interface{}) {
	t.Helper()
	if (err != nil) != expectError {
		if len(msgAndArgs) > 0 {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			t.Errorf("Error occurrence mismatch: got error=%v, expected error=%v", err != nil, expectError)
		}
	}
}

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, actual, expected interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if actual != expected {
		if len(msgAndArgs) > 0 {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			t.Errorf("Values not equal: got %v, expected %v", actual, expected)
		}
	}
}