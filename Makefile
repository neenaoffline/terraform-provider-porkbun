.PHONY: build test testacc install clean fmt lint

# Binary name
BINARY=terraform-provider-porkbun

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

# Version
VERSION?=0.1.0

# OS detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
	OS=linux
endif
ifeq ($(UNAME_S),Darwin)
	OS=darwin
endif

ifeq ($(UNAME_M),x86_64)
	ARCH=amd64
endif
ifeq ($(UNAME_M),arm64)
	ARCH=arm64
endif
ifeq ($(UNAME_M),aarch64)
	ARCH=arm64
endif

# Install path for local development
INSTALL_PATH=~/.terraform.d/plugins/registry.terraform.io/neena/porkbun/$(VERSION)/$(OS)_$(ARCH)

build:
	$(GOBUILD) -o $(BINARY) -v

test:
	$(GOTEST) -v ./...

# Run acceptance tests
# Requires: PORKBUN_API_KEY and PORKBUN_SECRET_API_KEY environment variables
testacc:
	TF_ACC=1 $(GOTEST) -v ./internal/provider -timeout 120m

# Run a specific acceptance test
# Usage: make testacc-one TEST=TestAccDNSRecordResource_A
testacc-one:
	TF_ACC=1 $(GOTEST) -v ./internal/provider -timeout 120m -run $(TEST)

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY) $(INSTALL_PATH)/

clean:
	$(GOCLEAN)
	rm -f $(BINARY)

fmt:
	$(GOFMT) ./...

lint:
	golangci-lint run ./...

# Generate documentation (if using tfplugindocs)
docs:
	go generate ./...

# Run all checks before committing
check: fmt lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the provider binary"
	@echo "  test       - Run unit tests"
	@echo "  testacc    - Run acceptance tests (requires API keys)"
	@echo "  testacc-one TEST=<name> - Run a single acceptance test"
	@echo "  install    - Install provider locally for development"
	@echo "  clean      - Clean build artifacts"
	@echo "  fmt        - Format Go code"
	@echo "  lint       - Run linter"
	@echo "  docs       - Generate documentation"
	@echo "  check      - Run all checks (fmt, lint, test)"
	@echo ""
	@echo "Environment variables for acceptance tests:"
	@echo "  PORKBUN_API_KEY        - Porkbun API key"
	@echo "  PORKBUN_SECRET_API_KEY - Porkbun secret API key"
	@echo "  PORKBUN_TEST_DOMAIN    - Domain to use for testing (must have API access enabled)"
