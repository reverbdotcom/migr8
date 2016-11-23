.PHONY: all test build linux

all: clean vet test build
clean:
	@rm -rf bin/*

build: clean
	@gb build

test:
	@gb test

vet:
	@go vet src/**/*.go

linux: clean
	@GOOS=linux gb build
