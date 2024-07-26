
NAME = pixo-dev
EXECUTABLE = ./bin/$(NAME)

all: lint build coverage docker cp

dev: build cp

.PHONY: build
build:
	@echo "🚧 Building..."
	@go build -o $(EXECUTABLE) main.go


.PHONY: cp
cp:
	@echo "🚀 Copying binary to /usr/local/bin/..."
	@sudo cp $(EXECUTABLE) /usr/local/bin/$(NAME)


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


.PHONY: test
mutate:
	@echo "👾 Running Mutation Tester..."
	@go-mutesting ./src/...


.PHONY: image
image:
	@echo "🐳 Building Docker image..."
	@docker build .. -t ${NAME} -f ./Dockerfile


.PHONY: container
container:
	@echo "🚀 Running Docker container..."
	@docker run --rm -it ${NAME} --version


.PHONY: container-debug
container-debug:
	@echo "🚀 Running Docker container in debug mode..."
	@docker run --rm -it --entrypoint /bin/bash ${NAME}


.PHONY: docker
docker: image container


.PHONY: clean
clean:
	rm -rf ${EXECUTABLE} || true
	docker rmi --force ${NAME}:latest || true
