
NAME = pixo-dev
EXECUTABLE = ./bin/$(NAME)

all: lint test build cp

dev: build cp

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

container:
	docker run --rm -it ${NAME} --version

container-debug:
	docker run --rm -it --entrypoint /bin/bash ${NAME}

docker: image container

clean:
	rm -rf ${EXECUTABLE} || true
	docker rmi --force ${NAME}:latest || true
