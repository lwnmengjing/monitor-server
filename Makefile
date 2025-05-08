.PHONY: build

build:
	CGO_ENABLED=0 go build -o application main.go
test:
	go test -v ./... -cover
deps:
	go mod tidy
run:
	go run main.go server
migrate:
	go run main.go migrate
generate:
	go generate ./...