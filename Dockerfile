# 1. Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o reflecta-api ./cmd/server


# 2. Run stage (smaller image)
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/reflecta-api .

EXPOSE 4000

CMD ["./reflecta-api"]
