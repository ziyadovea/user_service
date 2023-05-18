# Stage 1: Build stage
FROM golang:1.20-alpine3.18 AS builder

WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o myapp cmd/app/main.go cmd/app/cli.go

# Stage 2: Final stage
FROM alpine:3.18.0

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the previous stage
COPY --from=builder /app/myapp ./

# Set the entry point for the container
CMD ["./myapp", "-config", "/app/configs/dev.yaml"]