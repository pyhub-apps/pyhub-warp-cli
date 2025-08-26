package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
)

// NLICClient represents the National Law Information Center API client
type NLICClient struct {
	httpClient     *http.Client
	baseURL        string
	detailURL      string
	historyURL     string
	apiKey         string
	retryBaseDelay time.Duration
}

// NewNLICClient creates a new NLIC API client
func NewNLICClient(apiKey string) *NLICClient {
	return &NLICClient{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:        BaseURL,
		detailURL:      "https://www.law.go.kr/DRF/lawService.do",
		historyURL:     "https://www.law.go.kr/DRF/lawHistory.do",
		apiKey:         apiKey,
		retryBaseDelay: InitialRetryDelay,
	}
}

// NewNLICClientWithURL creates a new NLIC API client with custom URLs (for testing)
func NewNLICClientWithURL(apiKey, baseURL string) *NLICClient {
	return &NLICClient{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:        baseURL,
		detailURL:      baseURL, // Use same URL for testing
		historyURL:     baseURL, // Use same URL for testing
		apiKey:         apiKey,
		retryBaseDelay: InitialRetryDelay,
	}
}

// GetAPIType returns the API type
func (c *NLICClient) GetAPIType() APIType {
	return APITypeNLIC
}

// Search performs a law search
func (c *NLICClient) Search(ctx context.Context, req *UnifiedSearchRequest) (*SearchResponse, error) {
	// Set defaults
	if req.Type == "" {
		req.Type = "JSON"
	}
	if req.PageNo == 0 {
		req.PageNo = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// Build URL with parameters
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "law")
	params.Set("query", req.Query)
	params.Set("type", req.Type)
	params.Set("page", fmt.Sprintf("%d", req.PageNo))
	params.Set("display", fmt.Sprintf("%d", req.PageSize))

	// Add optional filters
	if req.LawType != "" {
		params.Set("법령구분", req.LawType)
	}
	if req.Department != "" {
		params.Set("소관부처", req.Department)
	}
	if req.Sort != "" {
		params.Set("sort", req.Sort)
	}

	fullURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	logger.Debug("API Request URL: %s", fullURL)

	// Perform request with retries
	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Parse response
	var searchResp SearchResponse
	if req.Type == "JSON" {
		if err := json.Unmarshal(body, &searchResp); err != nil {
			// Check if response is HTML (error page)
			bodyStr := string(body)
			if strings.HasPrefix(strings.TrimSpace(bodyStr), "<!DOCTYPE") || strings.HasPrefix(strings.TrimSpace(bodyStr), "<html") {
				// Parse HTML error message
				errorMsg := ParseHTMLError(bodyStr)
				logger.Debug("HTML error response detected: %s", errorMsg)
				// Check if it's an API key error
				if strings.Contains(errorMsg, "API 인증 실패") || strings.Contains(errorMsg, "API 키") {
					return nil, &APIKeyError{Message: errorMsg}
				}
				return nil, fmt.Errorf("%s", errorMsg)
			}
			// Log non-HTML parsing errors for debugging
			if len(bodyStr) > 500 {
				bodyStr = bodyStr[:500] + "..."
			}
			logger.Debug("API Response (first 500 chars): %s", bodyStr)
			return nil, fmt.Errorf("응답 데이터 파싱 실패: %w", err)
		}
	} else {
		if err := xml.Unmarshal(body, &searchResp); err != nil {
			return nil, fmt.Errorf("XML 파싱 실패: %w", err)
		}
	}

	// Set the page number if not returned by API
	if searchResp.Page == 0 {
		searchResp.Page = req.PageNo
	}

	return &searchResp, nil
}

