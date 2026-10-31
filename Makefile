.PHONY: init dev build lint fmt sqlc migrate-up migrate-down

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

sqlc:
	go tool sqlc generate

migrate-up:
	go tool goose up --env=.env.dev

migrate-down:
	go tool goose down --env=.env.dev
