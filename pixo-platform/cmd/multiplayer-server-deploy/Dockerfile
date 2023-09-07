FROM golang:1.20-alpine as builder

WORKDIR /build

COPY . .

RUN cd cmd/multiplayer-server-deploy && go build -o deploy


FROM alpine:latest

WORKDIR /client

COPY --from=builder /build/cmd/multiplayer-server-deploy/deploy /usr/bin/deploy

ENTRYPOINT ["/usr/bin/deploy"]
