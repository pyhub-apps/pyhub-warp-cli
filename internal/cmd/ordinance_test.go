package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pyhub-apps/sejong-cli/internal/api"
	"github.com/pyhub-apps/sejong-cli/internal/i18n"
	"github.com/spf13/cobra"
)

// MockOrdinanceClient is a mock implementation for testing
type MockOrdinanceClient struct {
	SearchFunc    func(ctx context.Context, req *api.UnifiedSearchRequest) (*api.SearchResponse, error)
	GetDetailFunc func(ctx context.Context, ordinanceID string) (*api.LawDetail, error)
}

func (m *MockOrdinanceClient) Search(ctx context.Context, req *api.UnifiedSearchRequest) (*api.SearchResponse, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, req)
	}
	return &api.SearchResponse{
		TotalCount: 1,
		Page:       1,
		Laws: []api.LawInfo{
			{
				ID:         "ORD001",
				Name:       "서울특별시 주차장 설치 및 관리 조례",
				Department: "서울특별시",
				PromulDate: "20230101",
				LawType:    "자치법규",
			},
		},
	}, nil
}

func (m *MockOrdinanceClient) GetDetail(ctx context.Context, ordinanceID string) (*api.LawDetail, error) {
	if m.GetDetailFunc != nil {
		return m.GetDetailFunc(ctx, ordinanceID)
	}
	return &api.LawDetail{
		LawInfo: api.LawInfo{
			ID:         ordinanceID,
			Name:       "서울특별시 주차장 설치 및 관리 조례",
			Department: "서울특별시",
			PromulDate: "20230101",
			LawType:    "자치법규",
		},
		Content: "조례 본문...",
	}, nil
}

func (m *MockOrdinanceClient) GetHistory(ctx context.Context, ordinanceID string) (*api.LawHistory, error) {
	return nil, nil
}

func (m *MockOrdinanceClient) GetAPIType() api.APIType {
	return api.APITypeELIS
}

func TestOrdinanceCommand(t *testing.T) {
	// Initialize i18n for testing
	i18n.Init()
	i18n.SetLanguage("ko")

	// Set up mock client
	mockClient := &MockOrdinanceClient{}
	testOrdinanceClient = mockClient
	defer func() { testOrdinanceClient = nil }()

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:       "No arguments shows help",
			args:       []string{},
			wantErr:    false,
			wantOutput: "자치법규정보시스템(ELIS)에서",
		},
		{
			name:        "Empty search query",
			args:        []string{""},
			wantErr:     true,
			errContains: "검색어",
		},
		{
			name:       "Valid search query",
			args:       []string{"주차"},
			wantErr:    false,
			wantOutput: "서울특별시 주차장",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test
			cmd := &cobra.Command{
				Use: "test",
			}

			// Initialize ordinance command
			initOrdinanceCmd()
			cmd.AddCommand(ordinanceCmd)

			// Set up output buffer
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			args := append([]string{"ordinance"}, tt.args...)
			cmd.SetArgs(args)

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

			// Check output if expected
			if tt.wantOutput != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantOutput) {
					t.Errorf("Output should contain %q, got %q", tt.wantOutput, output)
				}
			}
		})
	}
}

func TestOrdinanceSearchCommand(t *testing.T) {
	// Initialize i18n for testing
	i18n.Init()
	i18n.SetLanguage("ko")

	// Set up mock client
	mockClient := &MockOrdinanceClient{}
	testOrdinanceClient = mockClient
	defer func() { testOrdinanceClient = nil }()

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:        "No arguments",
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg",
		},
		{
			name:        "Empty search query",
			args:        []string{""},
			wantErr:     true,
			errContains: "검색어",
		},
		{
			name:       "Valid search",
			args:       []string{"주차"},
			wantErr:    false,
			wantOutput: "서울특별시 주차장",
		},
		{
			name:       "Search with region",
			args:       []string{"주차", "--region", "서울"},
			wantErr:    false,
			wantOutput: "서울특별시 주차장",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test
			cmd := &cobra.Command{
				Use: "test",
			}

			// Initialize ordinance command
			initOrdinanceCmd()
			cmd.AddCommand(ordinanceCmd)

			// Set up output buffer
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			args := append([]string{"ordinance", "search"}, tt.args...)
			cmd.SetArgs(args)

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

			// Check output if expected
			if tt.wantOutput != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantOutput) {
					t.Errorf("Output should contain %q, got %q", tt.wantOutput, output)
				}
			}
		})
	}
}

func TestOrdinanceDetailCommand(t *testing.T) {
	// Initialize i18n for testing
	i18n.Init()
	i18n.SetLanguage("ko")

	// Set up mock client
	mockClient := &MockOrdinanceClient{}
	testOrdinanceClient = mockClient
	defer func() { testOrdinanceClient = nil }()

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:        "No arguments",
			args:        []string{},
			wantErr:     true,
			errContains: "accepts 1 arg",
		},
		{
			name:        "Empty ordinance ID",
			args:        []string{""},
			wantErr:     true,
			errContains: "조례 ID",
		},
		{
			name:       "Valid ordinance ID",
			args:       []string{"ORD001"},
			wantErr:    false,
			wantOutput: "서울특별시 주차장",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test
			cmd := &cobra.Command{
				Use: "test",
			}

			// Initialize ordinance command
			initOrdinanceCmd()
			cmd.AddCommand(ordinanceCmd)

			// Set up output buffer
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Set args
			args := append([]string{"ordinance", "detail"}, tt.args...)
			cmd.SetArgs(args)

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

			// Check output if expected
			if tt.wantOutput != "" {
				output := buf.String()
				if !strings.Contains(output, tt.wantOutput) {
					t.Errorf("Output should contain %q, got %q", tt.wantOutput, output)
				}
			}
		})
	}
}
