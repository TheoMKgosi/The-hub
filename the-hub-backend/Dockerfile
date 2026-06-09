# Stage 1: Building
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o the-hub .


# Stage 2: Production
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/the-hub .

EXPOSE 8080

CMD ["./the-hub"]
