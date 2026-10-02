.PHONY: run routes build test

run:
	go run ./cmd/api

routes:
	@go run ./cmd/routes $(ARGS)

build:
	go build ./...

test:
	go test ./...
