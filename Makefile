#!/bin/bash
PROJECT="GAME"

SOURCE ?= $(shell find . -type f -name '*.go' -not -path '*/generated/*')

all: build

# gps is short for go project skeleton
build:
	go build -v -o tic-tac-toe .

test:
	go test -v -tags="sqlite json1" ./...
	@echo "===\033[32m OK \033[0m==="

fmt:
	@diff=$$(gofmt -d $(SOURCE)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;
