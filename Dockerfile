# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build both binaries
RUN CGO_ENABLED=0 go build -o server ./cmd/server
RUN CGO_ENABLED=0 go build -o cli ./cmd/cli

# Final runtime image
FROM alpine:3.17
WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Copy both binaries into final image
COPY --from=builder /app/server /app/server
COPY --from=builder /app/cli /app/cli

EXPOSE 8080
CMD ["./server"]
