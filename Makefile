.PHONY: build run clean \
        dev lint fmt \
        test test-coverage \
        create-user update-password quality \
        release release-darwin-arm64 release-darwin-amd64 \
        release-linux-arm64 release-linux-amd64 \
        release-windows-arm64 release-windows-amd64 \
        bump-major bump-minor bump-patch \
        docker-build test-e2e test-e2e-ui

# Variables
BINARY_NAME=dist/bin/captain
MAIN_FILE=main.go
VERSION := $(shell grep 'Version = "' system/system.go | cut -d'"' -f2)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build commands
frontend-build:
	cd frontend && $(MAKE) build

build: frontend-build
	mkdir -p dist/bin
	go build -o $(BINARY_NAME) $(MAIN_FILE)

clean:
	go clean
	rm -rf dist

# Run commands
run: build
	./$(BINARY_NAME) run

dev:
	CAPTAIN_DEBUG=true air

# Test commands
test:
	go test ./... -v

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-data-cleanup:
	rm -f testdata/*

# E2E test commands
test-e2e: test-data-cleanup build
	npx playwright test

test-e2e-ui: test-data-cleanup build
	npx playwright test --ui

# Code quality commands
frontend-format:
	cd frontend && $(MAKE) format

frontend-format-check:
	cd frontend && $(MAKE) format-check

lint:
	golangci-lint run

format: frontend-format
	go fmt ./...

check-format: frontend-format-check
	go fmt ./...

quality: check-format lint test test-e2e

# Release commands for each platform
release-darwin-arm64:
	mkdir -p dist/bin
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -o "dist/bin/captain-darwin-arm64" .
	mkdir -p dist/zip
	cd dist/bin && zip -r "../zip/captain-darwin-arm64.zip" "captain-darwin-arm64"

release-darwin-amd64:
	mkdir -p dist/bin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -o "dist/bin/captain-darwin-amd64" .
	mkdir -p dist/zip
	cd dist/bin && zip -r "../zip/captain-darwin-amd64.zip" "captain-darwin-amd64"

release-linux-arm64:
	mkdir -p dist/bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -o "dist/bin/captain-linux-arm64" .
	mkdir -p dist/zip
	cd dist/bin && zip -r "../zip/captain-linux-arm64.zip" "captain-linux-arm64"

release-linux-amd64:
	mkdir -p dist/bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -o "dist/bin/captain-linux-amd64" .
	mkdir -p dist/zip
	cd dist/bin && zip -r "../zip/captain-linux-amd64.zip" "captain-linux-amd64"

release-windows-arm64:
	mkdir -p dist/bin
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -o "dist/bin/captain-windows-arm64.exe" .
	mkdir -p dist/zip
	cd dist/bin && zip -r "../zip/captain-windows-arm64.zip" "captain-windows-arm64.exe"

release-windows-amd64:
	mkdir -p dist/bin
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -o "dist/bin/captain-windows-amd64.exe" .
	mkdir -p dist/zip
	cd dist/bin && zip -r "../zip/captain-windows-amd64.zip" "captain-windows-amd64.exe"

release-darwin: release-darwin-arm64 release-darwin-amd64

release-linux: release-linux-amd64 release-linux-arm64

release-windows: release-windows-amd64 release-windows-arm64

# Main release task that builds all platforms
release: test \
	release-darwin \
	release-linux \
	release-windows
	@echo "Release build complete. Binaries in dist/bin/ and archives in dist/zip/"

# Docker commands
docker-build:
	docker build -t captain:$(VERSION) -t captain:latest .

# User management commands
create-user: build
	$(BINARY_NAME) user create

update-password: build
	$(BINARY_NAME) user update-password

# Version management
bump-major:
	@echo "Current version: $(VERSION)"
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{print $$1+1".0.0"}') && \
	sed -i '' "s/Version = \".*\"/Version = \"$$NEW_VERSION\"/" system/system.go && \
	git add system/system.go && \
	git commit -m "Bump major version to $$NEW_VERSION" && \
	git tag -a "v$$NEW_VERSION" -m "Version $$NEW_VERSION"

bump-minor:
	@echo "Current version: $(VERSION)"
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{print $$1"."$$2+1".0"}') && \
	sed -i '' "s/Version = \".*\"/Version = \"$$NEW_VERSION\"/" system/system.go && \
	git add system/system.go && \
	git commit -m "Bump minor version to $$NEW_VERSION" && \
	git tag -a "v$$NEW_VERSION" -m "Version $$NEW_VERSION"

bump-patch:
	@echo "Current version: $(VERSION)"
	@NEW_VERSION=$$(echo "$(VERSION)" | awk -F. '{print $$1"."$$2"."$$3+1}') && \
	sed -i '' "s/Version = \".*\"/Version = \"$$NEW_VERSION\"/" system/system.go && \
	git add system/system.go && \
	git commit -m "Bump patch version to $$NEW_VERSION" && \
	git tag -a "v$$NEW_VERSION" -m "Version $$NEW_VERSION"

.DEFAULT_GOAL := build
