# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o mailhub-admin \
    ./cmd/mailhub-admin

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies (ca-certificates for HTTPS, tzdata for time)
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -u 1000 appuser

# Copy binary
COPY --from=builder /app/mailhub-admin /usr/local/bin/

# Copy web assets
COPY --from=builder /app/web /web

# Create data directory
RUN mkdir -p /data && chown appuser:appuser /data

USER appuser

EXPOSE 8080

ENV PORT=8080

CMD ["mailhub-admin"]
