package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/logger"
)

// PrecSearchResponse represents the precedent search response
type PrecSearchResponse struct {
	XMLName    xml.Name   `xml:"PrecSearch"`
	TotalCount int        `xml:"totalCnt"`
	Page       int        `xml:"page"`
	Precs      []PrecInfo `xml:"prec"`
}

// PrecInfo represents individual precedent information
type PrecInfo struct {
	ID         string `xml:"판례일련번호"`
	CaseName   string `xml:"사건명"`
	CaseNumber string `xml:"사건번호"`
	CourtName  string `xml:"법원명"`
	JudgeDate  string `xml:"선고일자"`
	CaseType   string `xml:"사건종류명"`
	DetailLink string `xml:"판례상세링크"`
}

// PrecClient represents the Precedent API client (판례 API 클라이언트)
type PrecClient struct {
	httpClient     *http.Client
	baseURL        string
	detailURL      string
	apiKey         string
	retryBaseDelay time.Duration
}

// NewPrecClient creates a new Precedent API client
func NewPrecClient(apiKey string) *PrecClient {
	return &PrecClient{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:        "https://www.law.go.kr/DRF/lawSearch.do",
		detailURL:      "https://www.law.go.kr/DRF/lawService.do",
		apiKey:         apiKey,
		retryBaseDelay: InitialRetryDelay,
	}
}

// GetAPIType returns the API type
func (c *PrecClient) GetAPIType() APIType {
	return APITypePrec
}

// Search performs a precedent search
func (c *PrecClient) Search(ctx context.Context, req *UnifiedSearchRequest) (*SearchResponse, error) {
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
	params.Set("target", "prec") // 판례 검색
	params.Set("query", req.Query)
	params.Set("type", req.Type)
	params.Set("page", fmt.Sprintf("%d", req.PageNo))
	params.Set("display", fmt.Sprintf("%d", req.PageSize))

	// Add optional filters
	if req.Sort != "" {
		params.Set("sort", req.Sort)
	}

	fullURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	logger.Debug("Precedent API Request URL: %s", fullURL)

	// Perform request with retries
	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Parse response based on type
	if strings.ToUpper(req.Type) == "JSON" {
		// JSON is not supported for precedent API, return empty result
		logger.Debug("JSON format not supported for precedent API, returning empty result")
		return &SearchResponse{TotalCount: 0, Page: req.PageNo, Laws: []LawInfo{}}, nil
	}

	// Parse XML response
	var precResponse PrecSearchResponse
	if err := xml.Unmarshal(body, &precResponse); err != nil {
		logger.Error("XML parsing failed: %v", err)
		return nil, fmt.Errorf("XML 파싱 실패: %w", err)
	}

	// Convert to SearchResponse
	response := &SearchResponse{
		TotalCount: precResponse.TotalCount,
		Page:       precResponse.Page,
		Laws:       make([]LawInfo, len(precResponse.Precs)),
	}

	for i, prec := range precResponse.Precs {
		response.Laws[i] = LawInfo{
			ID:         prec.ID,
			Name:       prec.CaseName,
			LawType:    prec.CaseType,
			Department: prec.CourtName,
			PromulDate: prec.JudgeDate,
			PromulNo:   prec.CaseNumber,
		}
	}

	return response, nil
}

// GetDetail retrieves detailed precedent information
func (c *PrecClient) GetDetail(ctx context.Context, precID string) (*LawDetail, error) {
	// Build URL with parameters
	params := url.Values{}
	params.Set("OC", c.apiKey)
	params.Set("target", "prec") // 판례 상세
	params.Set("ID", precID)
	params.Set("type", "JSON")

	fullURL := fmt.Sprintf("%s?%s", c.detailURL, params.Encode())
	logger.Debug("Precedent Detail API Request URL: %s", fullURL)

	// Perform request with retries
	body, err := c.doRequestWithRetry(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	// Debug: Log the raw response
	maxLen := 500
	if len(body) < maxLen {
		maxLen = len(body)
	}
	logger.Debug("Precedent Detail API Response (first %d chars): %s", maxLen, string(body[:maxLen]))

	// Parse response
	var detail LawDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		logger.Error("JSON parsing failed for precedent detail: %v", err)
		return nil, fmt.Errorf("판례 상세 정보 파싱 실패: %w", err)
	}

	return &detail, nil
}

// GetHistory retrieves precedent history (판례는 이력이 없으므로 미지원)
func (c *PrecClient) GetHistory(ctx context.Context, precID string) (*LawHistory, error) {
	return nil, fmt.Errorf("판례는 이력 조회를 지원하지 않습니다")
}

