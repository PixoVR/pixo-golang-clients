
all: container

run:
	go run main.go 1 gcr.io/some-image

build:
	go build -o multiplayer-server-deploy

image:
	docker build ../.. -t multiplayer-server-deploy:latest -f ./Dockerfile

container: image
	docker run -d multiplayer-server-deploy:latest