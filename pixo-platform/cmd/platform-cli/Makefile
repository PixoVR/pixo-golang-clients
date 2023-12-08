
NAME = pixo-dev
EXECUTABLE = ./bin/$(NAME)
MODULE_ID ?= 271
SERVER_VERSION ?= "1.03.01"

all: build cp

cp:
	sudo cp $(EXECUTABLE) /usr/local/bin/$(NAME)

lint:
	golangci-lint run

test:
	ginkgo -r -v

run:
	go run main.go

build:
	go build -o $(EXECUTABLE) main.go

image:
	docker build ../.. -t ${NAME} -f ./Dockerfile

container:
	docker run -d --rm -it ${NAME}

clean:
	rm -rf ${EXECUTABLE} || true
	docker rmi --force ${NAME}:latest || true

auth: build
	$(EXECUTABLE) -d auth login

deploy: build
	$(EXECUTABLE) -d -m ${MODULE_ID} mp serverVersions deploy --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest

matchmake: build
	$(EXECUTABLE) -d -m ${MODULE_ID} mp matchmake --server-version ${SERVER_VERSION}

mockserver: build
	$(EXECUTABLE) -d mp matchmake mockserver

