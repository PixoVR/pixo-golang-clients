
NAME = pixo-dev
EXECUTABLE = ./bin/$(NAME)
GO_BIN ?= $$(go env GOPATH)/bin

all: lint test build docker cp

dev: build cp

cp:
	sudo cp $(EXECUTABLE) /usr/local/bin/$(NAME)

lint:
	golangci-lint run --timeout 5m

update:
	@echo "🔄 Updating Pixo utilities..."
	@go get github.com/PixoVR/pixo-golang-server-utilities/pixo-platform@0.1.52
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


build:
	go build -o $(EXECUTABLE) main.go

image:
	docker build .. -t ${NAME} -f ./Dockerfile

container:
	docker run --rm -it ${NAME} --version

container-debug:
	docker run --rm -it --entrypoint /bin/bash ${NAME}

docker: image container

clean:
	rm -rf ${EXECUTABLE} || true
	docker rmi --force ${NAME}:latest || true
