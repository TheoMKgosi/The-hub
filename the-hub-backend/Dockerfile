# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy binaries
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Copy environment files (if they exist)
COPY --from=builder /app/.env* ./

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]