package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/config"
	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
)

// ELISClient handles ELIS (자치법규정보시스템) API requests
// Note: ELIS API is provided through the National Law Information Center API
type ELISClient struct {
	httpClient     *http.Client
	baseURL        string // 자치법규 목록 조회
	detailURL      string // 자치법규 본문 조회
	apiKey         string
	retryBaseDelay time.Duration
	maxRetries     int
}

// NewELISClient creates a new ELIS API client
func NewELISClient(apiKey string) *ELISClient {
	return &ELISClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:        "https://www.law.go.kr/DRF/lawSearch.do",  // 자치법규 목록
		detailURL:      "https://www.law.go.kr/DRF/lawService.do", // 자치법규 본문
		apiKey:         apiKey,
		retryBaseDelay: 500 * time.Millisecond,
		maxRetries:     3,
	}
}

// OrdinanceInfo represents local ordinance information
type OrdinanceInfo struct {
	ID         string `json:"자치법규ID" xml:"자치법규ID"`
	Name       string `json:"자치법규명" xml:"자치법규명"`
	LocalGov   string `json:"지자체명" xml:"지자체명"`
	Department string `json:"소관부서" xml:"소관부서"`
	PromulDate string `json:"공포일자" xml:"공포일자"`
	PromulNo   string `json:"공포번호" xml:"공포번호"`
	EffectDate string `json:"시행일자" xml:"시행일자"`
	Category   string `json:"구분" xml:"구분"` // 조례/규칙
	Status     string `json:"상태" xml:"상태"` // 현행/폐지
}

// OrdinanceDetail represents detailed ordinance information
type OrdinanceDetail struct {
	OrdinanceInfo
	Content     string    `json:"본문" xml:"본문"`
	Articles    []Article `json:"조문" xml:"조문"`
	Attachments []string  `json:"첨부파일" xml:"첨부파일"`
	RelatedLaws []string  `json:"관련법령" xml:"관련법령"`
}

// OrdinanceSearchResponse represents the search response for ordinances
type OrdinanceSearchResponse struct {
	OrdinSearch struct {
		ResultCode string                   `json:"resultCode"`
		ResultMsg  string                   `json:"resultMsg"`
		TotalCnt   string                   `json:"totalCnt"`
		Page       string                   `json:"page"`
		NumOfRows  string                   `json:"numOfRows"`
		Law        []map[string]interface{} `json:"law"`
	} `json:"OrdinSearch"`
}

// Search searches for local ordinances
func (c *ELISClient) Search(ctx context.Context, req *UnifiedSearchRequest) (*SearchResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "ordin") // 자치법규 대상

	// Add region to query if provided
	query := req.Query
	if req.Region != "" {
		query = req.Region + " " + query
	}
	params.Set("query", query)

	params.Set("page", fmt.Sprintf("%d", req.PageNo))
	params.Set("display", fmt.Sprintf("%d", req.PageSize))
	params.Set("type", "json")

	// Add sort order
	if req.Sort != "" {
		params.Set("sort", req.Sort)
	} else {
		params.Set("sort", "date") // 기본값: 날짜순
	}

	fullURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	logger.Debug("ELIS API Request URL: %s", fullURL)

	// Make request with retry logic
	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, fmt.Errorf("ELIS API 요청 실패: %w", err)
	}

	// Parse response
	var elisResp OrdinanceSearchResponse
	if err := json.Unmarshal(body, &elisResp); err != nil {
		// Check if response is HTML (error page)
		bodyStr := string(body)
		if strings.HasPrefix(strings.TrimSpace(bodyStr), "<!DOCTYPE") || strings.HasPrefix(strings.TrimSpace(bodyStr), "<html") {
			errorMsg := c.parseHTMLError(bodyStr)
			logger.Debug("HTML error response detected: %s", errorMsg)
			// Check if it's an API key error
			if strings.Contains(errorMsg, "API 인증 실패") || strings.Contains(errorMsg, "API 키") {
				return nil, &APIKeyError{Message: errorMsg}
			}
			return nil, fmt.Errorf("%s", errorMsg)
		}
		logger.Debug("Response body: %s", bodyStr)
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	// Check for API error
	if elisResp.OrdinSearch.ResultCode != "00" {
		return nil, fmt.Errorf("API 오류: %s", elisResp.OrdinSearch.ResultMsg)
	}

	// Parse total count and page
	totalCount := 0
	if tc, err := strconv.Atoi(elisResp.OrdinSearch.TotalCnt); err == nil {
		totalCount = tc
	}
	page := 1
	if p, err := strconv.Atoi(elisResp.OrdinSearch.Page); err == nil {
		page = p
	}

	// Convert to unified SearchResponse
	searchResp := &SearchResponse{
		TotalCount: totalCount,
		Page:       page,
		Laws:       make([]LawInfo, 0, len(elisResp.OrdinSearch.Law)),
	}

	for _, lawData := range elisResp.OrdinSearch.Law {
		law := LawInfo{
			LawType: "자치법규",
		}

		// Parse fields from map
		if v, ok := lawData["자치법규ID"].(string); ok {
			law.ID = v
		}
		if v, ok := lawData["자치법규명"].(string); ok {
			law.Name = v
		}
		if v, ok := lawData["자치법규일련번호"].(string); ok {
			law.SerialNo = v
		}
		if v, ok := lawData["공포일자"].(string); ok {
			law.PromulDate = v
		}
		if v, ok := lawData["공포번호"].(string); ok {
			law.PromulNo = v
		}
		if v, ok := lawData["자치법규종류"].(string); ok {
			law.Category = v
		}
		if v, ok := lawData["지자체기관명"].(string); ok {
			law.Department = v
		}
		if v, ok := lawData["시행일자"].(string); ok {
			law.EffectDate = v
		}

		searchResp.Laws = append(searchResp.Laws, law)
	}

	return searchResp, nil
}

