# 📋 Sejong CLI - 다음 단계 작업 가이드

## 🎯 현재 상황
- **Phase 1 완료**: 핵심 법령 서비스 구현 완성
- **구현된 API**: 13개 (법령 3, 자치법규 4, 판례 2, 행정규칙 2, 법령해석례 2)
- **구현률**: 10.4% (13/125)

## 🔴 우선순위: 높음 (즉시 필요)

### 1. 실제 API 테스트 및 버그 수정
```bash
# 실제 API 키 획득
https://open.law.go.kr/LSO/openApi/cuAskList.do

# 테스트 명령어
./sejong law search "개인정보"
./sejong precedent search "계약"
./sejong admrule search "공공기관"
./sejong interpretation search "근로시간"
```

**예상 이슈**:
- API 응답 형식 불일치 (JSON/XML)
- 필드명 차이
- 페이지네이션 파라미터 차이

### 2. 테스트 코드 작성
```go
// 우선 작성할 테스트 파일들
internal/api/nlic_test.go
internal/api/prec_test.go
internal/api/admrul_test.go
internal/api/expc_test.go
internal/cmd/law_test.go
```

## 🟡 우선순위: 중간 (기능 확장)

### 3. 통합 검색 기능
```bash
# 새 명령어 구조
sejong search "키워드" [flags]
  --type all|law|prec|admrul|expc  # 검색 대상 (기본: all)
  --format table|json               # 출력 형식
  --page 1                          # 페이지 번호
  --size 10                         # 페이지 크기
```

**구현 파일**:
- `internal/cmd/search.go` - 통합 검색 명령어
- `internal/api/unified_search.go` - 병렬 검색 로직

### 4. Phase 2: 헌재결정례 & 조약
```go
// internal/api/types.go
APITypeConst APIType = "const"   // 헌재결정례
APITypeTreaty APIType = "treaty" // 조약

// 새 클라이언트
internal/api/const.go
internal/api/treaty.go

// 새 명령어
internal/cmd/constitutional.go
internal/cmd/treaty.go
```

## 🟢 우선순위: 낮음 (품질 개선)

### 5. 캐싱 시스템
```go
// internal/cache/cache.go
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
}

// 캐시 키 형식
// search:{api_type}:{query}:{page}:{size}
// detail:{api_type}:{id}
```

### 6. 설정 프로파일
```yaml
# ~/.sejong/config.yaml 구조 개선
profiles:
  default:
    law:
      key: "API_KEY_1"
  work:
    law:
      key: "API_KEY_2"
active_profile: default
```

## 🐛 발견된 이슈

### Issue #1: API 응답 형식
- **문제**: XML/JSON 자동 감지 없음
- **해결**: Content-Type 헤더 확인 및 자동 파싱

### Issue #2: 페이지네이션 정보
- **현재**: `총 3개의 법령을 찾았습니다.`
- **개선**: `총 123개 중 1-10번째 결과 (1/13 페이지)`

### Issue #3: 에러 메시지 일관성
- **문제**: API 키 오류 시 처리 불일치
- **해결**: 통일된 에러 핸들러 구현

## 📅 추천 작업 순서

### Day 1-2 (즉시)
- [ ] 실제 API 키 획득
- [ ] 각 API 타입별 실제 테스트
- [ ] 발견된 파싱 오류 수정
- [ ] 기본 단위 테스트 작성

### Week 1 (단기)
- [ ] 통합 검색 기능 구현
- [ ] 페이지네이션 UI 개선
- [ ] 에러 처리 일관성 확보
- [ ] README 업데이트

### Week 2 (중기)
- [ ] Phase 2 시작 (헌재결정례)
- [ ] 조약 API 구현
- [ ] 캐싱 시스템 기본 구현
- [ ] 성능 측정 및 최적화

### Month 1 (장기)
- [ ] 영문 법령 API 추가
- [ ] 법령 비교 기능
- [ ] 설정 프로파일 시스템
- [ ] CI/CD 파이프라인 구축

## 🚀 빠른 시작 명령어

```bash
# 빌드
go build -o sejong cmd/sejong/main.go

# 테스트
go test ./...

# 실행
./sejong law search "test"

# 형식 검사
go fmt ./...
go vet ./...

# 릴리즈 빌드
goreleaser release --snapshot --clean
```

## 📝 메모

- 현재 구현된 모든 API는 `target` 파라미터로 구분됨
  - law (법령), prec (판례), admrul (행정규칙), expc (법령해석례)
- Base URL: `https://www.law.go.kr/DRF/lawSearch.do`
- 모든 API는 동일한 인증 방식 사용 (OC 파라미터)
- 에러 응답이 HTML로 오는 경우 처리 로직 구현됨

## 🔗 참고 자료

- API 가이드: https://open.law.go.kr/LSO/openApi/guideList.do
- API 키 발급: https://open.law.go.kr/LSO/openApi/cuAskList.do
- 공공데이터포털: https://www.data.go.kr/data/15000115/openapi.do

---

*이 문서는 작업 진행을 위한 임시 가이드입니다. `.gitignore`에 추가하여 커밋하지 마세요.*