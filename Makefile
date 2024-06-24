.PHONY: all test build linux

all: clean vet test build
clean:
	@rm -rf bin/*

build: clean
	@go build

test:
	@go test

vet:
	@go vet ./...

linux: clean
	@GOOS=linux go build
