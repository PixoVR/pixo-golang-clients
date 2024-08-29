
all: lint coverage

.PHONY: update
update:
	@echo "🔄 Updating Pixo utilities..."
	@go get github.com/PixoVR/pixo-golang-server-utilities/pixo-platform@0.1.98
	@go mod tidy


.PHONY: lint
lint:
	@echo "🪮 Linting..."
	@golangci-lint run --timeout 5m


.PHONY: unit
unit:
	@echo "🧪 Running tests with coverage..."
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...


.PHONY: calc
calc:
	@echo "🧪 Checking test coverage threshold..."
	@go-test-coverage --config=./.coverage.yaml


coverage: unit calc
cov: coverage


.PHONY: report
report:
	@echo "📊 Generating coverage report..."
	@go test -cover -v -coverprofile=prof.out ./...
	@go tool cover -html=prof.out

