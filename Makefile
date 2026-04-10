-include .env
export

.PHONY: build test lint format fix check

build:
	go build ./...

test:
	go test -race ./...

test-integration:
	go test -race -tags=integration ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run

format:
	go fmt ./...

fix:
	go fix ./...

check: format
	go vet ./...
	go test -race ./...
	gofmt -l .
