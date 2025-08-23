# Contributing to Sejong CLI

[한국어](#한국어) | [English](#english)

---

## 한국어

Sejong CLI 프로젝트에 기여해주셔서 감사합니다! 이 문서는 프로젝트에 기여하는 방법을 안내합니다.

### 🤝 기여 방법

#### 1. 이슈 등록

기여를 시작하기 전에 먼저 이슈를 등록해주세요:

- **버그 리포트**: 발견한 버그를 상세히 설명해주세요
- **기능 제안**: 새로운 기능이나 개선사항을 제안해주세요
- **문서 개선**: 문서의 오류나 개선점을 알려주세요

이슈 템플릿:
```markdown
## 설명
[문제나 제안사항을 명확히 설명]

## 재현 방법 (버그의 경우)
1. [첫 번째 단계]
2. [두 번째 단계]
3. [오류 발생]

## 예상 동작
[어떻게 동작해야 하는지 설명]

## 환경
- OS: [예: macOS 14.0]
- Go 버전: [예: 1.21]
- Sejong CLI 버전: [예: v1.2534.1]
```

#### 2. 개발 환경 설정

```bash
# 저장소 Fork
# GitHub에서 Fork 버튼을 클릭

# Fork한 저장소 클론
git clone https://github.com/YOUR_USERNAME/pyhub-sejong-cli.git
cd pyhub-sejong-cli

# 원본 저장소를 upstream으로 추가
git remote add upstream https://github.com/pyhub-kr/pyhub-sejong-cli.git

# 의존성 설치
go mod download

# 개발 브랜치 생성
git checkout -b feature/your-feature-name
```

#### 3. 코드 작성

##### 코딩 스타일

- Go 표준 포맷팅 사용: `go fmt ./...`
- 의미 있는 변수명과 함수명 사용
- 주석은 영어로 작성 (한국어 주석도 허용)
- 에러는 명시적으로 처리

##### 프로젝트 구조

```
internal/
├── api/        # API 클라이언트
├── cmd/        # CLI 명령어
├── config/     # 설정 관리
└── output/     # 출력 포맷터
```

##### 커밋 메시지 형식

```
<type>: <subject>

<body>

<footer>
```

타입:
- `feat`: 새로운 기능
- `fix`: 버그 수정
- `docs`: 문서 변경
- `style`: 코드 포맷팅
- `refactor`: 리팩토링
- `test`: 테스트 추가/수정
- `chore`: 빌드, 도구 설정 등

예시:
```
feat: add support for pagination in law search

- Add --page and --size flags
- Implement pagination logic in API client
- Update table output to show page info

Closes #123
```

#### 4. 테스트

```bash
# 테스트 실행
make test

# 커버리지 확인
make test-coverage

# 린트 검사
make lint

# 통합 테스트
go test -tags=integration ./...
```

테스트 작성 가이드:
- 모든 새 기능에 대한 단위 테스트 작성
- 테스트 커버리지 80% 이상 유지
- 테이블 드리븐 테스트 사용 권장

```go
func TestSearchLaw(t *testing.T) {
    tests := []struct {
        name    string
        query   string
        want    int
        wantErr bool
    }{
        {"valid query", "민법", 10, false},
        {"empty query", "", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 테스트 로직
        })
    }
}
```

#### 5. Pull Request 제출

```bash
# 변경사항 커밋
git add .
git commit -m "feat: your feature description"

# Fork한 저장소에 푸시
git push origin feature/your-feature-name
```

PR 체크리스트:
- [ ] 이슈 번호 참조 (예: `Closes #123`)
- [ ] 테스트 통과
- [ ] 문서 업데이트 (필요시)
- [ ] 커밋 메시지 컨벤션 준수
- [ ] 코드 리뷰 요청

### 📋 코드 리뷰 프로세스

1. **자동 검사**: CI/CD 파이프라인이 자동으로 테스트와 린트를 실행합니다
2. **코드 리뷰**: 메인테이너가 코드를 리뷰하고 피드백을 제공합니다
3. **수정**: 피드백에 따라 코드를 수정합니다
4. **승인**: 리뷰가 완료되면 PR이 머지됩니다

### 🏗️ 릴리스 프로세스

이 프로젝트는 HeadVer 버저닝을 사용합니다:
- 형식: `{head}.{yearweek}.{build}`
- 예: `1.2534.1` (head=1, 2025년 34주, 빌드 1)

릴리스 생성:
```bash
# 헤드 버전 확인
./scripts/headver.sh

# 태그 생성
git tag v1.2534.1

# 릴리스 (자동화됨)
git push origin v1.2534.1
```

### 🐛 버그 리포트

버그를 발견하면:
1. 기존 이슈를 먼저 확인해주세요
2. 재현 가능한 최소한의 예제를 제공해주세요
3. 환경 정보를 포함해주세요
4. 가능하면 해결 방법도 제안해주세요

### 💬 커뮤니티

- **Issues**: [GitHub Issues](https://github.com/pyhub-kr/pyhub-sejong-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/pyhub-kr/pyhub-sejong-cli/discussions)
- **PyHub Korea**: [https://github.com/pyhub-kr](https://github.com/pyhub-kr)

---

## English

Thank you for contributing to Sejong CLI! This document guides you through the contribution process.

### 🤝 How to Contribute

#### 1. Create an Issue

Before starting your contribution, please create an issue:

- **Bug Reports**: Describe the bug in detail
- **Feature Requests**: Propose new features or improvements
- **Documentation**: Report errors or suggest improvements

Issue Template:
```markdown
## Description
[Clear description of the problem or suggestion]

## Steps to Reproduce (for bugs)
1. [First step]
2. [Second step]
3. [Error occurs]

## Expected Behavior
[Describe how it should work]

## Environment
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.21]
- Sejong CLI version: [e.g., v1.2534.1]
```

#### 2. Development Setup

```bash
# Fork the repository
# Click the Fork button on GitHub

# Clone your fork
git clone https://github.com/YOUR_USERNAME/pyhub-sejong-cli.git
cd pyhub-sejong-cli

# Add upstream remote
git remote add upstream https://github.com/pyhub-kr/pyhub-sejong-cli.git

# Install dependencies
go mod download

# Create a feature branch
git checkout -b feature/your-feature-name
```

#### 3. Writing Code

##### Coding Style

- Use Go standard formatting: `go fmt ./...`
- Use meaningful variable and function names
- Write comments in English (Korean comments are also accepted)
- Handle errors explicitly

##### Project Structure

```
internal/
├── api/        # API client
├── cmd/        # CLI commands
├── config/     # Configuration management
└── output/     # Output formatters
```

##### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting
- `refactor`: Refactoring
- `test`: Adding/modifying tests
- `chore`: Build, tool configuration, etc.

Example:
```
feat: add support for pagination in law search

- Add --page and --size flags
- Implement pagination logic in API client
- Update table output to show page info

Closes #123
```

#### 4. Testing

```bash
# Run tests
make test

# Check coverage
make test-coverage

# Run linter
make lint

# Integration tests
go test -tags=integration ./...
```

Testing Guidelines:
- Write unit tests for all new features
- Maintain test coverage above 80%
- Use table-driven tests when appropriate

```go
func TestSearchLaw(t *testing.T) {
    tests := []struct {
        name    string
        query   string
        want    int
        wantErr bool
    }{
        {"valid query", "civil law", 10, false},
        {"empty query", "", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

#### 5. Submit Pull Request

```bash
# Commit your changes
git add .
git commit -m "feat: your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

PR Checklist:
- [ ] Reference issue number (e.g., `Closes #123`)
- [ ] Tests pass
- [ ] Documentation updated (if needed)
- [ ] Commit message convention followed
- [ ] Code review requested

### 📋 Code Review Process

1. **Automated Checks**: CI/CD pipeline automatically runs tests and linting
2. **Code Review**: Maintainers review code and provide feedback
3. **Updates**: Make changes based on feedback
4. **Approval**: PR is merged after review completion

### 🏗️ Release Process

This project uses HeadVer versioning:
- Format: `{head}.{yearweek}.{build}`
- Example: `1.2534.1` (head=1, year 2025 week 34, build 1)

Creating a Release:
```bash
# Check head version
./scripts/headver.sh

# Create tag
git tag v1.2534.1

# Release (automated)
git push origin v1.2534.1
```

### 🐛 Bug Reports

When you find a bug:
1. Check existing issues first
2. Provide a minimal reproducible example
3. Include environment information
4. Suggest a fix if possible

### 💬 Community

- **Issues**: [GitHub Issues](https://github.com/pyhub-kr/pyhub-sejong-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/pyhub-kr/pyhub-sejong-cli/discussions)
- **PyHub Korea**: [https://github.com/pyhub-kr](https://github.com/pyhub-kr)

---

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

Examples of behavior that contributes to creating a positive environment include:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project team. All complaints will be reviewed and investigated and will result in a response that is deemed necessary and appropriate to the circumstances.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/pyhub-kr">PyHub Korea</a>
</p>