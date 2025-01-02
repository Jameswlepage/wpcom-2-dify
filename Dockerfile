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
# Copy binaries into /usr/local/bin (owned by root), then we will run as appuser but not overwrite these
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/cli /usr/local/bin/cli
RUN chmod +x /usr/local/bin/server /usr/local/bin/cli

# Copy Air config
COPY .air.toml .
# Copy source code for hot reloading
COPY . .
COPY --from=builder /go/pkg /go/pkg

# Create tmp directory with proper permissions
RUN mkdir -p tmp && chmod 777 tmp

# Set PATH environment variable
ENV PATH="/usr/local/bin:${PATH}"

# Use non-root user
RUN adduser -D appuser
USER appuser

EXPOSE 8080
CMD ["air"]
