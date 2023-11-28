FROM golang:1.21.4 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/platform-cli/main.go


FROM alpine:latest

WORKDIR /client

COPY --from=builder /build/main /usr/bin/pixo-platform-cli

ENTRYPOINT ["/usr/bin/pixo-platform-cli"]
