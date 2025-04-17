FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory for the entire project
WORKDIR /app

# Copy go module files from the tools directory first for better layer caching
COPY tools/go.mod tools/go.sum ./tools/
WORKDIR /app/tools
RUN go mod download

# Go back to the root directory
WORKDIR /app

# Copy the whole project
COPY . .

# Build the application
WORKDIR /app/tools/validators-monitor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o validator-monitor .

# Use a smaller image for the final build
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/validator-monitor .

# Set executable permissions
RUN chmod +x /app/validator-monitor

# Run as non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose the API port
EXPOSE 8080

# Command to run
ENTRYPOINT ["/app/validator-monitor"]