FROM golang:1.21.4 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/platform-cli/main.go


FROM gcr.io/distroless/static-debian12

WORKDIR /client

COPY --from=builder /build/main /usr/bin/pixo

ENTRYPOINT ["/usr/bin/pixo"]
