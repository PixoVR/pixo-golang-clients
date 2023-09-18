FROM golang:1.20-alpine as builder

WORKDIR /build

COPY . .

RUN cd cmd/platform-cli && go build -o deploy


FROM alpine:latest

WORKDIR /client

COPY --from=builder /build/cmd/platform-cli/deploy /usr/bin/deploy

ENTRYPOINT ["/usr/bin/deploy"]
