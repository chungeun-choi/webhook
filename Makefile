# This is a simple Makefile for building and running the application
BINARY_NAME=myapp

# This is the default target, which will be executed by
BUILD_FLAGS=-ldflags "-s -w"

# Build the application
all: build

# Build the application
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(PKG)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64 $(PKG)

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(PKG)

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64 $(PKG)


# Build Image for platform
docker-build-darwin-arm64:
	GOARCH=arm64 GOOS=darwin go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64 $(PKG)
	docker build --build-arg BINARY=$(BINARY_NAME)-darwin-arm64 -t $(BINARY_NAME)-darwin-arm64 -f Dockerfile . --no-cache

docker-build-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(PKG)
	docker build --build-arg BINARY=$(BINARY_NAME)-darwin-amd64 -t $(BINARY_NAME)-darwin-amd64 -f Dockerfile . --no-cache

docker-build-linux-arm64:
	GOARCH=arm64 GOOS=linux go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64 $(PKG)
	docker build --build-arg BINARY=$(BINARY_NAME)-linux-arm64 -t $(BINARY_NAME)-linux-arm64 -f Dockerfile . --no-cache

docker-build-linux-amd64:
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(PKG)
	docker build --build-arg BINARY=$(BINARY_NAME)-linux-amd64 -t $(BINARY_NAME)-linux-amd64 -f Dockerfile . --no-cache


# Build Image for all platforms
docker-build: docker-build-linux-amd64 docker-build-linux-arm64 docker-build-darwin-amd64 docker-build-darwin-arm64

# Run the application
run:
	./$(BINARY_NAME)

# Dependency management
deps:
	go mod tidy
	go mod download

# Clean the build
clean:
	rm -f $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)-linux-arm64 $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-darwin-arm64
	go clean

.PHONY: all build build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 docker-build-linux-amd64 docker-build-linux-arm64 docker-build-darwin-amd64 docker-build-darwin-arm64 docker-build run fmt deps clean
