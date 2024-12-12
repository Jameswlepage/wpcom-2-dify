# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Install build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Final runtime image
FROM alpine:3.17
WORKDIR /app

# Add ca-certificates if needed for HTTPS requests
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /app/server /app/server
EXPOSE 8080

# Environment variables can be set via docker-compose or environment
CMD ["./server"]