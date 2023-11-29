FROM golang:1.21.4 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/platform-cli/main.go


FROM debian:stable

WORKDIR /client

RUN apt-get update \
        && apt-get upgrade -y \
        && apt-get install ca-certificates -y \
        && update-ca-certificates 2>/dev/null

COPY --from=builder /build/main /usr/bin/pixo

ENTRYPOINT ["/usr/bin/pixo"]
