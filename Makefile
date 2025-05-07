app_name := myapp
main_package := ./cmd/api

.PHONY: run build clean

run:
	go run $(main_package)

build:
	go build -o $(app_name) $(main_package)

clean:
	rm -f $(app_name)
