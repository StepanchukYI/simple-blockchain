BIN_NAME := $(or $(PROJECT_NAME),'blockchain')
.PHONY: all build clean test cover

all: build test

dep: # Download required dependencies
	go mod tidy
	go mod download
	go mod vendor

test: dep ## Run unit tests
	go test -race -count=1 -timeout=5s -v ./...

cover:
	go test -timeout=5m -cover -v ./...

build: dep ## Build the binary file
	go build -o ./bin/${BIN_NAME} -a .

clean: ## Remove previous build
	rm -f bin/$(BIN_NAME)
