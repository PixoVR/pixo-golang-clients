FROM golang:1.20-alpine as builder

WORKDIR /build

COPY . .

RUN cd cmd/multiplayer-server-deploy && go build -o deploy && ls -la


FROM alpine:latest

WORKDIR /client

COPY --from=builder /build/cmd/multiplayer-server-deploy/deploy .
RUN ls -la

ENTRYPOINT ["./deploy"]
