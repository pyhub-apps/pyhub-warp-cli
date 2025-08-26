package api

import (
	"context"
	"fmt"

	"github.com/pyhub-apps/pyhub-warp-cli/internal/config"
)

// CreateClient creates an API client for the specified type
func CreateClient(apiType APIType) (ClientInterface, error) {
	switch apiType {
	case APITypeNLIC:
		apiKey := config.GetNLICAPIKey()
		if apiKey == "" {
			return nil, fmt.Errorf("NLIC API 키가 설정되지 않았습니다. 'warp config set law.nlic.key YOUR_KEY' 명령으로 설정하세요")
		}
		// Use the dedicated NLIC client
		return NewNLICClient(apiKey), nil

	case APITypeELIS:
		// ELIS uses the same API key as NLIC (both from law.go.kr)
		apiKey := config.GetNLICAPIKey()
		if apiKey == "" {
			// Try ELIS-specific key first
			apiKey = config.GetString("law.elis.key")
			if apiKey == "" {
				return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 'warp config set law.key YOUR_KEY' 명령으로 설정하세요")
			}
		}
		return NewELISClient(apiKey), nil

	case APITypeAll:
		// Unified client for searching both NLIC and ELIS
		return NewUnifiedClient()

	case APITypePrec:
		// Precedent API client (판례)
		apiKey := config.GetNLICAPIKey()
		if apiKey == "" {
			return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 'warp config set law.key YOUR_KEY' 명령으로 설정하세요")
		}
		return NewPrecClient(apiKey), nil

	case APITypeAdmrul:
		// Administrative Rule API client (행정규칙)
		apiKey := config.GetNLICAPIKey()
		if apiKey == "" {
			return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 'warp config set law.key YOUR_KEY' 명령으로 설정하세요")
		}
		return NewAdmrulClient(apiKey), nil

	case APITypeExpc:
		// Legal Interpretation API client (법령해석례)
		apiKey := config.GetNLICAPIKey()
		if apiKey == "" {
			return nil, fmt.Errorf("API 키가 설정되지 않았습니다. 'warp config set law.key YOUR_KEY' 명령으로 설정하세요")
		}
		return NewExpcClient(apiKey), nil

	default:
		return nil, fmt.Errorf("알 수 없는 API 타입: %s", apiType)
	}
}

// CreateDefaultClient creates a client using the default (NLIC) API
func CreateDefaultClient() (ClientInterface, error) {
	return CreateClient(APITypeNLIC)
}

// LegacyClientWrapper wraps the existing Client to implement ClientInterface
type LegacyClientWrapper struct {
	*Client
}

// Search adapts the request to the legacy client
func (w *LegacyClientWrapper) Search(ctx context.Context, req *UnifiedSearchRequest) (*SearchResponse, error) {
	// Convert UnifiedSearchRequest to legacy SearchRequest
	legacyReq := &SearchRequest{
		Query:    req.Query,
		Type:     req.Type,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}
	return w.Client.Search(ctx, legacyReq)
}

// GetDetail is not yet implemented in the legacy client
func (w *LegacyClientWrapper) GetDetail(ctx context.Context, lawID string) (*LawDetail, error) {
	return nil, fmt.Errorf("상세 조회 기능은 아직 구현되지 않았습니다")
}

// GetHistory is not yet implemented in the legacy client
func (w *LegacyClientWrapper) GetHistory(ctx context.Context, lawID string) (*LawHistory, error) {
	return nil, fmt.Errorf("이력 조회 기능은 아직 구현되지 않았습니다")
}

// GetAPIType returns the API type
func (w *LegacyClientWrapper) GetAPIType() APIType {
	return APITypeNLIC
}
