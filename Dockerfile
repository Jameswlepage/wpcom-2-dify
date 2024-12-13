# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build both binaries
RUN CGO_ENABLED=0 go build -o server ./cmd/server
RUN CGO_ENABLED=0 go build -o cli ./cmd/cli

# Final runtime image
FROM golang:1.23-alpine
WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates git && update-ca-certificates

# Copy Air from builder
COPY --from=builder /go/bin/air /usr/local/bin/air
# Copy binaries into /usr/local/bin so they're not hidden by the volume mount
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/cli /usr/local/bin/cli
RUN chmod +x /usr/local/bin/server /usr/local/bin/cli

# Copy Air config
COPY .air.toml .
# Copy source code and dependencies for hot reloading
COPY . .
COPY --from=builder /go/pkg /go/pkg

# Ensure tmp directory exists for Air builds
RUN mkdir -p tmp

# Set PATH environment variable
ENV PATH="/usr/local/bin:${PATH}"

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV POSTGRES_HOST=dify-db
ENV POSTGRES_DB=${POSTGRES_DB}
ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Run the application
CMD ["air"]
