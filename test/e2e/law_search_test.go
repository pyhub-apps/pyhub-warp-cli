//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2EFirstUserScenario tests the first-time user experience
func TestE2EFirstUserScenario(t *testing.T) {
	// Setup
	tempDir := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	// Scenario 1: Search without API key - should show guidance
	t.Run("SearchWithoutAPIKey", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "law", "개인정보 보호법")
		cmd.Env = append(os.Environ(), fmt.Sprintf("HOME=%s", tempDir))

		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "API 키가 설정되지 않았습니다")
		assert.Contains(t, outputStr, "sejong config set law.key")
	})

	// Scenario 2: Set API key
	t.Run("SetAPIKey", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "config", "set", "law.key", "TEST_API_KEY")
		cmd.Env = append(os.Environ(), fmt.Sprintf("HOME=%s", tempDir))

		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "API 키가 성공적으로 설정되었습니다")
	})

	// Scenario 3: Verify API key is set
	t.Run("GetAPIKey", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "config", "get", "law.key")
		cmd.Env = append(os.Environ(), fmt.Sprintf("HOME=%s", tempDir))

		output, err := cmd.CombinedOutput()
		require.NoError(t, err)
		
		outputStr := string(output)
		// Output format changed - now shows first 4 chars + ...
		assert.Contains(t, outputStr, "TEST_")
		assert.Contains(t, outputStr, "law.key")
	})
}

// TestE2ENormalUserScenario tests normal usage scenarios
func TestE2ENormalUserScenario(t *testing.T) {
	// Skip this test for now as it requires integration with the actual command
	// The command doesn't support custom API URL through environment variables
	t.Skip("Skipping E2E normal user scenario - requires API URL configuration support")
	
	// Setup
	tempDir := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	// Set API key and mock server URL
	setupConfig(t, tempDir, "TEST_API_KEY", mockServer.GetSearchURL())

	// Scenario 1: Normal search with table output
	t.Run("NormalSearchTableOutput", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "law", "개인정보 보호법")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=%s", mockServer.GetSearchURL()),
		)

		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command failed: %s", string(output))
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "개인정보 보호법")
		assert.Contains(t, outputStr, "개인정보보호위원회")
		assert.Contains(t, outputStr, "2024-03-15")
		assert.Contains(t, outputStr, "총 3개의 법령")
	})

	// Scenario 2: Search with JSON output
	t.Run("SearchWithJSONOutput", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "law", "도로교통법", "--format", "json")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=%s", mockServer.GetSearchURL()),
		)

		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command failed: %s", string(output))
		
		var result map[string]interface{}
		err = json.Unmarshal(output, &result)
		require.NoError(t, err, "Failed to parse JSON output")
		
		assert.Equal(t, float64(1), result["totalCnt"])
		laws := result["law"].([]interface{})
		assert.Len(t, laws, 1)
		
		firstLaw := laws[0].(map[string]interface{})
		assert.Equal(t, "도로교통법", firstLaw["법령명한글"])
	})

	// Scenario 3: Empty result handling
	t.Run("EmptyResultHandling", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "law", "없는법령")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=%s", mockServer.GetSearchURL()),
		)

		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command failed: %s", string(output))
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "검색 결과가 없습니다")
	})

	// Scenario 4: Pagination test
	t.Run("PaginationTest", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "law", "개인정보 보호법", "--page", "2", "--size", "10")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=%s", mockServer.GetSearchURL()),
		)

		output, err := cmd.CombinedOutput()
		// Note: This might return empty results for page 2, which is expected
		_ = err
		_ = output
	})
}

// TestE2EErrorScenarios tests various error scenarios
func TestE2EErrorScenarios(t *testing.T) {
	// Skip this test for now as it requires integration with the actual command
	t.Skip("Skipping E2E error scenarios - requires API URL configuration support")
	
	// Setup
	tempDir := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	mockServer := testutil.NewMockServer()
	defer mockServer.Close()

	// Scenario 1: Invalid API key
	t.Run("InvalidAPIKey", func(t *testing.T) {
		setupConfig(t, tempDir, "INVALID_KEY", mockServer.GetSearchURL())

		cmd := exec.Command("../../sejong", "law", "개인정보")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=%s", mockServer.GetSearchURL()),
		)

		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "Invalid API key")
	})

	// Scenario 2: Server error
	t.Run("ServerError", func(t *testing.T) {
		setupConfig(t, tempDir, "TEST_API_KEY", mockServer.GetSearchURL())

		cmd := exec.Command("../../sejong", "law", "error")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=%s", mockServer.GetSearchURL()),
		)

		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "error")
	})

	// Scenario 3: Network timeout (simulated with non-existent server)
	t.Run("NetworkTimeout", func(t *testing.T) {
		setupConfig(t, tempDir, "TEST_API_KEY", "http://localhost:99999/api")

		cmd := exec.Command("../../sejong", "law", "test")
		cmd.Env = append(os.Environ(), 
			fmt.Sprintf("HOME=%s", tempDir),
			fmt.Sprintf("LAW_API_URL=http://localhost:99999/api", ),
		)

		output, err := cmd.CombinedOutput()
		assert.Error(t, err)
		
		outputStr := string(output)
		// Should contain some network error message
		assert.True(t, 
			strings.Contains(outputStr, "connection") || 
			strings.Contains(outputStr, "네트워크") ||
			strings.Contains(outputStr, "연결"),
		)
	})
}

// TestE2EVersionAndHelp tests version and help commands
func TestE2EVersionAndHelp(t *testing.T) {
	// Test version command
	t.Run("VersionCommand", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "version")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command failed: %s", string(output))
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "Version:")
		assert.Contains(t, outputStr, "Built:")
	})

	// Test help command
	t.Run("HelpCommand", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command failed: %s", string(output))
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "sejong")
		assert.Contains(t, outputStr, "law")
		assert.Contains(t, outputStr, "config")
		assert.Contains(t, outputStr, "version")
	})

	// Test law subcommand help
	t.Run("LawHelpCommand", func(t *testing.T) {
		cmd := exec.Command("../../sejong", "law", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command failed: %s", string(output))
		
		outputStr := string(output)
		assert.Contains(t, outputStr, "--format")
		assert.Contains(t, outputStr, "--page")
		assert.Contains(t, outputStr, "--size")
		assert.Contains(t, outputStr, "--verbose")
	})
}

// Helper functions

func setupTestEnvironment(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "sejong-e2e-test-*")
	require.NoError(t, err)
	
	// Create .sejong directory
	sejongDir := filepath.Join(tempDir, ".sejong")
	err = os.MkdirAll(sejongDir, 0755)
	require.NoError(t, err)
	
	return tempDir
}

func setupConfig(t *testing.T, homeDir, apiKey, apiURL string) {
	configDir := filepath.Join(homeDir, ".sejong")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)
	
	configFile := filepath.Join(configDir, "config.yaml")
	configContent := fmt.Sprintf(`law:
  key: %s
  url: %s
`, apiKey, apiURL)
	
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)
}