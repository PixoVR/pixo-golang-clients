FROM golang:1.20-alpine

WORKDIR /client

COPY . .

RUN cd cmd/multiplayer-server-update && go build -o multiplayer-server-update

CMD ["./cmd/multiplayer-server-update/multiplayer-server-update"]