// GetDetail retrieves detailed ordinance information
func (c *ELISClient) GetDetail(ctx context.Context, ordinanceID string) (*LawDetail, error) {
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "ordin")
	params.Set("MST", ordinanceID) // Use MST parameter for ordinance serial number
	params.Set("type", "json")

	fullURL := fmt.Sprintf("%s?%s", c.detailURL, params.Encode())
	logger.Debug("ELIS Detail API Request URL: %s", fullURL)

	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, fmt.Errorf("자치법규 상세 조회 실패: %w", err)
	}

	// Parse the response
	var detailResp map[string]interface{}
	if err := json.Unmarshal(body, &detailResp); err != nil {
		// Check if response is HTML (error page)
		bodyStr := string(body)
		if strings.HasPrefix(strings.TrimSpace(bodyStr), "<!DOCTYPE") || strings.HasPrefix(strings.TrimSpace(bodyStr), "<html") {
			errorMsg := c.parseHTMLError(bodyStr)
			logger.Debug("HTML error response detected: %s", errorMsg)
			// Check if it's an API key error
			if strings.Contains(errorMsg, "API 인증 실패") || strings.Contains(errorMsg, "API 키") {
				return nil, &APIKeyError{Message: errorMsg}
			}
			return nil, fmt.Errorf("%s", errorMsg)
		}
		return nil, fmt.Errorf("자치법규 상세 정보 파싱 실패: %w", err)
	}

	// Check for error message
	if lawMsg, ok := detailResp["Law"].(string); ok {
		if strings.Contains(lawMsg, "일치하는 자치법규가 없습니다") {
			return nil, fmt.Errorf("자치법규 상세 정보를 찾을 수 없습니다")
		}
	}

	// Parse LawService response
	lawService, ok := detailResp["LawService"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("잘못된 응답 형식")
	}

	// Parse basic info
	basicInfo, ok := lawService["자치법규기본정보"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("자치법규 기본 정보를 찾을 수 없습니다")
	}

	// Convert to LawDetail
	lawDetail := &LawDetail{
		LawInfo: LawInfo{
			LawType: "자치법규",
		},
	}

	// Parse fields from basicInfo
	if v, ok := basicInfo["자치법규ID"].(string); ok {
		lawDetail.LawInfo.ID = v
	}
	if v, ok := basicInfo["자치법규명"].(string); ok {
		lawDetail.LawInfo.Name = v
	}
	if v, ok := basicInfo["자치법규일련번호"].(string); ok {
		lawDetail.LawInfo.SerialNo = v
	}
	if v, ok := basicInfo["공포일자"].(string); ok {
		lawDetail.LawInfo.PromulDate = v
	}
	if v, ok := basicInfo["공포번호"].(string); ok {
		lawDetail.LawInfo.PromulNo = v
	}
	if v, ok := basicInfo["자치법규종류"].(string); ok {
		if v == "C0001" {
			lawDetail.LawInfo.Category = "조례"
		} else if v == "C0002" {
			lawDetail.LawInfo.Category = "규칙"
		} else {
			lawDetail.LawInfo.Category = v
		}
	}
	if v, ok := basicInfo["지자체기관명"].(string); ok {
		lawDetail.LawInfo.Department = v
	}
	if v, ok := basicInfo["시행일자"].(string); ok {
		lawDetail.LawInfo.EffectDate = v
	}

	// Parse content (articles)
	if articles, ok := lawService["조문단위"].([]interface{}); ok {
		lawDetail.Articles = make([]Article, 0, len(articles))
		for _, art := range articles {
			if artMap, ok := art.(map[string]interface{}); ok {
				article := Article{}
				if v, ok := artMap["조문번호"].(string); ok {
					article.Number = v
				}
				if v, ok := artMap["조문제목"].(string); ok {
					article.Title = v
				}
				if v, ok := artMap["조문내용"].(string); ok {
					article.Content = v
				}
				lawDetail.Articles = append(lawDetail.Articles, article)
			}
		}
	}

	return lawDetail, nil
}

