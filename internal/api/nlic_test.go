package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNLICClient_Search(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		responseStatus int
		expectedError  bool
		expectedCount  int
	}{
		{
			name: "successful search",
			responseBody: `{
				"totalCnt": 2,
				"page": 1,
				"law": [
					{"법령ID": "001", "법령명한글": "테스트법1"},
					{"법령ID": "002", "법령명한글": "테스트법2"}
				]
			}`,
			responseStatus: http.StatusOK,
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name:           "HTML error response",
			responseBody:   `<html><body><h1>Error</h1><p>인증 실패</p></body></html>`,
			responseStatus: http.StatusOK,
			expectedError:  true,
			expectedCount:  0,
		},
		{
			name:           "invalid JSON response",
			responseBody:   `{invalid json}`,
			responseStatus: http.StatusOK,
			expectedError:  true,
			expectedCount:  0,
		},
		{
			name:           "server error",
			responseBody:   "",
			responseStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := &NLICClient{
				httpClient:     &http.Client{Timeout: 5 * time.Second},
				baseURL:        server.URL,
				detailURL:      server.URL,
				apiKey:         "test-key",
				retryBaseDelay: time.Millisecond,
			}

			ctx := context.Background()
			req := &UnifiedSearchRequest{
				Query:    "test",
				PageNo:   1,
				PageSize: 10,
			}
			result, err := client.Search(ctx, req)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result.Laws) != tt.expectedCount {
					t.Errorf("expected %d items, got %d", tt.expectedCount, len(result.Laws))
				}
			}
		})
	}
}

func TestNLICClient_GetDetail(t *testing.T) {
	// Mock server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check API key
		if r.URL.Query().Get("OC") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check law ID (MST parameter for detail API)
		lawID := r.URL.Query().Get("MST")
		if lawID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return mock response based on law ID
		if lawID == "001234" {
			detail := LawDetail{
				LawInfo: LawInfo{
					ID:         "001234",
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
				Articles: []Article{
					{
						Number:  "제1조",
						Title:   "목적",
						Content: "이 법은 개인정보의 처리 및 보호에 관한 사항을 정함으로써...",
					},
					{
						Number:  "제2조",
						Title:   "정의",
						Content: "이 법에서 사용하는 용어의 뜻은 다음과 같다...",
					},
				},
				RelatedLaws: []string{"정보통신망법", "신용정보법"},
				Attachments: []string{"시행령.pdf", "시행규칙.pdf"},
			}

			// Wrap the response as the API expects
			response := map[string]interface{}{
				"법령": detail,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if lawID == "999999" {
			// Simulate not found
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &NLICClient{
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		baseURL:        server.URL,
		detailURL:      server.URL,
		historyURL:     server.URL,
		apiKey:         "test-api-key",
		retryBaseDelay: 100 * time.Millisecond,
	}

	tests := []struct {
		name         string
		lawID        string
		wantErr      bool
		wantName     string
		wantArticles int
	}{
		{
			name:         "Valid law ID",
			lawID:        "001234",
			wantErr:      false,
			wantName:     "개인정보 보호법",
			wantArticles: 2,
		},
		{
			name:         "Non-existent law ID",
			lawID:        "999999",
			wantErr:      true,
			wantName:     "",
			wantArticles: 0,
		},
		{
			name:         "Empty law ID",
			lawID:        "",
			wantErr:      true,
			wantName:     "",
			wantArticles: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := client.GetDetail(ctx, tt.lawID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Name != tt.wantName {
					t.Errorf("GetDetail() Name = %v, want %v", result.Name, tt.wantName)
				}
				if len(result.Articles) != tt.wantArticles {
					t.Errorf("GetDetail() Articles count = %v, want %v", len(result.Articles), tt.wantArticles)
				}
			}
		})
	}
}

func TestNLICClient_GetHistory(t *testing.T) {
	// Mock server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check API key
		if r.URL.Query().Get("OC") != "test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check law ID (MST parameter for history API)
		lawID := r.URL.Query().Get("MST")
		if lawID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return mock response based on law ID
		if lawID == "001234" {
			response := LawHistory{
				LawID:   "001234",
				LawName: "개인정보 보호법",
				Histories: []HistoryRecord{
					{
						Date:       "20110329",
						Type:       "제정",
						PromulNo:   "제10465호",
						EffectDate: "20110930",
						Reason:     "개인정보 보호 강화를 위한 법률 제정",
					},
					{
						Date:       "20200205",
						Type:       "일부개정",
						PromulNo:   "제16930호",
						EffectDate: "20200805",
						Reason:     "개인정보 처리자의 책임성 강화",
					},
					{
						Date:       "20230314",
						Type:       "일부개정",
						PromulNo:   "제19234호",
						EffectDate: "20230914",
						Reason:     "디지털 환경 변화에 따른 보호 강화",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if lawID == "999999" {
			// Simulate not found
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &NLICClient{
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		baseURL:        server.URL,
		detailURL:      server.URL,
		historyURL:     server.URL,
		apiKey:         "test-api-key",
		retryBaseDelay: 100 * time.Millisecond,
	}

	tests := []struct {
		name          string
		lawID         string
		wantErr       bool
		wantName      string
		wantHistories int
	}{
		{
			name:          "Valid law ID",
			lawID:         "001234",
			wantErr:       false,
			wantName:      "개인정보 보호법",
			wantHistories: 3,
		},
		{
			name:          "Non-existent law ID",
			lawID:         "999999",
			wantErr:       true,
			wantName:      "",
			wantHistories: 0,
		},
		{
			name:          "Empty law ID",
			lawID:         "",
			wantErr:       true,
			wantName:      "",
			wantHistories: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := client.GetHistory(ctx, tt.lawID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.LawName != tt.wantName {
					t.Errorf("GetHistory() LawName = %v, want %v", result.LawName, tt.wantName)
				}
				if len(result.Histories) != tt.wantHistories {
					t.Errorf("GetHistory() Histories count = %v, want %v", len(result.Histories), tt.wantHistories)
				}
			}
		})
	}
}

func TestNLICClient_RetryLogic(t *testing.T) {
	retryCount := 0
	maxRetries := 3

	// Mock server that simulates temporary failures
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++

		// Fail the first attempts, succeed on the last
		if retryCount < maxRetries {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// Success response
		response := SearchResponse{
			TotalCount: 1,
			Page:       1,
			Laws: []LawInfo{
				{
					ID:   "001",
					Name: "Test Law",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test server URL and short retry delay
	client := &NLICClient{
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		baseURL:        server.URL,
		detailURL:      server.URL,
		historyURL:     server.URL,
		apiKey:         "test-api-key",
		retryBaseDelay: 10 * time.Millisecond, // Short delay for testing
	}

	ctx := context.Background()
	req := &UnifiedSearchRequest{
		Query:    "test",
		PageNo:   1,
		PageSize: 10,
	}

	// Should succeed after retries
	result, err := client.Search(ctx, req)
	if err != nil {
		t.Errorf("Search() should succeed after retries, got error: %v", err)
	}

	if result == nil || result.TotalCount != 1 {
		t.Errorf("Search() unexpected result: %v", result)
	}

	if retryCount != maxRetries {
		t.Errorf("Expected %d retries, got %d", maxRetries, retryCount)
	}
}
