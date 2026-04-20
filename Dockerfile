# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /mocker ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /mocker .

# Copy docs for Swagger UI
COPY --from=builder /app/docs ./docs

# Create data directory
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment
ENV DATA_DIR=/data

# Run the application
CMD ["./mocker"]