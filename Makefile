main_pkg := ./cmd/api

.PHONY: run build lint fmt

run:
	go run $(main_pkg)

build:
	go build -o ./bin/ $(main_pkg)

lint:
	golangci-lint run

fmt:
	go mod tidy
	go fmt ./...
	golangci-lint fmt