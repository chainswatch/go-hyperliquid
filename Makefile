# Simple Makefile for go-hyperliquid SDK

BINARY_NAME=hyperliquid

.PHONY: all build test integration test-all fmt vet clean

all: build

build:
	go build ./...

# Run unit tests only (default)
test:
	go test -count=1 ./...

# Run integration tests (requires TEST_ADDRESS and TEST_PRIVATE_KEY)
integration:
	go test -count=1 -tags=integration ./...

# Run both unit and integration tests
test-all: test integration

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)
