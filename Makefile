# Notifications makefile

VERSION    ?= 0.0.0-dev
PORT       ?= 12345 
GO         ?= go

commands=deps start test dev

# Usage documentation 
usage:
	@printf "Commands: "
	@for command in $(commands); do \
		printf "\033[34m $$command\033[0m,"; done
	@printf "\b \n"
.PHONY: usage

# Install the required dependencies
install.deps:
	@echo "==> Install dep"
	@$(GO) get -u github.com/golang/dep/cmd/dep
	@echo "==> Install dependencies"
	@dep ensure
.PHONY: install.deps

# Install all the dependencies
deps:
	@echo "==> Install dependencies"
	@dep ensure
.PHONY: deps

# Start the service
start:
	@VERSION=$(VERSION) \
	  PORT=$(PORT) \
	  $(GO) run main.go
.PHONY: start

# Build all files.
build:
	@$(GO) build -o server ./...
.PHONY: build

# Build so the binary can be deployed.
build.up:
	@GOOS=linux GOARCH=amd64 $(GO) build -o server ./...
.PHONY: build
# Deploy to up
deploy:
	@echo "==> Deploying..."
	@echo "==> TK"
	@echo "==> Done"
.PHONY: deploy

# Run tests.
test:
	@VERSION=$(VERSION) \
	  PORT=$(PORT) \
	  HOST=$(HOST) \
	  $(GO) test -cover ./... && echo "\n==>\033[32m Ok\033[m\n"
.PHONY: test

# Run all tests.
test_all:
	@VERSION=$(VERSION) \
	  PORT=$(PORT) \
	  HOST=$(HOST) \
	  $(GO) test -tags intergration -v -cover ./... && echo "\n==>\033[32m Ok\033[m\n"
.PHONY: test_all

# Development environment helpers
vim:
	@vim `find . -name '*.go' -not -path './vendor/*' -or -name 'Makefile'`
.PHONY: vim

dev:
	@find . -name '*.go' -not -path './vendor/*' | entr -rcs 'make test && make start'
.PHONY: dev

# Clean.
#clean:
# 	#@rm -fr dist
# 	@git clean -f
#.PHONY: clean