// GetDetail retrieves detailed law information
func (c *NLICClient) GetDetail(ctx context.Context, lawID string) (*LawDetail, error) {
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "law")
	params.Set("MST", lawID) // 법령일련번호
	params.Set("type", "JSON")

	fullURL := fmt.Sprintf("%s?%s", c.detailURL, params.Encode())

	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Debug: Log the raw response
	maxLen := 2000
	if len(body) < maxLen {
		maxLen = len(body)
	}
	logger.Debug("Law Detail API Response (first %d chars): %s", maxLen, string(body[:maxLen]))

	// Parse the response using the correct structure
	var detailResp LawDetailResponse
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
		return nil, fmt.Errorf("응답 데이터 파싱 실패: %w", err)
	}

	// Convert the detailed response to our LawDetail structure
	detail := &LawDetail{
		LawInfo: LawInfo{
			SerialNo: lawID, // Use the provided law ID
		},
		Articles:                 make([]Article, 0),
		Tables:                   make([]Table, 0),
		SupplementaryProvisions:  make([]SupplementaryProvision, 0),
	}

	// Extract basic info if available
	if detailResp.Law.BasicInfo != nil {
		basicInfo := detailResp.Law.BasicInfo
		detail.LawInfo.ID = basicInfo.LawID
		detail.LawInfo.Name = basicInfo.LawNameKorean
		detail.LawInfo.PromulDate = basicInfo.PromulgationDate
		detail.LawInfo.PromulNo = basicInfo.PromulgationNumber
		detail.LawInfo.EffectDate = basicInfo.EffectiveDate
		detail.LawInfo.Category = basicInfo.RevisionType
		
		// Extract department info if available
		if basicInfo.Department.Content != "" {
			detail.LawInfo.Department = basicInfo.Department.Content
		}
		
		// Extract law type info if available
		if basicInfo.LawTypeInfo.Content != "" {
			detail.LawInfo.LawType = basicInfo.LawTypeInfo.Content
		}
	}

	// Process revision text
	if detailResp.Law.Revisions.Content != nil {
		detail.HasRevisionText = true
		// Convert revision content to string if it exists
		if revStr, ok := detailResp.Law.Revisions.Content.(string); ok {
			detail.RevisionText = revStr
		} else {
			// Mark that revision exists but don't try to parse complex structure
			detail.RevisionText = "(개정문 내용 있음)"
		}
	}

	// Convert tables from the API response
	for _, unit := range detailResp.Law.Tables.TableUnits {
		table := Table{
			Number: unit.TableNumber,
			Title:  unit.TableTitle,
		}
		
		// Handle table content which can be string or array
		if unit.TableContent != nil {
			if tableStr, ok := unit.TableContent.(string); ok {
				table.Content = tableStr
			} else {
				// If it's an array or other type, just mark as existing
				table.Content = "(별표 내용 있음)"
			}
		}
		
		detail.Tables = append(detail.Tables, table)
	}

	// Convert supplementary provisions from the API response
	for _, unit := range detailResp.Law.SupplementaryProvisions.ProvisionUnits {
		supp := SupplementaryProvision{
			Number:           unit.ProvisionNumber,
			PromulgationDate: unit.ProvisionDate,
		}
		
		// Handle provision content which can be string or array
		if unit.ProvisionContent != nil {
			if provStr, ok := unit.ProvisionContent.(string); ok {
				supp.Content = provStr
			} else if provArr, ok := unit.ProvisionContent.([]interface{}); ok {
				// If it's an array, join the elements
				var contentParts []string
				for _, part := range provArr {
					if str, ok := part.(string); ok {
						contentParts = append(contentParts, str)
					}
				}
				supp.Content = strings.Join(contentParts, "\n")
			} else {
				// If it's something else, just mark as existing
				supp.Content = "(부칙 내용 있음)"
			}
		}
		
		detail.SupplementaryProvisions = append(detail.SupplementaryProvisions, supp)
	}

	// Convert articles from the API response
	for _, unit := range detailResp.Law.ArticlesRaw.ArticleUnits {
		article := Article{
			Number:     unit.ArticleNumber,
			Title:      unit.ArticleTitle,
			Content:    unit.ArticleContent,
			EffectDate: unit.ArticleEffectDate,
		}
		detail.Articles = append(detail.Articles, article)
		
		// If basic info wasn't in the main structure, try to get from article units
		if detail.LawInfo.ID == "" && unit.LawID != "" {
			detail.LawInfo.ID = unit.LawID
		}
		if detail.LawInfo.Name == "" && unit.LawNameKorean != "" {
			detail.LawInfo.Name = unit.LawNameKorean
		}
		if detail.LawInfo.SerialNo == "" && unit.LawSerialNo != "" {
			detail.LawInfo.SerialNo = unit.LawSerialNo
		}
	}

	// Set the law key if available and ID is still empty
	if detailResp.Law.LawKey != "" && detail.LawInfo.ID == "" {
		detail.LawInfo.ID = detailResp.Law.LawKey
	}

	return detail, nil
}

// GetHistory retrieves law amendment history
func (c *NLICClient) GetHistory(ctx context.Context, lawID string) (*LawHistory, error) {
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "law")
	params.Set("MST", lawID) // 법령일련번호
	params.Set("type", "JSON")

	fullURL := fmt.Sprintf("%s?%s", c.historyURL, params.Encode())

	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var history LawHistory
	if err := json.Unmarshal(body, &history); err != nil {
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
		// Try parsing as a wrapper
		var wrapper struct {
			History *LawHistory `json:"연혁" xml:"연혁"`
		}
		if err2 := json.Unmarshal(body, &wrapper); err2 != nil {
			return nil, fmt.Errorf("응답 데이터 파싱 실패: %w", err)
		}
		if wrapper.History != nil {
			return wrapper.History, nil
		}
	}

	// Fill in the law ID if not returned
	if history.LawID == "" {
		history.LawID = lawID
	}

	return &history, nil
}

