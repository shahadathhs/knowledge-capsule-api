# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o server main.go

# Final lightweight image
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY data ./data

EXPOSE 8080
CMD ["./server"]
