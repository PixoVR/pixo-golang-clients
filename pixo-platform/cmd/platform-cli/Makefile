
NAME=pixo
EXECUTABLE=./bin/$(NAME)

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