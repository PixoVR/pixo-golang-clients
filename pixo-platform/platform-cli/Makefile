
NAME = pixo-dev
EXECUTABLE = ./bin/$(NAME)

all: lint test build docker cp

dev: build cp

cp:
	sudo cp $(EXECUTABLE) /usr/local/bin/$(NAME)

lint:
	golangci-lint run --timeout 5m

test:
	@echo "🧪 Running tests with coverage..."
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

coverage: test
	@echo "🧪 Checking test coverage threshold..."
	@go-test-coverage --config=./.coverage.yaml

report:
	@echo "📊 Generating coverage report..."
	@go tool cover -html=cover.out

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
