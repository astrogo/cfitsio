## simple makefile to log workflow
.PHONY: all test clean build

all: build test
	@echo "## bye."

build:
	@echo "build github.com/sbinet/go-cfitsio"
	@go get -v ./...

test: build
	@go test -v

## EOF
