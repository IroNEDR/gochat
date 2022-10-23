# syntax=docker/dockerfile:1
FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY tls ./tls/
COPY ui/ ./ui/
RUN ls -la
RUN go build -o ./bin/gochat ./cmd/web/main.go
CMD ["./bin/gochat"]