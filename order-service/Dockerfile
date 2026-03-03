FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o order-service ./cmd/app

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/order-service .

COPY --from=builder /app/internal/infrastructure/web ./internal/infrastructure/web

EXPOSE 8081
CMD ["./order-service"]