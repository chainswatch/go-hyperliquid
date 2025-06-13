# Simple Makefile for go-hyperliquid SDK

BINARY_NAME=hyperliquid

.PHONY: all build test fmt vet clean

all: build

build:
	go build ./...

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)
