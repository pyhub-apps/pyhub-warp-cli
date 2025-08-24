package api

import (
	"context"
)

// APIType represents the type of law API
type APIType string

const (
	// APITypeNLIC represents National Law Information Center API (국가법령정보센터)
	APITypeNLIC APIType = "nlic"
	// APITypeELIS represents Local Regulations Information System API (자치법규정보시스템)
	APITypeELIS APIType = "elis"
	// APITypeAll represents all APIs combined (통합 검색)
	APITypeAll APIType = "all"
	// APITypePrec represents Precedent API (판례)
	APITypePrec APIType = "prec"
	// APITypeAdmrul represents Administrative Rule API (행정규칙)
	APITypeAdmrul APIType = "admrul"
	// APITypeExpc represents Legal Interpretation API (법령해석례)
	APITypeExpc APIType = "expc"
)

// ClientInterface represents a unified API client interface for law information
type ClientInterface interface {
	// Search performs a law search
	Search(ctx context.Context, req *UnifiedSearchRequest) (*SearchResponse, error)
	// GetDetail retrieves detailed law information
	GetDetail(ctx context.Context, lawID string) (*LawDetail, error)
	// GetHistory retrieves law amendment history
	GetHistory(ctx context.Context, lawID string) (*LawHistory, error)
	// GetAPIType returns the API type
	GetAPIType() APIType
}

// UnifiedSearchRequest represents a unified search request for multi-API support
type UnifiedSearchRequest struct {
	Query      string            // Search query
	PageNo     int               // Page number (1-based)
	PageSize   int               // Results per page
	Type       string            // Response type (JSON/XML)
	Region     string            // Region filter (for ELIS)
	LawType    string            // Law type filter
	Department string            // Department filter
	DateFrom   string            // Date range start (YYYYMMDD)
	DateTo     string            // Date range end (YYYYMMDD)
	Sort       string            // Sort order
	Extras     map[string]string // API-specific extra parameters
}

// LawDetail represents detailed law information
type LawDetail struct {
	LawInfo
	Content     string    `json:"조문내용" xml:"조문내용"`
	Articles    []Article `json:"조문" xml:"조문"`
	Attachments []string  `json:"첨부파일" xml:"첨부파일"`
	RelatedLaws []string  `json:"관련법령" xml:"관련법령"`
}

// Article represents a law article
type Article struct {
	Number     string `json:"조문번호" xml:"조문번호"`
	Title      string `json:"조문제목" xml:"조문제목"`
	Content    string `json:"조문내용" xml:"조문내용"`
	EffectDate string `json:"시행일자" xml:"시행일자"`
}

// LawHistory represents law amendment history
type LawHistory struct {
	LawID     string          `json:"법령ID" xml:"법령ID"`
	LawName   string          `json:"법령명" xml:"법령명"`
	Histories []HistoryRecord `json:"연혁" xml:"연혁"`
}

// HistoryRecord represents a single history record
type HistoryRecord struct {
	Date       string `json:"일자" xml:"일자"`
	Type       string `json:"구분" xml:"구분"` // 제정, 일부개정, 전부개정, 폐지
	Reason     string `json:"개정이유" xml:"개정이유"`
	PromulNo   string `json:"공포번호" xml:"공포번호"`
	EffectDate string `json:"시행일자" xml:"시행일자"`
}
