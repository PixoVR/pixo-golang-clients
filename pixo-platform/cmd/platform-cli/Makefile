
NAME = pixo-dev
EXECUTABLE = ./bin/$(NAME)
MODULE_ID ?= 271
SERVER_VERSION ?= "1.03.01"

all: lint test build cp

new: build cp

cp:
	sudo cp $(EXECUTABLE) /usr/local/bin/$(NAME)

lint:
	golangci-lint run

test:
	ginkgo -r -v --cover ./...

build:
	go build -o $(EXECUTABLE) main.go

image:
	docker build ../.. -t ${NAME} -f ./Dockerfile

container: image
	docker run -d --rm -it ${NAME} -- auth login

clean:
	rm -rf ${EXECUTABLE} || true
	docker rmi --force ${NAME}:latest || true