// GetHistory retrieves ordinance amendment history
func (c *ELISClient) GetHistory(ctx context.Context, ordinanceID string) (*LawHistory, error) {
	// 자치법규의 경우 제/개정 이력 조회가 별도로 제공되지 않을 수 있음
	// 추후 API 확인 후 구현
	return nil, fmt.Errorf("자치법규 이력 조회는 현재 지원되지 않습니다")
}

// FilterByRegion filters ordinances by local government
func (c *ELISClient) FilterByRegion(ctx context.Context, region string, query string, pageNo int, pageSize int) (*SearchResponse, error) {
	req := &UnifiedSearchRequest{
		Query:    query,
		Region:   region,
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	return c.Search(ctx, req)
}

// doRequestWithRetry performs HTTP request with retry logic
func (c *ELISClient) doRequestWithRetry(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := c.retryBaseDelay * time.Duration(1<<uint(attempt-1))
			logger.Debug("Retrying after %v (attempt %d/%d)", delay, attempt+1, c.maxRetries)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("요청 생성 실패: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP 요청 실패: %w", err)
			continue
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode == http.StatusServiceUnavailable ||
			resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("서버 에러: HTTP %d", resp.StatusCode)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("클라이언트 에러: HTTP %d", resp.StatusCode)
		}

		// Read response body
		body := make([]byte, 0, 1024*1024) // Pre-allocate 1MB
		buf := make([]byte, 32*1024)       // 32KB buffer
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				body = append(body, buf[:n]...)
			}
			if err != nil {
				if err.Error() != "EOF" && !strings.Contains(err.Error(), "closed") {
					return nil, fmt.Errorf("응답 읽기 실패: %w", err)
				}
				break
			}
		}

		return body, nil
	}

	return nil, fmt.Errorf("최대 재시도 횟수 초과: %w", lastErr)
}

// GetELISAPIKey gets the ELIS API key from config
func GetELISAPIKey() string {
	return config.GetString("law.elis.key")
}

// GetAPIType returns the API type
func (c *ELISClient) GetAPIType() APIType {
	return APITypeELIS
}

// parseHTMLError extracts meaningful error message from HTML error page
func (c *ELISClient) parseHTMLError(html string) string {
	// Common patterns for error messages in HTML pages
	patterns := []struct {
		start string
		end   string
	}{
		// Look for error messages in common formats
		{"<title>", "</title>"},
		{"<h1>", "</h1>"},
		{"<h2>", "</h2>"},
		{"class=\"error\"", "</"},
		{"class=\"message\"", "</"},
		{"오류", "</"},
		{"에러", "</"},
		{"Error", "</"},
		{"인증", "</"},
	}

	htmlLower := strings.ToLower(html)

	// Check for authentication/key related issues
	if strings.Contains(htmlLower, "인증") || strings.Contains(htmlLower, "auth") ||
		strings.Contains(htmlLower, "key") || strings.Contains(htmlLower, "키") {
		return "API 인증 실패: API 키가 유효하지 않거나 만료되었습니다. 'warp config set law.key YOUR_API_KEY' 명령으로 올바른 API 키를 설정하세요"
	}

	// Check for rate limit
	if strings.Contains(htmlLower, "limit") || strings.Contains(htmlLower, "제한") {
		return "API 호출 제한 초과: 일일 호출 한도를 초과했습니다. 잠시 후 다시 시도하세요"
	}

	// Check for service unavailable
	if strings.Contains(htmlLower, "maintenance") || strings.Contains(htmlLower, "점검") {
		return "서비스 점검 중: 자치법규정보시스템 API가 현재 점검 중입니다"
	}

	// Try to extract error message from patterns
	for _, pattern := range patterns {
		startIdx := strings.Index(htmlLower, pattern.start)
		if startIdx != -1 {
			startIdx += len(pattern.start)
			endIdx := strings.Index(htmlLower[startIdx:], pattern.end)
			if endIdx != -1 {
				msg := html[startIdx : startIdx+endIdx]
				msg = strings.TrimSpace(msg)
				// Remove any remaining HTML tags
				msg = stripHTMLTags(msg)
				if len(msg) > 10 && len(msg) < 200 {
					// Check if it's just a generic title
					if strings.Contains(msg, "국가법령정보") && !strings.Contains(msg, "오류") {
						// Generic title, return more specific message
						return "API 인증 실패: API 키가 올바르지 않습니다. 'warp config set law.key YOUR_API_KEY' 명령으로 유효한 API 키를 설정하세요"
					}
					return fmt.Sprintf("API 오류: %s", msg)
				}
			}
		}
	}

	// Default message if no specific error found
	return "API 요청 실패: 서버가 예상하지 못한 응답을 반환했습니다. API 키를 확인하거나 잠시 후 다시 시도하세요"
}

// Ensure ELISClient implements ClientInterface
var _ ClientInterface = (*ELISClient)(nil)
