include .env

.PHONY: run lint
run:
	go run cmd/main.go

lint:
	golangci-lint run -v ./...