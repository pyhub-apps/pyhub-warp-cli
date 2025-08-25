package api

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/pyhub-apps/sejong-cli/internal/config"
	"github.com/pyhub-apps/sejong-cli/internal/logger"
)

// UnifiedClient handles unified search across multiple APIs
type UnifiedClient struct {
	nlicClient *NLICClient
	elisClient *ELISClient
}

// NewUnifiedClient creates a new unified API client
func NewUnifiedClient() (*UnifiedClient, error) {
	// Get API key (both NLIC and ELIS use the same key from law.go.kr)
	apiKey := config.GetNLICAPIKey()
	if apiKey == "" {
		// Try legacy key path
		apiKey = config.GetString("law.key")
		if apiKey == "" {
			return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 'sejong config set law.key YOUR_KEY' 명령으로 설정하세요")
		}
	}

	return &UnifiedClient{
		nlicClient: NewNLICClient(apiKey),
		elisClient: NewELISClient(apiKey),
	}, nil
}

// Search performs parallel search across all APIs
func (c *UnifiedClient) Search(ctx context.Context, req *UnifiedSearchRequest) (*SearchResponse, error) {
	// Create channels for results
	type searchResult struct {
		source   string
		response *SearchResponse
		err      error
	}

	resultsChan := make(chan searchResult, 2)
	var wg sync.WaitGroup

	// Search NLIC in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Debug("Starting NLIC search for: %s", req.Query)
		resp, err := c.nlicClient.Search(ctx, req)
		resultsChan <- searchResult{
			source:   "NLIC",
			response: resp,
			err:      err,
		}
	}()

	// Search ELIS in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Debug("Starting ELIS search for: %s", req.Query)
		resp, err := c.elisClient.Search(ctx, req)
		resultsChan <- searchResult{
			source:   "ELIS",
			response: resp,
			err:      err,
		}
	}()

	// Wait for all searches to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var allLaws []LawInfo
	totalCount := 0
	errors := []error{}

	for result := range resultsChan {
		if result.err != nil {
			logger.Error("%s search error: %v", result.source, result.err)
			errors = append(errors, fmt.Errorf("%s: %w", result.source, result.err))
			continue
		}

		if result.response != nil {
			logger.Debug("%s returned %d results", result.source, len(result.response.Laws))

			// Add source information to each law
			for i := range result.response.Laws {
				if result.source == "ELIS" {
					// Mark ELIS results
					result.response.Laws[i].Source = "자치법규"
				} else {
					// Mark NLIC results
					result.response.Laws[i].Source = "국가법령"
				}
			}

			allLaws = append(allLaws, result.response.Laws...)
			totalCount += result.response.TotalCount
		}
	}

	// If all searches failed, return error
	if len(errors) == 2 {
		return nil, fmt.Errorf("모든 API 검색 실패: %v", errors)
	}

	// Sort results by date (newest first)
	sort.Slice(allLaws, func(i, j int) bool {
		// Sort by promulgation date, newest first
		return allLaws[i].PromulDate > allLaws[j].PromulDate
	})

	// Apply pagination
	startIdx := (req.PageNo - 1) * req.PageSize
	endIdx := startIdx + req.PageSize

	if startIdx > len(allLaws) {
		startIdx = len(allLaws)
	}
	if endIdx > len(allLaws) {
		endIdx = len(allLaws)
	}

	paginatedLaws := allLaws[startIdx:endIdx]

	// Create unified response
	response := &SearchResponse{
		TotalCount: totalCount,
		Page:       req.PageNo,
		Laws:       paginatedLaws,
	}

	logger.Info("통합 검색 완료: 총 %d개 결과 (NLIC+ELIS)", len(allLaws))

	return response, nil
}

// GetDetail retrieves detailed information (tries NLIC first, then ELIS)
func (c *UnifiedClient) GetDetail(ctx context.Context, lawID string) (*LawDetail, error) {
	// Try NLIC first (for national laws)
	detail, err := c.nlicClient.GetDetail(ctx, lawID)
	if err == nil {
		return detail, nil
	}

	logger.Debug("NLIC detail failed, trying ELIS: %v", err)

	// Try ELIS (for local ordinances)
	detail, err = c.elisClient.GetDetail(ctx, lawID)
	if err == nil {
		return detail, nil
	}

	return nil, fmt.Errorf("법령/조례 상세 정보를 찾을 수 없습니다 (ID: %s)", lawID)
}

// GetHistory retrieves law history (only NLIC supports this)
func (c *UnifiedClient) GetHistory(ctx context.Context, lawID string) (*LawHistory, error) {
	// Only NLIC supports history
	return c.nlicClient.GetHistory(ctx, lawID)
}

// GetAPIType returns the API type
func (c *UnifiedClient) GetAPIType() APIType {
	return APITypeAll
}

// SearchWithOptions performs search with specific API selection
func (c *UnifiedClient) SearchWithOptions(ctx context.Context, req *UnifiedSearchRequest, includeNLIC, includeELIS bool) (*SearchResponse, error) {
	if !includeNLIC && !includeELIS {
		return nil, fmt.Errorf("최소 하나의 API를 선택해야 합니다")
	}

	// If only one API selected, use single client
	if includeNLIC && !includeELIS {
		return c.nlicClient.Search(ctx, req)
	}
	if includeELIS && !includeNLIC {
		return c.elisClient.Search(ctx, req)
	}

	// Both selected, use unified search
	return c.Search(ctx, req)
}

// Ensure UnifiedClient implements ClientInterface
var _ ClientInterface = (*UnifiedClient)(nil)
