
NAME = pixo
EXECUTABLE = ./bin/$(NAME)

all: lint build coverage docker cp

dev: build cp

.PHONY: build
build:
	@echo "🚧 Building..."
	@go build -o $(EXECUTABLE) main.go


.PHONY: cp
cp:
	@echo "🚀 Copying binary to /usr/local/bin/pixo-dev..."
	@sudo cp $(EXECUTABLE) /usr/local/bin/pixo-dev


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
	@go tool cover -html=cover.out


.PHONY: mutate
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
	@echo "🧹 Cleaning up..."
	@rm -rf ${EXECUTABLE} || true
	@docker rmi --force ${NAME}:latest > /dev/null 2>&1 || true
	@rm -rf ./cover.out || true
	@rm -rf ./coverprofile.out || true
	@rm -rf .pixo || true
	@mkdir .pixo || true
	@touch .pixo/config.yaml || true
