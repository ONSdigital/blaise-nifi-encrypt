.PHONY: build test lint lint-fix

build:
	go build ./...

test:
	go test ./... -race -cover -covermode=atomic

lint:
	test -z "$(gofmt -l .)"
	golangci-lint run ./... --out-format=colored-line-number

lint-fix:
	gofmt -w .
	golangci-lint run ./... --fix --out-format=colored-line-number
