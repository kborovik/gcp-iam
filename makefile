.EXPORT_ALL_VARIABLES:
.ONESHELL:
.SILENT:

MAKEFLAGS += --no-builtin-rules --no-builtin-variables

default: help

install:
	$(call header,Install binary)
	go install
	cp tools/gcp-iam.fish ~/.config/fish/completions/

###############################################################################
# GO Tests
###############################################################################

test: go-tidy test-config test-db test-update ## Run GO tests

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

###############################################################################
# Git Management
###############################################################################

commit: ## Commit Changes
	git commit -m "$(shell date +%Y.%m.%d-%H%M)"

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
