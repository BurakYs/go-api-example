.PHONY: run build fmt

run:
	go run .

build:
	go build -o ./bin/ .

fmt:
	go mod tidy
	go fmt ./...
