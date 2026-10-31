.PHONY: init dev build

init:
	brew install pre-commit
	pre-commit install
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

dev:
	go tool air run

build:
	go build cmd/...

lint:
	golangci-lint run ./...

fmt:
	go tool gofumpt -l -w .
