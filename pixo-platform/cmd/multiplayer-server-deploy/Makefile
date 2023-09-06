
all: container

run:
	go run main.go 1 gcr.io/some-image

build:
	go build -o multiplayer-server-deploy

image:
	docker build ../.. -t multiplayer-server-deploy:latest -f ./Dockerfile

container: clean image
	docker run -d --rm -it multiplayer-server-deploy:latest -- ./deploy

clean:
	rm -f multiplayer-server-deploy || true
	docker rmi --force multiplayer-server-deploy:latest || true