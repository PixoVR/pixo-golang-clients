
GO_BIN ?= $$(go env GOPATH)/bin

all: lint test

lint:
	golangci-lint run --timeout 5m

update:
	@echo "🔄 Updating Pixo utilities..."
	@go get github.com/PixoVR/pixo-golang-server-utilities/pixo-platform@0.1.53
	@go mod tidy

test:
	@echo "🧪 Running tests with coverage..."
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

coverage: test
	@echo "🧪 Checking test coverage threshold..."
	@${GO_BIN}/go-test-coverage --config=./.coverage.yaml

report:
	@echo "📊 Generating coverage report..."
	@IS_LOCAL=true go test -cover -v -coverprofile=prof.out ./...
	@go tool cover -html=prof.out
