main_pkg := ./cmd/api

.PHONY: run build fmt

run:
	go run $(main_pkg)

build:
	go build -o ./bin/ $(main_pkg)

fmt:
	go fmt ./...