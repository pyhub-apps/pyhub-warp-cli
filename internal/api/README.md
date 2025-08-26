# API Client Module

국가법령정보센터 Open API 클라이언트 모듈입니다.

## 주요 기능

- HTTP 클라이언트 초기화 및 관리
- API 요청/응답 처리
- 자동 재시도 로직 (Exponential Backoff)
- 타임아웃 처리
- JSON/XML 응답 파싱
- 컨텍스트 기반 취소 지원

## 사용 예시

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/pyhub-apps/pyhub-warp-cli/internal/api"
)

func main() {
    // API 클라이언트 생성
    client, err := api.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    // 검색 요청
    req := &api.SearchRequest{
        Query:    "개인정보 보호법",
        Type:     "JSON",
        PageNo:   1,
        PageSize: 10,
    }
    
    // 검색 실행
    ctx := context.Background()
    resp, err := client.Search(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    // 결과 출력
    fmt.Printf("총 %d개의 법령을 찾았습니다.\n", resp.TotalCount)
    for _, law := range resp.Laws {
        fmt.Printf("- %s (%s)\n", law.Name, law.Department)
    }
}
```

## 에러 처리

API 클라이언트는 다음과 같은 에러를 반환할 수 있습니다:

- **네트워크 에러**: 서버에 연결할 수 없는 경우
- **타임아웃 에러**: 요청이 10초 내에 완료되지 않는 경우
- **파싱 에러**: 응답이 유효한 JSON/XML 형식이 아닌 경우
- **API 에러**: 서버에서 에러 응답을 반환한 경우
- **설정 에러**: API 키가 설정되지 않은 경우

## 재시도 로직

- 최대 3회 재시도
- Exponential Backoff 적용 (1초, 2초, 4초)
- 네트워크 에러 및 5xx 서버 에러 시 재시도

## 테스트

```bash
go test ./internal/api/... -v
```