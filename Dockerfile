# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies and build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations and security flags
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags '-w -s -extldflags "-static"' \
    -buildmode=pie \
    -trimpath \
    -o main cmd/api/main.go

# Final stage - Production optimized
FROM alpine:3.18

# Install runtime dependencies and security updates
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user for security
RUN adduser -D -s /bin/sh -u 1001 appuser

# Create app directory with proper permissions
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Create logs directory
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (configurable via environment)
EXPOSE 8080

# Health check with timeout and retries
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Add labels for better container management
LABEL maintainer="tolstoyjustin@gmail.com" \
      version="1.1.0" \
      description="Go Voting Blockchain - Phase 1 Security Features" \
      org.opencontainers.image.title="Go Voting Blockchain" \
      org.opencontainers.image.description="Secure blockchain-based voting system with JWT authentication" \
      org.opencontainers.image.version="1.1.0" \
      org.opencontainers.image.authors="tolstoyjustin@gmail.com" \
      org.opencontainers.image.url="https://github.com/tolstoyjustin/go-voting-blockchain" \
      org.opencontainers.image.documentation="https://github.com/tolstoyjustin/go-voting-blockchain/blob/main/README.md" \
      org.opencontainers.image.source="https://github.com/tolstoyjustin/go-voting-blockchain"

# Run the application with proper signal handling
CMD ["./main"]