.EXPORT_ALL_VARIABLES:
.ONESHELL:
.SILENT:

MAKEFLAGS += --no-builtin-rules --no-builtin-variables

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -ldflags "-X github.com/kborovik/gcp-iam/version.Version=$(VERSION) \
                    -X github.com/kborovik/gcp-iam/version.GitCommit=$(GIT_COMMIT) \
                    -X github.com/kborovik/gcp-iam/version.BuildDate=$(BUILD_DATE)"

default: help

dist:
	mkdir -p $(@)

build: dist ## Build binaries
	$(call header,Build binaries for all platforms)
	for platform in "linux/amd64" "darwin/arm64"; do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name="gcp-iam-$(VERSION)-$$GOOS-$$GOARCH"; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build \
			-ldflags "-w -s \
				-X github.com/kborovik/gcp-iam/version.Version=$(VERSION) \
				-X github.com/kborovik/gcp-iam/version.GitCommit=$(GIT_COMMIT) \
				-X github.com/kborovik/gcp-iam/version.BuildDate=$(BUILD_DATE)" \
			-o "dist/$$output_name" .; \
	done
	cd dist
	sha256sum gcp-iam-$(VERSION)-* > checksums.txt
	gpg --default-key E4AFCA7FBB19FC029D519A524AEBB5178D5E96C1 --detach-sign --armor checksums.txt

install: ## Install binary
	$(call header,Install binary)
	go install $(LDFLAGS)
	cp gcp-iam.fish ~/.config/fish/completions/

.PHONY: version
version: ## Show current version
	echo $(VERSION)

release: build ## Create GitHub release
	$(call header,Create GitHub release $(VERSION))
	gh release create $(VERSION) \
		--generate-notes \
		--title "Release $(VERSION)" \
		--attach dist/gcp-iam-$(VERSION)-linux-amd64 \
		--attach dist/gcp-iam-$(VERSION)-darwin-arm64 \
		--attach dist/checksums.txt \
		--attach dist/checksums.txt.asc

###############################################################################
# GO Tests
###############################################################################

test: go-tidy test-config test-db test-update test-cli ## Test all modules

go-tidy:
	go fmt && go vet && go mod tidy

test-config:
	$(call header,Test Module Config)
	cd config && go test -v

test-db:
	$(call header,Test Module Database)
	cd db && go test -v

test-update:
	$(call header,Test Module Update)
	cd update && go test -v

test-cli:
	$(call header,Test CLI Command)
	go test -v

###############################################################################
# Colors and Headers
###############################################################################

TERM := xterm-256color

black := $$(tput setaf 0)
red := $$(tput setaf 1)
green := $$(tput setaf 2)
yellow := $$(tput setaf 3)
blue := $$(tput setaf 4)
magenta := $$(tput setaf 5)
cyan := $$(tput setaf 6)
white := $$(tput setaf 7)
reset := $$(tput sgr0)

define header
echo "$(blue)==> $(1) <==$(reset)"
endef

define var
echo "$(magenta)$(1)$(white): $(yellow)$(2)$(reset)"
endef

help:
	echo "$(blue)Usage: $(green)make [recipe]$(reset)"
	echo "$(blue)Recipes:$(reset)"
	awk 'BEGIN {FS = ":.*?## "; sort_cmd = "sort"} /^[a-zA-Z0-9_-]+:.*?## / \
	{ printf "  \033[33m%-15s\033[0m %s\n", $$1, $$2 | sort_cmd; } \
	END {close(sort_cmd)}' $(MAKEFILE_LIST)

prompt:
	printf "$(magenta)Continue $(white)? $(cyan)(yes/no)$(reset)"
	read -p ": " answer && [ "$$answer" = "yes" ] || exit 127
