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
	Content                 string                   `json:"조문내용" xml:"조문내용"`
	Articles                []Article                `json:"조문" xml:"조문"`
	Attachments             []string                 `json:"첨부파일" xml:"첨부파일"`
	RelatedLaws             []string                 `json:"관련법령" xml:"관련법령"`
	RevisionText            string                   `json:"개정문" xml:"개정문"`     // 개정문 내용
	Tables                  []Table                  `json:"별표" xml:"별표"`       // 별표 목록
	SupplementaryProvisions []SupplementaryProvision `json:"부칙" xml:"부칙"`       // 부칙 목록
	HasRevisionText         bool                     `json:"개정문존재" xml:"개정문존재"` // 개정문 존재 여부
}

// LawDetailResponse represents the actual API response structure for law detail
type LawDetailResponse struct {
	Law LawDetailContent `json:"법령" xml:"법령"`
}

// LawDetailContent represents the content structure returned by the detail API
type LawDetailContent struct {
	LawKey                  string                         `json:"법령키" xml:"법령키"`
	BasicInfo               *BasicInfo                     `json:"기본정보" xml:"기본정보"`
	Revisions               RevisionContent                `json:"개정문" xml:"개정문"`
	Tables                  TableContent                   `json:"별표" xml:"별표"`
	ArticlesRaw             ArticlesContent                `json:"조문" xml:"조문"`
	SupplementaryProvisions SupplementaryProvisionsContent `json:"부칙" xml:"부칙"`
}

// BasicInfo represents basic law information
type BasicInfo struct {
	LawID              string         `json:"법령ID" xml:"법령ID"`
	LawNameKorean      string         `json:"법령명_한글" xml:"법령명_한글"`
	LawNameHanja       string         `json:"법령명_한자" xml:"법령명_한자"`
	LawNameAbbr        string         `json:"법령명약칭" xml:"법령명약칭"`
	PromulgationDate   string         `json:"공포일자" xml:"공포일자"`
	PromulgationNumber string         `json:"공포번호" xml:"공포번호"`
	EffectiveDate      string         `json:"시행일자" xml:"시행일자"`
	RevisionType       string         `json:"제개정구분" xml:"제개정구분"`
	Department         DepartmentInfo `json:"소관부처" xml:"소관부처"`
	LawTypeInfo        LawTypeInfo    `json:"법종구분" xml:"법종구분"`
}

// DepartmentInfo represents department information
type DepartmentInfo struct {
	Content string `json:"content" xml:"content"`
	Code    string `json:"소관부처코드" xml:"소관부처코드"`
}

// LawTypeInfo represents law type information
type LawTypeInfo struct {
	Content string `json:"content" xml:"content"`
	Code    string `json:"법종구분코드" xml:"법종구분코드"`
}

// RevisionContent represents revision content structure
type RevisionContent struct {
	Content interface{} `json:"개정문내용" xml:"개정문내용"` // Can be string or nested structure
}

// TableContent represents table content structure
type TableContent struct {
	TableUnits []TableUnit `json:"별표단위" xml:"별표단위"`
}

// TableUnit represents a single table unit from the API
type TableUnit struct {
	TableKey     string      `json:"별표키" xml:"별표키"`
	TableNumber  string      `json:"별표번호" xml:"별표번호"`
	TableTitle   string      `json:"별표제목" xml:"별표제목"`
	TableContent interface{} `json:"별표내용" xml:"별표내용"` // Can be string or array
	TableHoYN    string      `json:"별표서식여부" xml:"별표서식여부"`
}

// SupplementaryProvisionsContent represents supplementary provisions structure
type SupplementaryProvisionsContent struct {
	ProvisionUnits []SupplementaryProvisionUnit `json:"부칙단위" xml:"부칙단위"`
}

// SupplementaryProvisionUnit represents a single supplementary provision
type SupplementaryProvisionUnit struct {
	ProvisionKey     string      `json:"부칙키" xml:"부칙키"`
	ProvisionNumber  string      `json:"부칙번호" xml:"부칙번호"`
	ProvisionContent interface{} `json:"부칙내용" xml:"부칙내용"` // Can be string or array
	ProvisionDate    string      `json:"부칙공포일자" xml:"부칙공포일자"`
}

// Table represents a law table/appendix
type Table struct {
	Number  string `json:"번호" xml:"번호"`
	Title   string `json:"제목" xml:"제목"`
	Content string `json:"내용" xml:"내용"`
}

// ArticlesContent represents articles content structure
type ArticlesContent struct {
	ArticleUnits []ArticleUnit `json:"조문단위" xml:"조문단위"`
}

// ArticleUnit represents a single article unit from the API
type ArticleUnit struct {
	ArticleKey        string      `json:"조문키" xml:"조문키"`
	ArticleNumber     string      `json:"조문번호" xml:"조문번호"`
	ArticleYN         string      `json:"조문여부" xml:"조문여부"`
	ArticleContent    string      `json:"조문내용" xml:"조문내용"`
	ArticleReference  string      `json:"조문참고자료" xml:"조문참고자료"`
	ArticleEffectDate string      `json:"조문시행일자" xml:"조문시행일자"`
	ArticleTitle      string      `json:"조문제목" xml:"조문제목"`
	ArticleChangeYN   string      `json:"조문변경여부" xml:"조문변경여부"`
	ArticleMoveBefore string      `json:"조문이동이전" xml:"조문이동이전"`
	ArticleMoveAfter  string      `json:"조문이동이후" xml:"조문이동이후"`
	Paragraphs        interface{} `json:"항" xml:"항"` // Can be array or object
	ArticleHistory    interface{} `json:"조문이동이력" xml:"조문이동이력"`
	LawID             string      `json:"법령ID" xml:"법령ID"`
	LawNameKorean     string      `json:"법령명한글" xml:"법령명한글"`
	LawSerialNo       string      `json:"법령일련번호" xml:"법령일련번호"`
}

// SupplementaryProvision represents a supplementary provision (부칙)
type SupplementaryProvision struct {
	Number           string `json:"번호" xml:"번호"`
	Content          string `json:"내용" xml:"내용"`
	PromulgationDate string `json:"공포일자" xml:"공포일자"`
	PromulgationNo   string `json:"공포번호" xml:"공포번호"`
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
