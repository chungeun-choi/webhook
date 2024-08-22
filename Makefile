# 프로젝트의 바이너리 이름
BINARY_NAME=myapp

# Go 패키지 디렉터리
PKG=./...

# Go 빌드 플래그
BUILD_FLAGS=-ldflags "-s -w"

# 기본 타겟: 프로젝트 빌드 (기본적으로 호스트 아키텍처에 맞게 빌드)
all: build

# Linux AMD64용 바이너리 빌드
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(PKG)

# Linux ARM64용 바이너리 빌드
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64 $(PKG)

# Darwin (macOS) AMD64용 바이너리 빌드
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(PKG)

# Darwin (macOS) ARM64용 바이너리 빌드 (Apple Silicon)
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64 $(PKG)

# 모든 아키텍처 및 OS에 대해 빌드
build: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

# 기본 아키텍처에서 프로젝트 실행
run:
	./$(BINARY_NAME)

# 코드 포맷팅
fmt:
	go fmt $(PKG)

# 의존성 모듈 다운로드
deps:
	go mod tidy
	go mod download

# 바이너리 및 임시 파일 삭제
clean:
	rm -f $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)-linux-arm64 $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-darwin-arm64
	go clean

# 전체 빌드, 포맷, 의존성 관리 및 클린
.PHONY: all build build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 run fmt deps clean
