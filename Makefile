# Notifications makefile

VERSION    ?= 0.0.0-dev
PORT       ?= 12345 
GO         ?= go

default: help
.PHONY: default

help: # This
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

deps: ## Install dependencies
	@echo "==> Install dependencies"
	@go mod tidy
.PHONY: deps

lint:
	$(call cyan, "Linting...")
	$(call check-dependency,golint)
	@golint ./... | grep -v unexported || true
	@go vet ./... 2>&1 || echo ''
.PHONY: lint

test: lint
	$(call cyan, "Testing...")
	@VERSION=$(VERSION) \
	  PORT=$(PORT) \
	  HOST=$(HOST) \
	  $(GO) test -short -cover ./... && echo "\n==>\033[32m Ok\033[m\n"
.PHONY: test

start: test ## Run locally
	$(call cyan, "Running...")
	@VERSION=$(VERSION) \
	  PORT=$(PORT) \
	  $(GO) run main.go
.PHONY: start

watch: ## Run locally and monitor for changes
	$(call check-dependency,entr)
	@find . -type f -not -path './vendor/*' -a -not -path '*/\.*' -a -not -path './docs/*' \
	| entr -rcs 'make start -f ./Makefile'
.PHONY: watch

define check-dependency
	$(if $(shell command -v $(1)),,$(error Make sure $(1) is installed))
endef

define red
	@tput setaf 1
	@echo $1
	@tput sgr0
endef

define cyan
	@tput setaf 6
	@echo $1
	@tput sgr0
endef
