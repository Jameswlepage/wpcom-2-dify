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
# Copy both binaries into final image
COPY --from=builder /app/server /app/server
COPY --from=builder /app/cli /app/cli
# Copy Air config
COPY .air.toml .
# Copy source code and dependencies for hot reloading
COPY . .
COPY --from=builder /go/pkg /go/pkg

EXPOSE 8080
CMD ["air"]
