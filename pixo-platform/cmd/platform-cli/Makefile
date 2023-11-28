
all: clean test container
EXECUTABLE=./bin/platform-cli

test:
	ginkgo -r -v

run:
	go run main.go mp serverVersions deploy --module-id=1 --image=gcr.io/some-image

build:
	go build -o $(EXECUTABLE) main.go

image:
	docker build ../.. -t platform-cli -f ./Dockerfile

container:
	docker run -d --rm -it platform-cli -- /usr/bin/pixo-platform-cli mp serverVersions deploy --module-id=1 --image=gcr.io/some-image

clean:
	rm -f platform-cli || true
	docker rmi --force platform-cli:latest || true