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
- 📋 **다양한 출력 형식**: 테이블 형식 또는 JSON 형식 지원
- ⚡ **간편한 설정**: 한 번의 API 키 설정으로 계속 사용
- 📄 **페이지네이션**: 대량의 검색 결과를 페이지별로 조회
- 🎯 **스마트 온보딩**: 처음 사용자를 위한 친절한 안내
- 🔄 **자동 재시도**: 네트워크 오류 시 자동 재시도
- 🌈 **컬러 출력**: 가독성 높은 컬러 터미널 출력

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
👉 [https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9](https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9)

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

# 페이지 지정
sejong law "민법" --page 2 --size 20
```

### 📚 명령어 가이드

#### 법령 검색

```bash
# 기본 검색
sejong law "검색어"

# 출력 형식 지정
sejong law "검색어" --format json  # JSON 형식
sejong law "검색어" --format table # 테이블 형식 (기본값)

# 페이지네이션
sejong law "검색어" --page 2 --size 20

# 상세 로그 출력
sejong law "검색어" --verbose
sejong law "검색어" -v  # 단축 옵션
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
sejong config --help
```

### 📊 출력 예제

#### 테이블 형식 (기본)

```text
총 3개의 법령을 찾았습니다.

번호  법령명                                          법령구분   소관부처        시행일자
----------------------------------------------------------------------------------------------------
1     개인정보 보호법                                  법률      개인정보보호위원회  2024-03-15
2     개인정보 보호법 시행령                            대통령령   개인정보보호위원회  2024-03-15
3     개인정보 보호법 시행규칙                          부령      개인정보보호위원회  2024-03-15
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
      "시행일자": "20240315"
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
- 📋 **Multiple Output Formats**: Support for table and JSON formats
- ⚡ **Simple Configuration**: One-time API key setup for continuous use
- 📄 **Pagination**: Browse large search results page by page
- 🎯 **Smart Onboarding**: Friendly guidance for first-time users
- 🔄 **Auto Retry**: Automatic retry on network errors
- 🌈 **Color Output**: Readable colored terminal output

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
👉 [https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9](https://www.law.go.kr/LSW/opn/prvsn/opnPrvsnInfoP.do?mode=9)

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
sejong law "search term" --format json  # JSON format
sejong law "search term" --format table # Table format (default)

# Pagination
sejong law "search term" --page 2 --size 20

# Verbose logging
sejong law "search term" --verbose
sejong law "search term" -v  # Short option
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
sejong config --help
```

### 📊 Output Examples

#### Table Format (Default)

```text
Found 3 laws in total.

No.   Law Name                                        Type      Department              Effective Date
----------------------------------------------------------------------------------------------------
1     Personal Information Protection Act             Law       Personal Information    2024-03-15
                                                                 Protection Commission
2     Personal Information Protection Act             Decree    Personal Information    2024-03-15
      Enforcement Decree                                        Protection Commission
3     Personal Information Protection Act             Rule      Personal Information    2024-03-15
      Enforcement Rules                                         Protection Commission
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