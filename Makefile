-include .env
export

# Tool versions.
GOLANGCI_LINT_VERSION := v2.14.0
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available targets
	@awk 'BEGIN { FS = ":.*## " } /^[a-zA-Z_-]+:.*## / { printf "  %-14s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build vet fix format lint test example

build: ## Compile every package
	go build ./...

vet: ## Run go vet on every package
	go vet ./...

fix: ## Modernize code with go fix
	go fix ./...

format: ## Format code with gofmt and goimports
	$(GOLANGCI_LINT) fmt

lint: ## Run static analysis with golangci-lint
	$(GOLANGCI_LINT) run

test: ## Run tests with the race detector
	go test -race ./...

example: ## Run the solve example (make example provider=capsolver)
	go run ./examples/solve -provider $(or $(provider),capsolver)
