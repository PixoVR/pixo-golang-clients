FROM golang:1.22.1 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pixo platform-cli/main.go


FROM gcr.io/distroless/static-debian12

WORKDIR /client

COPY --from=builder /build/pixo /usr/bin/pixo

ENTRYPOINT ["/usr/bin/pixo"]
