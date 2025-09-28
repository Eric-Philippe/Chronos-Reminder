FROM golang:1.25-alpine AS builder

WORKDIR /app

# Build args passed by Docker Buildx
ARG TARGETOS
ARG TARGETARCH

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary for the target platform
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o chronos-reminder ./cmd/chronos/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/chronos-reminder .

# Copy assets
COPY --from=builder /app/assets ./assets

EXPOSE 8080

CMD ["./chronos-reminder"]
