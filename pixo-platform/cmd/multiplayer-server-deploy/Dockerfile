FROM golang:1.20-alpine

WORKDIR /client

COPY . .

RUN cd cmd/multiplayer-server-deploy && go build -o multiplayer-server-deploy

CMD ["./cmd/multiplayer-server-deploy/multiplayer-server-deploy"]
