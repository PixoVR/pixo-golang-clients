
NAME=platform-cli
EXECUTABLE=./bin/$(NAME)

all: clean test container

test:
	ginkgo -r -v

run:
	go run main.go mp serverVersions deploy --module-id=1 --image=gcr.io/some-image

build:
	go build -o $(EXECUTABLE) main.go

image:
	docker build ../.. -t ${NAME} -f ./Dockerfile

container:
	docker run -d --rm -it ${NAME} -- /usr/bin/pixo-platform-cli mp serverVersions deploy --module-id=1 --image=gcr.io/some-image

clean:
	rm -f ${NAME} || true
	docker rmi --force ${NAME}:latest || true