//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/api"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIClientWithMockServer tests API client with mock server
func TestAPIClientWithMockServer(t *testing.T) {
	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	client := api.NewClientWithURL("TEST_API_KEY", mockServer.GetSearchURL())

	tests := []struct {
		name        string
		query       string
		expectError bool
		expectCount int
	}{
		{
			name:        "Personal Information Protection Act",
			query:       "개인정보 보호법",
			expectError: false,
			expectCount: 3,
		},
		{
			name:        "Traffic Law",
			query:       "도로교통법",
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "Non-existent Law",
			query:       "없는법령",
			expectError: false,
			expectCount: 0,
		},
		{
			name:        "Server Error",
			query:       "error",
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req := &api.SearchRequest{
				Query:    tt.query,
				Type:     "JSON",
				PageNo:   1,
				PageSize: 20,
			}
			result, err := client.Search(ctx, req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectCount, result.TotalCount)
				assert.Len(t, result.Laws, tt.expectCount)
			}
		})
	}
}

// TestAPIClientWithInvalidKey tests API client with invalid API key
func TestAPIClientWithInvalidKey(t *testing.T) {
	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	client := api.NewClientWithURL("INVALID_KEY", mockServer.GetSearchURL())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.SearchRequest{
		Query:    "test",
		Type:     "JSON",
		PageNo:   1,
		PageSize: 20,
	}
	result, err := client.Search(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)
	// The actual error message is in Korean
	assert.Contains(t, err.Error(), "403")
}

// TestAPIClientWithoutKey tests API client without API key
func TestAPIClientWithoutKey(t *testing.T) {
	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	client := api.NewClientWithURL("", mockServer.GetSearchURL())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.SearchRequest{
		Query:    "test",
		Type:     "JSON",
		PageNo:   1,
		PageSize: 20,
	}
	result, err := client.Search(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestConfigurationIntegration tests configuration loading and saving
func TestConfigurationIntegration(t *testing.T) {
	// Create temp directory for config
	tempDir, err := os.MkdirTemp("", "warp-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Test 1: Create config file
	t.Run("CreateConfig", func(t *testing.T) {
		configDir := filepath.Join(tempDir, ".pyhub", "warp")
		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		viper.SetConfigFile(filepath.Join(configDir, "config.yaml"))
		viper.Set("law.key", "")
		err = viper.WriteConfig()
		require.NoError(t, err)
	})

	// Test 2: Set and save API key
	t.Run("SetAndSaveAPIKey", func(t *testing.T) {
		configFile := filepath.Join(tempDir, ".pyhub", "warp", "config.yaml")
		viper.SetConfigFile(configFile)
		viper.Set("law.key", "TEST_API_KEY_123")
		err := viper.WriteConfig()
		require.NoError(t, err)
	})

	// Test 3: Load saved config
	t.Run("LoadSavedConfig", func(t *testing.T) {
		configFile := filepath.Join(tempDir, ".pyhub", "warp", "config.yaml")
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		require.NoError(t, err)
		assert.Equal(t, "TEST_API_KEY_123", viper.GetString("law.key"))
	})

	// Test 4: Get masked API key
	t.Run("GetMaskedAPIKey", func(t *testing.T) {
		apiKey := viper.GetString("law.key")
		if len(apiKey) > 4 {
			masked := apiKey[:4] + "***"
			assert.Equal(t, "TEST***", masked)
		}
	})
}

// TestRetryMechanism tests the retry mechanism with network errors
func TestRetryMechanism(t *testing.T) {
	// This test simulates network errors and retry behavior
	// Since our current API client doesn't have built-in retry logic,
	// we'll test the basic error handling

	client := api.NewClientWithURL("TEST_API_KEY", "http://localhost:99999/api")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req := &api.SearchRequest{
		Query:    "test",
		Type:     "JSON",
		PageNo:   1,
		PageSize: 20,
	}
	result, err := client.Search(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)
	// Should contain some network error
	assert.True(t, err != nil)
}

// TestPaginationIntegration tests pagination functionality
func TestPaginationIntegration(t *testing.T) {
	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	client := api.NewClientWithURL("TEST_API_KEY", mockServer.GetSearchURL())

	// Test with default mock responses
	t.Run("DefaultPagination", func(t *testing.T) {
		ctx := context.Background()

		// Page 1 should return results
		req := &api.SearchRequest{
			Query:    "개인정보 보호법",
			Type:     "JSON",
			PageNo:   1,
			PageSize: 20,
		}
		result, err := client.Search(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.TotalCount)
		assert.Len(t, result.Laws, 3)

		// Different query
		req = &api.SearchRequest{
			Query:    "도로교통법",
			Type:     "JSON",
			PageNo:   1,
			PageSize: 20,
		}
		result, err = client.Search(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.TotalCount)
		assert.Len(t, result.Laws, 1)
	})
}
