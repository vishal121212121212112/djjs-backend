# syntax=docker/dockerfile:1
# 1. Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy Go module definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# (Optional) Add this to ensure missing deps are added
RUN go mod tidy

# Build static binary for Linux (optional: disable CGO to allow minimal base)
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -o /app/main ./main.go

# 2. Final stage — minimal runtime
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
