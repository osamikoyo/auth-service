FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/bin/main ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/main ./main

CMD ["./main"]