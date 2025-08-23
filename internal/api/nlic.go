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

	// Perform request with retries
	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Parse response
	var searchResp SearchResponse
	if req.Type == "JSON" {
		if err := json.Unmarshal(body, &searchResp); err != nil {
			return nil, fmt.Errorf("JSON 파싱 실패: %w", err)
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
	params.Set("MST", lawID)  // 법령일련번호
	params.Set("type", "JSON")

	fullURL := fmt.Sprintf("%s?%s", c.detailURL, params.Encode())

	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Parse the response - the API returns the detail in a wrapper
	var wrapper struct {
		Law *LawDetail `json:"법령" xml:"법령"`
	}
	
	if err := json.Unmarshal(body, &wrapper); err != nil {
		// Try direct unmarshal if wrapper fails
		var detail LawDetail
		if err2 := json.Unmarshal(body, &detail); err2 != nil {
			return nil, fmt.Errorf("JSON 파싱 실패: %w", err)
		}
		return &detail, nil
	}

	if wrapper.Law == nil {
		return nil, fmt.Errorf("법령 상세 정보를 찾을 수 없습니다")
	}

	return wrapper.Law, nil
}

// GetHistory retrieves law amendment history
func (c *NLICClient) GetHistory(ctx context.Context, lawID string) (*LawHistory, error) {
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "law")
	params.Set("MST", lawID)  // 법령일련번호
	params.Set("type", "JSON")

	fullURL := fmt.Sprintf("%s?%s", c.historyURL, params.Encode())

	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var history LawHistory
	if err := json.Unmarshal(body, &history); err != nil {
		// Try parsing as a wrapper
		var wrapper struct {
			History *LawHistory `json:"연혁" xml:"연혁"`
		}
		if err2 := json.Unmarshal(body, &wrapper); err2 != nil {
			return nil, fmt.Errorf("JSON 파싱 실패: %w", err)
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
		Error *ErrorInfo `json:"error" xml:"error"`
		ErrorMsg string `json:"errorMsg" xml:"errorMsg"`
		ErrorCode string `json:"errorCode" xml:"errorCode"`
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