.PHONY: test lint

GO ?= go

test:
	$(GO) test ./...

lint:
	golangci-lint run ./...
