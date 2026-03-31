# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o saleapp-backend ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/saleapp-backend .

# Copy environment file (should be overridden via env vars in production)
# Only copy if .env.example exists (use dummy to avoid build error if missing)
RUN if [ -f .env.example ]; then cp .env.example .env; fi

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENV GIN_MODE=release
CMD ["./saleapp-backend"]
