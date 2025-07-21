FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /app/.env ./.env

# Установка tzdata и bash для дебага (можно удалить в production)
RUN apk --no-cache add tzdata bash

EXPOSE 8080

CMD ["/app/bin"]