// doRequestWithRetry performs an HTTP request with retry logic
func (c *NLICClient) doRequestWithRetry(ctx context.Context, url string) ([]byte, error) {
	var lastErr error
	retryDelay := c.retryBaseDelay

	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			select {
			case <-time.After(retryDelay):
				retryDelay *= 2
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		body, err := c.doRequest(ctx, url)
		if err != nil {
			lastErr = err
			// Only retry on network errors or 5xx server errors
			if !c.shouldRetry(err) {
				return nil, err
			}
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("요청 실패 (재시도 %d회 초과): %w", MaxRetries, lastErr)
}

// doRequest performs a single HTTP request
func (c *NLICClient) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Do not retry on explicit context cancellation or deadline
		if ctx.Err() != nil {
			return nil, fmt.Errorf("요청이 취소되었거나 시간 초과되었습니다: %w", ctx.Err())
		}
		// Treat other transport-level problems as retryable
		return nil, &RetryableError{Err: fmt.Errorf("네트워크 에러: %w", err)}
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleHTTPError(resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", err)
	}

	// Check for API errors in the response
	if c.hasAPIError(body) {
		return nil, c.parseAPIError(body)
	}

	return body, nil
}

// handleHTTPError converts HTTP status codes to appropriate errors
func (c *NLICClient) handleHTTPError(statusCode int) error {
	switch statusCode {
	case http.StatusTooManyRequests:
		return &RetryableError{Err: fmt.Errorf("레이트 리밋: HTTP 429 (잠시 후 다시 시도하세요)")}
	case http.StatusRequestTimeout:
		return &RetryableError{Err: fmt.Errorf("요청 타임아웃: HTTP 408")}
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return &RetryableError{Err: fmt.Errorf("일시적 서버 오류: HTTP %d", statusCode)}
	case http.StatusInternalServerError:
		return &RetryableError{Err: fmt.Errorf("내부 서버 오류: HTTP 500")}
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("인증 실패: HTTP %d - API 키를 확인하세요", statusCode)
	default:
		if statusCode >= 500 {
			return &RetryableError{Err: fmt.Errorf("서버 에러: HTTP %d", statusCode)}
		}
		return fmt.Errorf("클라이언트 에러: HTTP %d", statusCode)
	}
}

// shouldRetry determines if an error is retryable
func (c *NLICClient) shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's explicitly marked as retryable
	var retryableErr *RetryableError
	if errors.As(err, &retryableErr) {
		return true
	}

	// Also check for network-related errors in the message
	return strings.Contains(err.Error(), "네트워크") || strings.Contains(err.Error(), "타임아웃")
}

// hasAPIError checks if the response contains an API error
func (c *NLICClient) hasAPIError(body []byte) bool {
	// Simple check for error indicators in the response
	bodyStr := string(body)
	return strings.Contains(bodyStr, `"error"`) || strings.Contains(bodyStr, "<error>") ||
		strings.Contains(bodyStr, `"에러"`) || strings.Contains(bodyStr, "ERROR")
}

// parseAPIError extracts API error from response body
func (c *NLICClient) parseAPIError(body []byte) error {
	var errResp struct {
		Error     *ErrorInfo `json:"error" xml:"error"`
		ErrorMsg  string     `json:"errorMsg" xml:"errorMsg"`
		ErrorCode string     `json:"errorCode" xml:"errorCode"`
	}

	// Try JSON first
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Error != nil {
			return fmt.Errorf("API 에러: %s - %s", errResp.Error.Code, errResp.Error.Message)
		}
		if errResp.ErrorMsg != "" {
			return fmt.Errorf("API 에러: %s", errResp.ErrorMsg)
		}
	}

	return fmt.Errorf("알 수 없는 API 에러")
}

// parseHTMLError extracts meaningful error message from HTML error page
func (c *NLICClient) parseHTMLError(html string) string {
	// Check for specific error messages in the HTML response
	if strings.Contains(html, "미신청된 목록/본문에 대한 접근입니다") {
		return "API 사용 권한이 없습니다. https://open.law.go.kr 에서 로그인 후 [OPEN API] -> [OPEN API 신청]에서 필요한 법령 종류를 체크해주세요"
	}
	
	if strings.Contains(html, "페이지 접속에 실패하였습니다") {
		return "API 접속 실패: API 키를 확인하거나 서비스 상태를 점검해주세요"
	}

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
		return "API 인증 실패: 이메일 ID가 올바르지 않습니다. 'warp config set law.key YOUR_EMAIL_ID' 명령으로 이메일 @ 앞부분을 설정하세요"
	}

	// Check for rate limit
	if strings.Contains(htmlLower, "limit") || strings.Contains(htmlLower, "제한") {
		return "API 호출 제한 초과: 일일 호출 한도를 초과했습니다. 잠시 후 다시 시도하세요"
	}

	// Check for service unavailable
	if strings.Contains(htmlLower, "maintenance") || strings.Contains(htmlLower, "점검") {
		return "서비스 점검 중: 국가법령정보센터 API가 현재 점검 중입니다"
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

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(s string) string {
	// Simple regex-like approach to remove HTML tags
	result := ""
	inTag := false
	for _, ch := range s {
		if ch == '<' {
			inTag = true
		} else if ch == '>' {
			inTag = false
		} else if !inTag {
			result += string(ch)
		}
	}
	return strings.TrimSpace(result)
}
