

all: build test

build:
	@echo "Building..."
	@go build -o bin/ ./...

example:
	@echo "Building example..."
	@go build -o bin/example ./examples
	@./bin/example

unit:
	@echo "Running unit tests..."
	@ginkgo vertex

test: example unit