// doRequestWithRetry performs HTTP request with retry logic
func (c *PrecClient) doRequestWithRetry(ctx context.Context, url string) ([]byte, error) {
	var lastErr error
	delay := c.retryBaseDelay

	for i := 0; i < MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		body, err := c.doRequest(ctx, url)
		if err == nil {
			// Check for API error in response
			if c.hasAPIError(body) {
				return nil, c.parseAPIError(body)
			}
			return body, nil
		}

		lastErr = err
		if !c.shouldRetry(err) {
			return nil, err
		}

		if i < MaxRetries-1 {
			logger.Debug("Retrying after %v (attempt %d/%d)", delay, i+1, MaxRetries)
			time.Sleep(delay)
			delay *= 2
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// doRequest performs a single HTTP request
func (c *PrecClient) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTML error response
	if strings.HasPrefix(strings.TrimSpace(string(body)), "<!DOCTYPE") ||
		strings.HasPrefix(strings.TrimSpace(string(body)), "<html") {
		errMsg := c.parseHTMLError(string(body))
		return nil, &APIKeyError{Message: errMsg}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleHTTPError(resp.StatusCode)
	}

	return body, nil
}

// handleHTTPError handles HTTP status errors
func (c *PrecClient) handleHTTPError(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return &APIKeyError{Message: "API 인증 실패: API 키가 유효하지 않거나 만료되었습니다"}
	case http.StatusForbidden:
		return &APIKeyError{Message: "API 접근 권한이 없습니다"}
	case http.StatusNotFound:
		return fmt.Errorf("요청한 판례를 찾을 수 없습니다")
	case http.StatusTooManyRequests:
		return fmt.Errorf("API 요청 한도를 초과했습니다. 잠시 후 다시 시도해주세요")
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return fmt.Errorf("서버 오류가 발생했습니다. 잠시 후 다시 시도해주세요")
	default:
		return fmt.Errorf("HTTP 오류: %d", statusCode)
	}
}

// shouldRetry determines if the error is retryable
func (c *PrecClient) shouldRetry(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "서버 오류") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused")
}

// hasAPIError checks if the response contains an API error
func (c *PrecClient) hasAPIError(body []byte) bool {
	// Check for common error patterns in JSON/XML responses
	bodyStr := string(body)
	return strings.Contains(bodyStr, "\"errorCode\"") ||
		strings.Contains(bodyStr, "<errorCode>")
}

// parseAPIError parses API error from response body
func (c *PrecClient) parseAPIError(body []byte) error {
	// Try JSON first
	var jsonErr struct {
		ErrorCode    string `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal(body, &jsonErr); err == nil && jsonErr.ErrorCode != "" {
		if jsonErr.ErrorCode == "AUTH_ERROR" || strings.Contains(jsonErr.ErrorMessage, "인증") {
			return &APIKeyError{Message: fmt.Sprintf("API 인증 오류: %s", jsonErr.ErrorMessage)}
		}
		return fmt.Errorf("API 오류 [%s]: %s", jsonErr.ErrorCode, jsonErr.ErrorMessage)
	}

	// Try XML
	var xmlErr struct {
		ErrorCode    string `xml:"errorCode"`
		ErrorMessage string `xml:"errorMessage"`
	}
	if err := xml.Unmarshal(body, &xmlErr); err == nil && xmlErr.ErrorCode != "" {
		if xmlErr.ErrorCode == "AUTH_ERROR" || strings.Contains(xmlErr.ErrorMessage, "인증") {
			return &APIKeyError{Message: fmt.Sprintf("API 인증 오류: %s", xmlErr.ErrorMessage)}
		}
		return fmt.Errorf("API 오류 [%s]: %s", xmlErr.ErrorCode, xmlErr.ErrorMessage)
	}

	return fmt.Errorf("알 수 없는 API 오류")
}

// parseHTMLError extracts meaningful error message from HTML response
func (c *PrecClient) parseHTMLError(html string) string {
	htmlLower := strings.ToLower(html)

	// Check for authentication/key related issues
	if strings.Contains(htmlLower, "인증") || strings.Contains(htmlLower, "auth") ||
		strings.Contains(htmlLower, "key") || strings.Contains(htmlLower, "키") {
		return "API 인증 실패: API 키가 유효하지 않거나 만료되었습니다. 'warp config set law.key <API_KEY>' 명령으로 유효한 API 키를 설정해주세요."
	}

	// Check for service errors
	if strings.Contains(htmlLower, "서비스") || strings.Contains(htmlLower, "service") {
		return "서비스 일시 중단: 국가법령정보센터 서비스가 일시적으로 이용할 수 없습니다. 잠시 후 다시 시도해주세요."
	}

	// Check for not found errors
	if strings.Contains(htmlLower, "not found") || strings.Contains(htmlLower, "404") {
		return "요청한 판례를 찾을 수 없습니다."
	}

	// Default error message
	return "API 요청 실패: 국가법령정보센터 API에서 오류가 발생했습니다. API 키를 확인하거나 잠시 후 다시 시도해주세요."
}
