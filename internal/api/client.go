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

	"github.com/pyhub-kr/pyhub-sejong-cli/internal/config"
)

const (
	// BaseURL is the National Law Information Center API endpoint
	BaseURL = "https://www.law.go.kr/DRF/lawSearch.do"
	
	// Default timeout for HTTP requests
	DefaultTimeout = 10 * time.Second
	
	// Maximum retry attempts
	MaxRetries = 3
	
	// Initial retry delay
	InitialRetryDelay = 1 * time.Second
	
	// Supported response types
	TypeJSON = "JSON"
	TypeXML  = "XML"
	
	// Default target for API requests
	DefaultTarget = "law"
)

// Client represents the API client for National Law Information Center
type Client struct {
	httpClient     *http.Client
	baseURL        string
	apiKey         string
	retryBaseDelay time.Duration // Configurable retry delay for testing
}

// SearchRequest represents the search request parameters
type SearchRequest struct {
	Query      string `json:"query"`
	Type       string `json:"type"`       // "XML" or "JSON"
	PageNo     int    `json:"page_no"`
	PageSize   int    `json:"page_size"`
}

// SearchResponse represents the search response
type SearchResponse struct {
	TotalCount int         `json:"totalCnt" xml:"totalCnt"`
	Page       int         `json:"page" xml:"page"`
	Laws       []LawInfo   `json:"law" xml:"law"`
	Error      *ErrorInfo  `json:"error,omitempty" xml:"error,omitempty"`
}

// LawInfo represents individual law information
type LawInfo struct {
	ID           string `json:"법령ID" xml:"법령ID"`
	Name         string `json:"법령명한글" xml:"법령명한글"`
	NameAbbrev   string `json:"법령명약칭" xml:"법령명약칭"`
	SerialNo     string `json:"법령일련번호" xml:"법령일련번호"`
	PromulDate   string `json:"공포일자" xml:"공포일자"`
	PromulNo     string `json:"공포번호" xml:"공포번호"`
	Category     string `json:"제개정구분명" xml:"제개정구분명"`
	Department   string `json:"소관부처명" xml:"소관부처명"`
	EffectDate   string `json:"시행일자" xml:"시행일자"`
	LawType      string `json:"법령구분명" xml:"법령구분명"`
}

// ErrorInfo represents API error information
type ErrorInfo struct {
	Code    string `json:"code" xml:"code"`
	Message string `json:"message" xml:"message"`
}

// RetryableError represents an error that can be retried
type RetryableError struct {
	err error
}

func (e *RetryableError) Error() string {
	return e.err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.err
}

// NewClient creates a new API client
func NewClient() (*Client, error) {
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 'sejong config set law.key YOUR_KEY' 명령으로 설정하세요")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:        BaseURL,
		apiKey:         apiKey,
		retryBaseDelay: InitialRetryDelay,
	}, nil
}

// NewClientWithConfig creates a new API client with custom configuration
func NewClientWithConfig(apiKey string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL:        BaseURL,
		apiKey:         apiKey,
		retryBaseDelay: InitialRetryDelay,
	}
}

// Search performs a law search with the given query
func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// Set defaults
	if req.Type == "" {
		req.Type = TypeJSON
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
	params.Set("target", DefaultTarget)
	params.Set("query", req.Query)
	params.Set("type", req.Type)
	params.Set("page", fmt.Sprintf("%d", req.PageNo))
	params.Set("display", fmt.Sprintf("%d", req.PageSize))

	fullURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Perform request with retries
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

		resp, err := c.doRequest(ctx, fullURL)
		if err != nil {
			lastErr = err
			// Only retry on network errors or 5xx server errors
			if !c.shouldRetry(err) {
				return nil, err
			}
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("요청 실패 (재시도 %d회 초과): %w", MaxRetries, lastErr)
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(ctx context.Context, url string) (*SearchResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RetryableError{fmt.Errorf("네트워크 에러: %w", err)}
	}
	defer resp.Body.Close()

	// Check HTTP status first
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= 500 {
			// Server error - retryable
			return nil, &RetryableError{fmt.Errorf("서버 에러: HTTP %d", resp.StatusCode)}
		}
		// Client error - not retryable
		return nil, fmt.Errorf("클라이언트 에러: HTTP %d", resp.StatusCode)
	}

	// Read response body only if status is OK
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", err)
	}

	// Parse response based on content type
	var searchResp SearchResponse
	contentType := resp.Header.Get("Content-Type")
	
	if strings.Contains(contentType, "json") {
		if err := json.Unmarshal(body, &searchResp); err != nil {
			return nil, fmt.Errorf("JSON 파싱 실패: %w", err)
		}
	} else if strings.Contains(contentType, "xml") {
		if err := xml.Unmarshal(body, &searchResp); err != nil {
			return nil, fmt.Errorf("XML 파싱 실패: %w", err)
		}
	} else {
		// Try JSON first, then XML
		if err := json.Unmarshal(body, &searchResp); err != nil {
			if err := xml.Unmarshal(body, &searchResp); err != nil {
				return nil, fmt.Errorf("응답 파싱 실패 (JSON/XML 형식이 아님)")
			}
		}
	}

	// Check for API errors
	if searchResp.Error != nil {
		return nil, fmt.Errorf("API 에러: %s - %s", searchResp.Error.Code, searchResp.Error.Message)
	}

	return &searchResp, nil
}

// shouldRetry determines if an error is retryable
func (c *Client) shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a RetryableError
	var retryableErr *RetryableError
	if errors.As(err, &retryableErr) {
		return true
	}
	
	return false
}