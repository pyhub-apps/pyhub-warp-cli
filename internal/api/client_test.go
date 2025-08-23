package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// Test with API key
	client := NewClientWithConfig("test-api-key", 10*time.Second)
	if client == nil {
		t.Error("Expected client to be created")
	}
	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key to be 'test-api-key', got %s", client.apiKey)
	}
}

func TestClient_Search(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check required parameters
		query := r.URL.Query()
		if query.Get("OC") == "" {
			t.Error("Missing API key parameter")
		}
		if query.Get("target") != "law" {
			t.Errorf("Expected target to be 'law', got %s", query.Get("target"))
		}
		if query.Get("query") == "" {
			t.Error("Missing query parameter")
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"totalCnt": 2,
			"page": 1,
			"law": [
				{
					"법령ID": "001234",
					"법령명한글": "개인정보 보호법",
					"법령명약칭": "개인정보보호법",
					"법령일련번호": "017969",
					"공포일자": "20200205",
					"공포번호": "16930",
					"제개정구분명": "일부개정",
					"소관부처명": "개인정보보호위원회",
					"시행일자": "20200805",
					"법령구분명": "법률"
				},
				{
					"법령ID": "001235",
					"법령명한글": "개인정보 보호법 시행령",
					"법령명약칭": "개인정보보호법시행령",
					"법령일련번호": "031380",
					"공포일자": "20210202",
					"공포번호": "31380",
					"제개정구분명": "일부개정",
					"소관부처명": "개인정보보호위원회",
					"시행일자": "20210202",
					"법령구분명": "대통령령"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithConfig("test-api-key", 10*time.Second)
	client.baseURL = server.URL

	// Test search
	ctx := context.Background()
	req := &SearchRequest{
		Query:    "개인정보",
		Type:     "JSON",
		PageNo:   1,
		PageSize: 10,
	}

	resp, err := client.Search(ctx, req)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Verify response
	if resp.TotalCount != 2 {
		t.Errorf("Expected total count to be 2, got %d", resp.TotalCount)
	}
	if len(resp.Laws) != 2 {
		t.Errorf("Expected 2 laws, got %d", len(resp.Laws))
	}
	if resp.Laws[0].Name != "개인정보 보호법" {
		t.Errorf("Expected first law name to be '개인정보 보호법', got %s", resp.Laws[0].Name)
	}
}

func TestClient_SearchWithRetry(t *testing.T) {
	attempts := 0
	// Create a mock server that fails first 2 attempts
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Simulate server error (retryable)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Success on third attempt
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"totalCnt": 0, "page": 1, "law": []}`))
	}))
	defer server.Close()

	// Create client with short retry delay for faster testing
	client := NewClientWithConfig("test-api-key", 10*time.Second)
	client.baseURL = server.URL
	client.retryBaseDelay = 5 * time.Millisecond // Speed up test

	// Test search with retry
	ctx := context.Background()
	req := &SearchRequest{
		Query: "test",
		Type:  "JSON",
	}

	resp, err := client.Search(ctx, req)
	if err != nil {
		t.Fatalf("Search failed after retries: %v", err)
	}

	// Verify exactly 3 attempts (initial + 2 retries)
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
	if resp.TotalCount != 0 {
		t.Errorf("Expected total count to be 0, got %d", resp.TotalCount)
	}
}

func TestClient_SearchTimeout(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with short timeout
	client := NewClientWithConfig("test-api-key", 100*time.Millisecond)
	client.baseURL = server.URL
	client.retryBaseDelay = 5 * time.Millisecond // Speed up test

	// Test search with timeout
	ctx := context.Background()
	req := &SearchRequest{
		Query: "test",
	}

	_, err := client.Search(ctx, req)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Check if it's a timeout error
	if !isTimeoutError(err) {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// isTimeoutError checks if an error is a timeout error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	// Check for net.Error timeout
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	// Check for context deadline exceeded
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	return false
}

func TestClient_SearchWithCancel(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	client := NewClientWithConfig("test-api-key", 10*time.Second)
	client.baseURL = server.URL

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	// Test search with cancelled context
	req := &SearchRequest{
		Query: "test",
	}

	_, err := client.Search(ctx, req)
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}

	// Check if it's a context canceled error
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestClient_SearchNoRetryOn4xx(t *testing.T) {
	attempts := 0
	// Create a mock server that returns 400 Bad Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	// Create client with short retry delay
	client := NewClientWithConfig("test-api-key", 10*time.Second)
	client.baseURL = server.URL
	client.retryBaseDelay = 5 * time.Millisecond

	// Test search
	ctx := context.Background()
	req := &SearchRequest{
		Query: "test",
	}

	_, err := client.Search(ctx, req)
	if err == nil {
		t.Fatal("Expected error for 400 status")
	}

	// Should not retry on 4xx errors
	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retry on 4xx), got %d", attempts)
	}
}

func TestClient_SearchXMLResponse(t *testing.T) {
	// Create a mock server that returns XML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<SearchResponse>
			<totalCnt>1</totalCnt>
			<page>1</page>
			<law>
				<법령ID>001234</법령ID>
				<법령명한글>테스트 법령</법령명한글>
				<법령명약칭>테스트법</법령명약칭>
				<법령일련번호>12345</법령일련번호>
				<공포일자>20231201</공포일자>
				<공포번호>1234</공포번호>
				<제개정구분명>제정</제개정구분명>
				<소관부처명>테스트부</소관부처명>
				<시행일자>20240101</시행일자>
				<법령구분명>법률</법령구분명>
			</law>
		</SearchResponse>`))
	}))
	defer server.Close()

	// Create client
	client := NewClientWithConfig("test-api-key", 10*time.Second)
	client.baseURL = server.URL

	// Test search with XML response
	ctx := context.Background()
	req := &SearchRequest{
		Query: "test",
		Type:  TypeXML,
	}

	resp, err := client.Search(ctx, req)
	if err != nil {
		t.Fatalf("Failed to parse XML response: %v", err)
	}

	// Verify response
	if resp.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", resp.TotalCount)
	}
	if len(resp.Laws) != 1 {
		t.Errorf("Expected 1 law, got %d", len(resp.Laws))
	}
	if resp.Laws[0].Name != "테스트 법령" {
		t.Errorf("Expected law name '테스트 법령', got '%s'", resp.Laws[0].Name)
	}
}

func TestClient_SearchEmptyQuery(t *testing.T) {
	// Create client
	client := NewClientWithConfig("test-api-key", 10*time.Second)

	// Test with empty query
	ctx := context.Background()
	req := &SearchRequest{
		Query: "",
	}

	// For now, empty query is allowed (API will handle it)
	// In future, we might want to add validation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query != "" {
			t.Errorf("Expected empty query, got '%s'", query)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"totalCnt": 0, "page": 1, "law": []}`))
	}))
	defer server.Close()

	client.baseURL = server.URL
	resp, err := client.Search(ctx, req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.TotalCount != 0 {
		t.Errorf("Expected 0 results for empty query, got %d", resp.TotalCount)
	}
}
