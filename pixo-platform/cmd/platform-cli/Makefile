
NAME = pixo
EXECUTABLE = ./bin/$(NAME)
SERVER_VERSION ?= "1.3.1"

all: build cp

cp:
	sudo cp $(EXECUTABLE) /usr/local/bin/$(NAME)

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
	rm -f ${NAME} || true
	docker rmi --force ${NAME}:latest || true

deploy: build
	$(EXECUTABLE) -d mp serverVersions deploy --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest

matchmake: build
	$(EXECUTABLE) -d mp matchmake --server-version ${SERVER_VERSION}