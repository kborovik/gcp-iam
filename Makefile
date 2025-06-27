.EXPORT_ALL_VARIABLES:
.ONESHELL:
.SILENT:
.PHONY: help build test clean fmt vet lint deps check run install dev

MAKEFLAGS += --no-builtin-rules --no-builtin-variables

# Get package name from go.mod
PACKAGE_NAME := $(shell go list -m | awk -F'/' '{print $$NF}')

help: deps
	echo "$(PACKAGE_NAME) - GCP IAM management tool"
	echo ""
	echo "Available targets:"
	echo "  check    - Run all checks (fmt, vet, lint, test)"
	echo "  install  - Install binary to GOPATH/bin"
	echo "  build    - Build binary in dist/"
	echo "  release  - Create GitHub release"

# Set dependencies
deps:
	mkdir -p dist/
	command -v staticcheck >/dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest
	command -v modernize >/dev/null 2>&1 || go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest
	go mod download
	go mod tidy

# Run all checks
check: fmt vet lint test
	echo "All checks passed!"

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run staticcheck linter
lint:
	staticcheck ./...
	modernize -category=efaceany -test ./...

# Run modernize linter
modernize:
	modernize -category=efaceany -fix -test ./...

# Run tests
test: deps
	go test -v ./...

LDFLAGS := -ldflags="-w -s"

# Install binary to GOPATH/bin
install:
	go install $(LDFLAGS) .
	cp $(PACKAGE_NAME).fish ~/.config/fish/completions/

# Build (optimized)
build: clean deps check
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(PACKAGE_NAME)-linux-amd64 .
	cd dist/
	sha256sum $(PACKAGE_NAME)-* > checksums.txt
	gpg --default-key E4AFCA7FBB19FC029D519A524AEBB5178D5E96C1 --detach-sign --armor checksums.txt
	echo "Release binaries built successfully!"

release: build
	$(eval VERSION := $(shell go run . --version | cut -f3 -d' '))
	gh release create $(VERSION) --generate-notes --title "Release $(VERSION)" dist/*

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f $(PACKAGE_NAME)
	go clean
