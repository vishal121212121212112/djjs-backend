# syntax=docker/dockerfile:1

# 1. Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy Go module definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy source (including main.go in root)
COPY . .

# Build static binary for Linux (optional: disable CGO to allow minimal base)
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -o /app/main ./main.go

# 2. Final stage — minimal runtime
FROM alpine:latest
WORKDIR /app

# Optionally install CA certificates if needed
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

# Add a non-root user (for security)
RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
