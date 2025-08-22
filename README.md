# Sejong-CLI

> 🏛️ 터미널에서 빠르게 검색하는 대한민국 법령 정보

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/badge/Release-1.2534.6--dev-orange)](https://github.com/pyhub-kr/pyhub-sejong-cli/releases)

## 📖 소개

**Sejong-CLI**는 국가법령정보센터 오픈 API를 활용하여 터미널에서 한국 법령 정보를 빠르고 쉽게 검색할 수 있는 명령줄 도구입니다.

### ✨ 주요 기능

- 🔍 **빠른 법령 검색**: 터미널에서 즉시 법령 정보 조회
- 📋 **다양한 출력 형식**: 테이블 형식 또는 JSON 형식 지원
- ⚡ **간편한 설정**: 한 번의 API 키 설정으로 계속 사용

### 👥 이런 분들께 추천합니다

- 법률/규제 관련 서비스를 개발하는 **개발자**
- 법령 데이터를 분석하는 **연구원**
- 법령 정보를 자주 확인하는 **법률 전문가**

## 🚀 설치

### Go 빌드 (개발 버전)

```bash
# 저장소 클론
git clone https://github.com/pyhub-kr/pyhub-sejong-cli.git
cd pyhub-sejong-cli

# 빌드
go build -o sejong cmd/sejong/main.go

# 실행
./sejong --help
```

### 바이너리 다운로드 (예정)

향후 [Releases](https://github.com/pyhub-kr/pyhub-sejong-cli/releases) 페이지에서 각 OS별 바이너리를 다운로드할 수 있습니다.

## 🎯 빠른 시작

### 1. API 키 발급

국가법령정보센터에서 오픈 API 인증키를 발급받으세요:
👉 [https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9](https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9)

### 2. API 키 설정

```bash
sejong config set law.key YOUR_API_KEY
```

### 3. 첫 검색

```bash
# 법령 검색
sejong law "개인정보 보호법"

# JSON 형식으로 출력
sejong law "도로교통법" --format json
```

## 📚 사용법

### 기본 명령어 구조

```
sejong [command] [arguments] [flags]
```

### 주요 명령어

#### 법령 검색
```bash
# 기본 검색
sejong law "검색어"

# JSON 출력
sejong law "검색어" --format json

# 상세 로그 출력
sejong law "검색어" --verbose
```

#### 설정 관리
```bash
# API 키 설정
sejong config set law.key YOUR_API_KEY

# API 키 확인
sejong config get law.key
```

#### 도움말
```bash
# 전체 도움말
sejong --help

# 버전 확인
sejong --version
```

## 📊 개발 상태

### 구현 완료 ✅
- [x] 프로젝트 초기 구조 설정
- [x] Cobra 기반 CLI 프레임워크
- [x] 기본 명령어 구조

### 개발 중 🚧
- [ ] Viper 기반 설정 관리
- [ ] 국가법령정보센터 API 연동
- [ ] 법령 검색 기능

### 계획됨 📋
- [ ] 테이블/JSON 출력 포맷터
- [ ] 에러 처리 및 로깅 개선
- [ ] 단위 테스트 및 통합 테스트
- [ ] GitHub Actions CI/CD
- [ ] 멀티 플랫폼 바이너리 릴리스

전체 개발 계획은 [Issues](https://github.com/pyhub-kr/pyhub-sejong-cli/issues)에서 확인할 수 있습니다.

## 🤝 기여하기

Sejong-CLI는 오픈소스 프로젝트입니다. 기여를 환영합니다!

### 기여 방법

1. 이슈를 먼저 등록해주세요
2. Fork 후 feature 브랜치 생성
3. 변경사항 커밋
4. Pull Request 제출

### 개발 환경 설정

```bash
# 의존성 설치
go mod download

# 테스트 실행
go test ./...

# 빌드
go build -o sejong cmd/sejong/main.go
```

## 📄 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 🙏 감사의 말

- [국가법령정보센터](https://www.law.go.kr) - 오픈 API 제공
- [Cobra](https://github.com/spf13/cobra) - CLI 프레임워크
- [Viper](https://github.com/spf13/viper) - 설정 관리
- [tablewriter](https://github.com/olekukonko/tablewriter) - 테이블 출력

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/pyhub-kr">PyHub Korea</a>
</p>