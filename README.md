# Sejong CLI

> 🏛️ 터미널에서 빠르게 검색하는 대한민국 법령 정보

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://github.com/pyhub-kr/pyhub-sejong-cli/actions/workflows/test.yml/badge.svg)](https://github.com/pyhub-kr/pyhub-sejong-cli/actions/workflows/test.yml)
[![Build](https://github.com/pyhub-kr/pyhub-sejong-cli/actions/workflows/build.yml/badge.svg)](https://github.com/pyhub-kr/pyhub-sejong-cli/actions/workflows/build.yml)

## 📑 목차 / Table of Contents

### 한국어
- [소개](#-소개)
- [주요 기능](#-주요-기능)
- [설치](#-설치)
- [빠른 시작](#-빠른-시작)
- [명령어 가이드](#-명령어-가이드)
- [출력 예제](#-출력-예제)
- [개발](#️-개발)
- [문제 해결](#-문제-해결)
- [기여하기](#-기여하기)
- [라이선스](#-라이선스)

### English
- [Introduction](#-introduction)
- [Key Features](#-key-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [License](#-license)
- [Acknowledgments](#-acknowledgments)

---

## 한국어

### 📖 소개

**Sejong CLI**는 국가법령정보센터 오픈 API를 활용하여 터미널에서 한국 법령 정보를 빠르고 쉽게 검색할 수 있는 명령줄 도구입니다.

### ✨ 주요 기능

- 🔍 **빠른 법령 검색**: 터미널에서 즉시 법령 정보 조회
- 📖 **법령 상세 조회**: 법령ID로 조문, 별표, 부칙 등 상세 정보 확인
- 📜 **법령 이력 조회**: 법령의 제/개정 이력 및 시행 이력 추적
- ⚖️ **판례 검색**: 대법원 및 각급 법원의 판례 검색 및 상세 조회
- 📜 **행정규칙 검색**: 정부 부처의 고시, 훈령, 예규 등 검색
- 📚 **법령해석례 검색**: 법령 적용과 해석에 대한 정부 공식 견해 조회
- 🏛️ **자치법규 검색**: 지방자치단체의 조례 및 규칙 검색
- 📋 **다양한 출력 형식**: Table, JSON, Markdown, CSV, HTML, HTML-Simple 지원
- ⚡ **간편한 설정**: 한 번의 API 키 설정으로 계속 사용
- 📄 **페이지네이션**: 대량의 검색 결과를 페이지별로 조회
- 🎯 **스마트 온보딩**: 처음 사용자를 위한 친절한 안내
- 🔄 **자동 재시도**: 네트워크 오류 시 자동 재시도
- 🌈 **향상된 테이블 출력**: Box-drawing 문자와 컬러로 가독성 극대화

### 🚀 설치

#### 바이너리 다운로드 (권장)

최신 릴리스는 [Releases](https://github.com/pyhub-kr/pyhub-sejong-cli/releases) 페이지에서 다운로드할 수 있습니다.

##### macOS (Apple Silicon)
```bash
# 최신 버전 다운로드
curl -LO https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Darwin_arm64.tar.gz
tar -xzf pyhub-sejong-cli_Darwin_arm64.tar.gz
sudo mv sejong /usr/local/bin/
```

##### macOS (Intel)
```bash
curl -LO https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Darwin_x86_64.tar.gz
tar -xzf pyhub-sejong-cli_Darwin_x86_64.tar.gz
sudo mv sejong /usr/local/bin/
```

##### Windows
```powershell
# PowerShell에서 실행
Invoke-WebRequest -Uri https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Windows_x86_64.zip -OutFile sejong.zip
Expand-Archive -Path sejong.zip -DestinationPath .
# sejong.exe를 PATH에 추가하거나 원하는 위치로 이동
```

##### Linux
```bash
curl -LO https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Linux_x86_64.tar.gz
tar -xzf pyhub-sejong-cli_Linux_x86_64.tar.gz
sudo mv sejong /usr/local/bin/
```

#### 소스에서 빌드

Go 1.21 이상이 필요합니다.

```bash
# 저장소 클론
git clone https://github.com/pyhub-kr/pyhub-sejong-cli.git
cd pyhub-sejong-cli

# 빌드
make build

# 또는 직접 빌드
go build -o sejong ./cmd/sejong

# 설치
make install
```

### 🎯 빠른 시작

#### 1. API 키 발급

국가법령정보센터에서 오픈 API 인증키를 발급받으세요:
👉 [https://open.law.go.kr/LSO/openApi/cuAskList.do](https://open.law.go.kr/LSO/openApi/cuAskList.do)

#### 2. API 키 설정

```bash
sejong config set law.key YOUR_API_KEY
```

#### 3. 첫 검색

```bash
# 법령 검색
sejong law "개인정보 보호법"

# JSON 형식으로 출력
sejong law "도로교통법" --format json

# Markdown 형식으로 출력 (문서 작성용)
sejong law "민법" --format markdown > laws.md

# CSV 형식으로 저장 (Excel에서 열기)
sejong law "상법" --format csv > laws.csv

# 페이지 지정 (기본 50건씩)
sejong law "민법" --page 2

# 법령 상세 조회
sejong law detail 011357

# 법령 이력 조회
sejong law history 011357
```

### 📚 명령어 가이드

#### 법령 검색

```bash
# 기본 검색
sejong law "검색어"

# 출력 형식 지정
sejong law "검색어" --format json       # JSON 형식
sejong law "검색어" --format table      # 테이블 형식 (기본값)
sejong law "검색어" --format markdown   # Markdown 형식
sejong law "검색어" --format csv        # CSV 형식 (Excel 호환)
sejong law "검색어" --format html       # HTML 형식
sejong law "검색어" --format html-simple # HTML 형식 (CSS 없음, LLM AI용)

# 페이지네이션
sejong law "검색어" --page 2 --size 50

# 검색 소스 지정
sejong law "검색어" --source all   # 통합 검색 (국가법령 + 자치법규)
sejong law "검색어" --source nlic  # 국가법령만
sejong law "검색어" --source elis  # 자치법규만

# 상세 로그 출력
sejong law "검색어" --verbose
sejong law "검색어" -v  # 단축 옵션
```

#### 법령 상세 조회

```bash
# 기본 상세 조회
sejong law detail 법령ID

# 조문 포함
sejong law detail 법령ID --articles

# 별표 포함
sejong law detail 법령ID --tables

# 부칙 포함
sejong law detail 법령ID --addendum

# 모두 포함
sejong law detail 법령ID --articles --tables --addendum

# JSON 형식으로 출력
sejong law detail 법령ID --format json
```

#### 법령 이력 조회

```bash
# 기본 이력 조회
sejong law history 법령ID

# 최근 N개만 조회
sejong law history 법령ID --limit 10

# JSON 형식으로 출력
sejong law history 법령ID --format json
```

#### 판례 검색

```bash
# 기본 검색
sejong precedent search "계약 해지"

# 또는 단축 명령어 사용
sejong prec search "손해배상"

# JSON 형식으로 출력
sejong precedent search "부당이득" --format json

# 페이지네이션
sejong precedent search "계약" --page 2 --size 20

# 판례 상세 조회
sejong precedent detail 12345
```

#### 행정규칙 검색

```bash
# 기본 검색
sejong admrule search "공공기관"

# 단축 명령어 사용
sejong admr search "개인정보"
sejong rule search "행정처분"

# JSON 형식으로 출력
sejong admrule search "고시" --format json

# 페이지네이션
sejong admrule search "훈령" --page 2 --size 20

# 행정규칙 상세 조회
sejong admrule detail 12345
```

#### 법령해석례 검색

```bash
# 기본 검색
sejong interpretation search "근로시간"

# 단축 명령어 사용
sejong interp search "휴가"
sejong expc search "임금"

# JSON 형식으로 출력
sejong interpretation search "퇴직금" --format json

# 페이지네이션
sejong interpretation search "근로계약" --page 2 --size 20

# 법령해석례 상세 조회
sejong interpretation detail 12345
```

#### 자치법규 (조례/규칙) 검색

```bash
# 기본 검색
sejong ordinance search "주차 조례"

# 단축 명령어 사용
sejong ord search "건축 조례"

# JSON 형식으로 출력
sejong ordinance search "환경" --format json

# 페이지네이션
sejong ordinance search "교통" --page 2 --size 50

# 자치법규 상세 조회
sejong ordinance detail ORD123456
```

#### 설정 관리

```bash
# API 키 설정
sejong config set law.key YOUR_API_KEY

# API 키 확인 (마스킹된 출력)
sejong config get law.key

# 설정 파일 경로 확인
sejong config path
```

#### 버전 및 도움말

```bash
# 버전 정보
sejong version

# 전체 도움말
sejong --help
sejong -h

# 명령별 도움말
sejong law --help
sejong precedent --help
sejong admrule --help
sejong interpretation --help
sejong config --help
```

### 📊 출력 예제

#### 테이블 형식 (기본) - 향상된 버전

```text
총 3개의 법령을 찾았습니다.

│──────│────────│──────────────────────────────│──────────│────────────────────│────────────│
│ 번호 │ 법령ID │ 법령명                         │ 법령구분 │ 소관부처           │ 시행일자   │
│──────│────────│──────────────────────────────│──────────│────────────────────│────────────│
│ 1    │ 011357 │ 개인정보 보호법                │ 법률     │ 개인정보보호위원회 │ 2025-03-13 │
│ 2    │ 011468 │ 개인정보 보호법 시행령        │ 대통령령 │ 개인정보보호위원회 │ 2025-07-01 │
│──────│────────│──────────────────────────────│──────────│────────────────────│────────────│
```

#### Markdown 형식

```markdown
## 검색 결과

총 **3**개의 법령을 찾았습니다.

| 번호 | 법령ID | 법령명 | 법령구분 | 소관부처 | 시행일자 |
| --- | --- | --- | --- | --- | --- |
| 1 | 011357 | 개인정보 보호법 | 법률 | 개인정보보호위원회 | 2025-03-13 |
| 2 | 011468 | 개인정보 보호법 시행령 | 대통령령 | 개인정보보호위원회 | 2025-07-01 |
```

#### CSV 형식 (Excel 호환)

```csv
번호,법령ID,법령명,법령구분,소관부처,시행일자
1,011357,개인정보 보호법,법률,개인정보보호위원회,2025-03-13
2,011468,개인정보 보호법 시행령,대통령령,개인정보보호위원회,2025-07-01
```

#### HTML Simple 형식 (LLM AI용)

```html
<h2>검색 결과</h2>
<p>총 <strong>2</strong>개의 법령을 찾았습니다.</p>
<table>
  <thead>
    <tr>
      <th>번호</th>
      <th>법령ID</th>
      <th>법령명</th>
      <th>법령구분</th>
      <th>소관부처</th>
      <th>시행일자</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>1</td>
      <td>011357</td>
      <td>개인정보 보호법</td>
      <td>법률</td>
      <td>개인정보보호위원회</td>
      <td>2025-03-13</td>
    </tr>
  </tbody>
</table>
```

#### JSON 형식

```json
{
  "totalCnt": 3,
  "page": 1,
  "law": [
    {
      "법령ID": "173995",
      "법령명한글": "개인정보 보호법",
      "법령구분명": "법률",
      "소관부처명": "개인정보보호위원회",
      "시혹일자": "20240315"
    }
  ]
}
```

### 🛠️ 개발

#### 개발 환경 설정

```bash
# 의존성 설치
go mod download

# 테스트 실행
make test

# 테스트 커버리지
make test-coverage

# 코드 포맷팅
make fmt

# 린트 검사
make lint
```

#### 빌드

```bash
# 현재 플랫폼용 빌드
make build

# 개발 빌드 (race detector 포함)
make dev

# 모든 플랫폼용 빌드 (릴리스 스냅샷)
make release-snapshot
```

### 🐛 문제 해결

#### API 키가 설정되지 않음

```bash
# API 키가 올바르게 설정되었는지 확인
sejong config get law.key

# API 키 재설정
sejong config set law.key YOUR_NEW_API_KEY
```

#### 네트워크 오류

- 인터넷 연결 상태를 확인하세요
- 방화벽이나 프록시 설정을 확인하세요
- API 서버 상태를 확인하세요: [https://www.law.go.kr](https://www.law.go.kr)

#### 권한 오류 (macOS/Linux)

```bash
# 실행 권한 부여
chmod +x sejong

# sudo를 사용하여 시스템 경로에 설치
sudo mv sejong /usr/local/bin/
```

### 🤝 기여하기

기여를 환영합니다! [CONTRIBUTING.md](CONTRIBUTING.md)를 참조하세요.

1. 이슈를 먼저 등록해주세요
2. Fork 후 feature 브랜치 생성 (`git checkout -b feature/AmazingFeature`)
3. 변경사항 커밋 (`git commit -m 'Add some AmazingFeature'`)
4. 브랜치에 Push (`git push origin feature/AmazingFeature`)
5. Pull Request 제출

### 📄 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

---

## English

### 📖 Introduction

**Sejong CLI** is a command-line tool that allows you to quickly and easily search Korean law information from the terminal using the National Law Information Center Open API.

### ✨ Key Features

- 🔍 **Fast Law Search**: Instantly search law information from terminal
- 📖 **Law Details**: View detailed information including articles, tables, and addenda
- 📜 **Law History**: Track enactment and amendment history of laws
- ⚖️ **Precedent Search**: Search and view court precedents from all levels
- 📜 **Administrative Rule Search**: Search government notices, directives, and regulations
- 📚 **Legal Interpretation Search**: View official government interpretations of laws
- 🏛️ **Local Ordinance Search**: Search local government ordinances and rules
- 📋 **Multiple Output Formats**: Table, JSON, Markdown, CSV, HTML, HTML-Simple formats
- ⚡ **Simple Configuration**: One-time API key setup for continuous use
- 📄 **Pagination**: Browse large search results page by page
- 🎯 **Smart Onboarding**: Friendly guidance for first-time users
- 🔄 **Auto Retry**: Automatic retry on network errors
- 🌈 **Enhanced Table Output**: Box-drawing characters with color for better readability

### 🚀 Installation

#### Download Binary (Recommended)

Download the latest release from the [Releases](https://github.com/pyhub-kr/pyhub-sejong-cli/releases) page.

##### macOS (Apple Silicon)
```bash
curl -LO https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Darwin_arm64.tar.gz
tar -xzf pyhub-sejong-cli_Darwin_arm64.tar.gz
sudo mv sejong /usr/local/bin/
```

##### macOS (Intel)
```bash
curl -LO https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Darwin_x86_64.tar.gz
tar -xzf pyhub-sejong-cli_Darwin_x86_64.tar.gz
sudo mv sejong /usr/local/bin/
```

##### Windows
```powershell
# Run in PowerShell
Invoke-WebRequest -Uri https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Windows_x86_64.zip -OutFile sejong.zip
Expand-Archive -Path sejong.zip -DestinationPath .
# Add sejong.exe to PATH or move to desired location
```

##### Linux
```bash
curl -LO https://github.com/pyhub-kr/pyhub-sejong-cli/releases/latest/download/pyhub-sejong-cli_Linux_x86_64.tar.gz
tar -xzf pyhub-sejong-cli_Linux_x86_64.tar.gz
sudo mv sejong /usr/local/bin/
```

#### Build from Source

Requires Go 1.21 or higher.

```bash
git clone https://github.com/pyhub-kr/pyhub-sejong-cli.git
cd pyhub-sejong-cli
make build
make install
```

### 🎯 Quick Start

#### 1. Get API Key

Get your Open API authentication key from the National Law Information Center:
👉 [https://open.law.go.kr/LSO/openApi/cuAskList.do](https://open.law.go.kr/LSO/openApi/cuAskList.do)

#### 2. Configure API Key

```bash
sejong config set law.key YOUR_API_KEY
```

#### 3. First Search

```bash
# Search laws
sejong law "personal information"

# Output in JSON format
sejong law "traffic law" --format json

# Specify page
sejong law "civil law" --page 2 --size 20
```

### 📚 Command Guide

#### Law Search

```bash
# Basic search
sejong law "search term"

# Specify output format
sejong law "search term" --format json       # JSON format
sejong law "search term" --format table      # Table format (default)
sejong law "search term" --format markdown   # Markdown format
sejong law "search term" --format csv        # CSV format (Excel compatible)
sejong law "search term" --format html       # HTML format
sejong law "search term" --format html-simple # HTML format without CSS (for LLM AI)

# Pagination
sejong law "search term" --page 2 --size 50

# Search source
sejong law "search term" --source all   # Unified search
sejong law "search term" --source nlic  # National laws only
sejong law "search term" --source elis  # Local ordinances only

# Verbose logging
sejong law "search term" --verbose
sejong law "search term" -v  # Short option
```

#### Law Details

```bash
# Basic detail view
sejong law detail LAW_ID

# Include articles
sejong law detail LAW_ID --articles

# Include tables
sejong law detail LAW_ID --tables

# Include addenda
sejong law detail LAW_ID --addendum

# Output in JSON format
sejong law detail LAW_ID --format json
```

#### Law History

```bash
# Basic history view
sejong law history LAW_ID

# Limit number of records
sejong law history LAW_ID --limit 10

# Output in JSON format
sejong law history LAW_ID --format json
```

#### Precedent Search

```bash
# Basic search
sejong precedent search "contract termination"

# Or use alias
sejong prec search "damages"

# Output in JSON format
sejong precedent search "unjust enrichment" --format json

# Pagination
sejong precedent search "contract" --page 2 --size 20

# View precedent details
sejong precedent detail 12345
```

#### Administrative Rule Search

```bash
# Basic search
sejong admrule search "public institution"

# Use aliases
sejong admr search "personal information"
sejong rule search "administrative action"

# Output in JSON format
sejong admrule search "notice" --format json

# Pagination
sejong admrule search "directive" --page 2 --size 20

# View administrative rule details
sejong admrule detail 12345
```

#### Legal Interpretation Search

```bash
# Basic search
sejong interpretation search "working hours"

# Use aliases
sejong interp search "vacation"
sejong expc search "wages"

# Output in JSON format
sejong interpretation search "retirement" --format json

# Pagination
sejong interpretation search "employment" --page 2 --size 20

# View legal interpretation details
sejong interpretation detail 12345
```

#### Local Ordinance Search

```bash
# Basic search
sejong ordinance search "parking ordinance"

# Use alias
sejong ord search "building ordinance"

# Output in JSON format
sejong ordinance search "environment" --format json

# Pagination
sejong ordinance search "traffic" --page 2 --size 50

# View ordinance details
sejong ordinance detail ORD123456
```

#### Configuration Management

```bash
# Set API key
sejong config set law.key YOUR_API_KEY

# Check API key (masked output)
sejong config get law.key

# Check configuration file path
sejong config path
```

#### Version and Help

```bash
# Version information
sejong version

# General help
sejong --help
sejong -h

# Command-specific help
sejong law --help
sejong precedent --help
sejong admrule --help
sejong interpretation --help
sejong config --help
```

### 📊 Output Examples

#### Enhanced Table Format (Default)

```text
Found 2 laws in total.

│──────│────────│──────────────────────────────│──────────│────────────────────│────────────│
│ No.  │ Law ID │ Law Name                     │ Type     │ Department         │ Date       │
│──────│────────│──────────────────────────────│──────────│────────────────────│────────────│
│ 1    │ 011357 │ Personal Information         │ Law      │ Privacy Commission │ 2025-03-13 │
│      │        │ Protection Act               │          │                    │            │
│──────│────────│──────────────────────────────│──────────│────────────────────│────────────│
```

#### Markdown Format

```markdown
## Search Results

Found **2** laws in total.

| No. | Law ID | Law Name | Type | Department | Date |
| --- | --- | --- | --- | --- | --- |
| 1 | 011357 | Personal Information Protection Act | Law | Privacy Commission | 2025-03-13 |
| 2 | 011468 | Personal Information Protection Act Enforcement Decree | Decree | Privacy Commission | 2025-07-01 |
```

#### CSV Format

```csv
No.,Law ID,Law Name,Type,Department,Date
1,011357,Personal Information Protection Act,Law,Privacy Commission,2025-03-13
2,011468,Personal Information Protection Act Enforcement Decree,Decree,Privacy Commission,2025-07-01
```

#### HTML Simple Format (for LLM AI)

```html
<h2>Search Results</h2>
<p>Found <strong>2</strong> laws in total.</p>
<table>
  <thead>
    <tr>
      <th>No.</th>
      <th>Law ID</th>
      <th>Law Name</th>
      <th>Type</th>
      <th>Department</th>
      <th>Date</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>1</td>
      <td>011357</td>
      <td>Personal Information Protection Act</td>
      <td>Law</td>
      <td>Privacy Commission</td>
      <td>2025-03-13</td>
    </tr>
  </tbody>
</table>
```

#### JSON Format

```json
{
  "totalCnt": 3,
  "page": 1,
  "law": [
    {
      "법령ID": "173995",
      "법령명한글": "Personal Information Protection Act",
      "법령구분명": "Law",
      "소관부처명": "Personal Information Protection Commission",
      "시행일자": "20240315"
    }
  ]
}
```

### 🛠️ Development

#### Development Environment Setup

```bash
# Install dependencies
go mod download

# Run tests
make test

# Test coverage
make test-coverage

# Code formatting
make fmt

# Lint check
make lint
```

#### Build

```bash
# Build for current platform
make build

# Development build (with race detector)
make dev

# Build for all platforms (release snapshot)
make release-snapshot
```

### 🐛 Troubleshooting

#### API Key Not Set

```bash
# Check if API key is properly set
sejong config get law.key

# Reset API key
sejong config set law.key YOUR_NEW_API_KEY
```

#### Network Errors

- Check your internet connection
- Verify firewall or proxy settings
- Check API server status: [https://www.law.go.kr](https://www.law.go.kr)

#### Permission Errors (macOS/Linux)

```bash
# Grant execution permission
chmod +x sejong

# Install to system path with sudo
sudo mv sejong /usr/local/bin/
```

### 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md).

1. Create an issue first
2. Fork and create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Submit a Pull Request

### 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### 🙏 Acknowledgments

- [National Law Information Center](https://www.law.go.kr) - Open API Provider
- [Cobra](https://github.com/spf13/cobra) - CLI Framework
- [Viper](https://github.com/spf13/viper) - Configuration Management
- [tablewriter](https://github.com/olekukonko/tablewriter) - Table Output

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/pyhub-kr">PyHub Korea</a>
</p>