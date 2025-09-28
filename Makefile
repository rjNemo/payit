SHELL := /bin/bash
GOFILES := $(shell find . -path './.modcache' -prune -o -path './.cache' -prune -o -name '*.go' -print)
GOCACHE ?= $(CURDIR)/.cache
GOMODCACHE ?= $(CURDIR)/.modcache
GOFLAGS ?= -count=1

.PHONY: dev fmt lint test clean tidy

dev:
	@echo "Starting development environment"
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) air -c .air.toml

fmt:
	@echo "Formatting Go files"
	@gofmt -w $(GOFILES)

lint:
	@echo "Running linters"
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go vet ./...
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) golangci-lint run --timeout=5m

test:
	@echo "Running tests with coverage"
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test $(GOFLAGS) -cover ./...

clean:
	@rm -rf $(GOCACHE) $(GOMODCACHE)

tidy:
	@echo "Tidying go.mod"
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go mod tidy
