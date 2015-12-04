.PHONY: all test build linux

all: clean test build
clean:
	@rm -rf bin/*

build: clean
	@gb build

test:
	@gb test

linux: clean
	@GOOS=linux gb build
