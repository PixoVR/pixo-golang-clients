
all: clean test container

test:
	ginkgo -r -v

run:
	go run main.go 1 gcr.io/some-image

build:
	go build -o deploy

image:
	docker build ../.. -t platform-cli:latest -f ./Dockerfile

container: image
	docker run -d --rm -it platform-cli:latest -- ./deploy

clean:
	rm -f platform-cli || true
	docker rmi --force platform-cli:latest || true