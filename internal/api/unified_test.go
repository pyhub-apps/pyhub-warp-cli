package api

import (
	"context"
	"testing"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/config"
)

func TestUnifiedClient_Search(t *testing.T) {
	// This test requires API key to be set
	apiKey := config.GetNLICAPIKey()
	if apiKey == "" {
		t.Skip("API key not set, skipping unified search test")
	}

	client, err := NewUnifiedClient()
	if err != nil {
		t.Fatalf("Failed to create unified client: %v", err)
	}

	ctx := context.Background()
	req := &UnifiedSearchRequest{
		Query:    "환경",
		PageNo:   1,
		PageSize: 5,
		Type:     "json",
	}

	// Test unified search
	result, err := client.Search(ctx, req)
	if err != nil {
		t.Errorf("Unified search failed: %v", err)
		return
	}

	if result == nil {
		t.Error("Expected non-nil result")
		return
	}

	// Check that we have results
	if len(result.Laws) == 0 {
		t.Error("Expected at least one result")
		return
	}

	// Check that results have source information
	hasNLIC := false
	hasELIS := false
	for _, law := range result.Laws {
		if law.Source == "국가법령" {
			hasNLIC = true
		}
		if law.Source == "자치법규" {
			hasELIS = true
		}
	}

	t.Logf("Unified search returned %d results (NLIC: %v, ELIS: %v)",
		len(result.Laws), hasNLIC, hasELIS)
}

func TestUnifiedClient_SearchWithOptions(t *testing.T) {
	// This test requires API key to be set
	apiKey := config.GetNLICAPIKey()
	if apiKey == "" {
		t.Skip("API key not set, skipping unified search test")
	}

	client, err := NewUnifiedClient()
	if err != nil {
		t.Fatalf("Failed to create unified client: %v", err)
	}

	ctx := context.Background()
	req := &UnifiedSearchRequest{
		Query:    "환경",
		PageNo:   1,
		PageSize: 5,
		Type:     "json",
	}

	tests := []struct {
		name        string
		includeNLIC bool
		includeELIS bool
		wantErr     bool
	}{
		{
			name:        "Both APIs",
			includeNLIC: true,
			includeELIS: true,
			wantErr:     false,
		},
		{
			name:        "NLIC only",
			includeNLIC: true,
			includeELIS: false,
			wantErr:     false,
		},
		{
			name:        "ELIS only",
			includeNLIC: false,
			includeELIS: true,
			wantErr:     false,
		},
		{
			name:        "Neither API",
			includeNLIC: false,
			includeELIS: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.SearchWithOptions(ctx, req, tt.includeNLIC, tt.includeELIS)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchWithOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != nil {
				t.Logf("%s returned %d results", tt.name, len(result.Laws))
			}
		})
	}
}

func TestUnifiedClient_GetDetail(t *testing.T) {
	// This test requires API key to be set
	apiKey := config.GetNLICAPIKey()
	if apiKey == "" {
		t.Skip("API key not set, skipping unified detail test")
	}

	client, err := NewUnifiedClient()
	if err != nil {
		t.Fatalf("Failed to create unified client: %v", err)
	}

	ctx := context.Background()

	// Test with a non-existent ID (should fail for both APIs)
	detail, err := client.GetDetail(ctx, "INVALID_ID_12345")
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
	if detail != nil {
		t.Error("Expected nil detail for invalid ID")
	}
}
