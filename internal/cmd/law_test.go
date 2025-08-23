package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/api"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/config"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/i18n"
	"github.com/pyhub-kr/pyhub-sejong-cli/internal/testutil"
	"github.com/spf13/cobra"
)

func TestLawCommand(t *testing.T) {
	// Initialize i18n for testing (Korean by default)
	if err := i18n.Init(); err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}
	i18n.SetLanguage("ko")
	
	// Initialize law command
	initLawCmd()
	
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "No arguments",
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg(s), received 0",
		},
		{
			name:        "Empty search query",
			args:        []string{""},
			wantErr:     true,
			errContains: "검색어를 입력해주세요",
		},
		{
			name:        "Multiple arguments",
			args:        []string{"arg1", "arg2"},
			wantErr:     true,
			errContains: "accepts 1 arg(s), received 2",
		},
		{
			name:        "Whitespace only query",
			args:        []string{"   "},
			wantErr:     true,
			errContains: "검색어를 입력해주세요",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(lawCmd)

			// Set args
			args := append([]string{"law"}, tt.args...)
			cmd.SetArgs(args)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Execute command
			err := cmd.Execute()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error should contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestLawCommandFlags(t *testing.T) {
	// Test that flags are properly registered
	if lawCmd.Flag("format") == nil {
		t.Error("format flag not registered")
	}
	if lawCmd.Flag("page") == nil {
		t.Error("page flag not registered")
	}
	if lawCmd.Flag("size") == nil {
		t.Error("size flag not registered")
	}

	// Test default values
	formatFlag := lawCmd.Flag("format")
	if formatFlag.DefValue != "table" {
		t.Errorf("format flag default = %s, want table", formatFlag.DefValue)
	}

	pageFlag := lawCmd.Flag("page")
	if pageFlag.DefValue != "1" {
		t.Errorf("page flag default = %s, want 1", pageFlag.DefValue)
	}

	sizeFlag := lawCmd.Flag("size")
	if sizeFlag.DefValue != "10" {
		t.Errorf("size flag default = %s, want 10", sizeFlag.DefValue)
	}
}

func TestLawCommandWithAPIKey(t *testing.T) {
	// Setup test environment
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-law-test-*")
	defer cleanup()

	// Mock API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check API key
		if r.URL.Query().Get("OC") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return mock response
		response := api.SearchResponse{
			TotalCount: 2,
			Page:       1,
			Laws: []api.LawInfo{
				{
					ID:         "001",
					Name:       "개인정보 보호법",
					NameAbbrev: "개인정보보호법",
					SerialNo:   "12345",
					PromulDate: "20110329",
					PromulNo:   "제10465호",
					Category:   "제정",
					Department: "개인정보보호위원회",
					EffectDate: "20110930",
					LawType:    "법률",
				},
				{
					ID:         "002",
					Name:       "정보통신망 이용촉진 및 정보보호 등에 관한 법률",
					NameAbbrev: "정보통신망법",
					SerialNo:   "67890",
					PromulDate: "20200610",
					PromulNo:   "제17344호",
					Category:   "일부개정",
					Department: "과학기술정보통신부",
					EffectDate: "20201210",
					LawType:    "법률",
				},
			},
		}

		// Check format type
		if r.URL.Query().Get("type") == "JSON" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	// Reset config and set test path
	config.ResetConfig()
	config.SetTestConfigPath(tempDir)

	// Initialize config
	if err := config.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Set test API key
	if err := config.SetAPIKey("test-api-key"); err != nil {
		t.Fatalf("Failed to set API key: %v", err)
	}

	tests := []struct {
		name       string
		args       []string
		format     string
		wantOutput string
	}{
		{
			name:       "Table format",
			args:       []string{"law", "개인정보", "--format", "table"},
			format:     "table",
			wantOutput: "개인정보 보호법",
		},
		{
			name:       "JSON format",
			args:       []string{"law", "개인정보", "--format", "json"},
			format:     "json",
			wantOutput: `"법령명한글": "개인정보 보호법"`,
		},
		{
			name:       "Default format (table)",
			args:       []string{"law", "개인정보"},
			format:     "table",
			wantOutput: "개인정보 보호법",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock API client with test server URL
			testAPIClient = api.NewClientWithURL("test-api-key", server.URL)
			defer func() { testAPIClient = nil }()

			// Create a new root command for testing
			cmd := &cobra.Command{Use: "test"}
			cmd.AddCommand(lawCmd)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			// Check output
			output := buf.String()
			if !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Output should contain %q, got %q", tt.wantOutput, output)
			}
		})
	}
}

func TestLawCommandNoAPIKey(t *testing.T) {
	// Setup test environment without API key
	tempDir, cleanup := testutil.CreateTempDir(t, "sejong-law-test-nokey-*")
	defer cleanup()

	// Reset config and set test path
	config.ResetConfig()
	config.SetTestConfigPath(tempDir)

	// Initialize config without API key
	if err := config.Initialize(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Create a new root command for testing
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(lawCmd)

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Set args
	cmd.SetArgs([]string{"law", "개인정보"})

	// Execute command - should return error about missing API key
	err := cmd.Execute()
	if err == nil {
		t.Error("Execute() should return error when API key is not set")
	}

	// Check that error message contains API key setup instruction
	output := buf.String()
	if !strings.Contains(output, "API 키가 설정되지 않았습니다") {
		t.Errorf("Output should contain API key error message, got %q", output)
	}
}

func TestSearchLaws(t *testing.T) {
	// Test with mock client
	mockClient := &mockAPIClient{
		searchFunc: func(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
			return &api.SearchResponse{
				TotalCount: 1,
				Page:       1,
				Laws: []api.LawInfo{
					{
						ID:         "001",
						Name:       "테스트 법률",
						NameAbbrev: "테스트법",
						SerialNo:   "12345",
						PromulDate: "20240101",
						PromulNo:   "제1호",
						Category:   "제정",
						Department: "테스트부",
						EffectDate: "20240201",
						LawType:    "법률",
					},
				},
			}, nil
		},
	}

	// Test table output
	var buf bytes.Buffer
	err := searchLaws(mockClient, "테스트", "table", 1, 10, &buf, false)
	if err != nil {
		t.Errorf("searchLaws() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "테스트 법률") {
		t.Errorf("Table output should contain law name, got %q", output)
	}

	// Test JSON output
	buf.Reset()
	err = searchLaws(mockClient, "테스트", "json", 1, 10, &buf, false)
	if err != nil {
		t.Errorf("searchLaws() error = %v", err)
	}

	output = buf.String()
	var result api.SearchResponse
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("JSON output should be valid JSON, got %q", output)
	}
	if result.TotalCount != 1 {
		t.Errorf("JSON output TotalCount = %d, want 1", result.TotalCount)
	}
}

// Mock API client for testing
type mockAPIClient struct {
	searchFunc func(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error)
}

func (m *mockAPIClient) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, req)
	}
	return &api.SearchResponse{}, nil
}
