include .env

.PHONY: run lint run-env
run:
	go run cmd/main.go

lint:
	golangci-lint run -v ./...

run-env:
	docker-compose up -f build/docker-compose.yaml
