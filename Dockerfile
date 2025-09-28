FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o chronos-reminder ./cmd/chronos/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/chronos-reminder .

# Copy the assets directory
COPY --from=builder /app/assets ./assets

# Expose port (if needed)
EXPOSE 8080

# Command to run the application
CMD ["./chronos-reminder"]
