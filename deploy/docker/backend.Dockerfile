FROM golang:1.23 AS builder
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the backend
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/main_server/main.go

# Create final image with only the built binary
FROM alpine:latest
WORKDIR /root/

# Copy the built backend binary from the builder stage
COPY --from=builder /app/backend .

# Expose backend port
EXPOSE 8080

# Run the backend binary
CMD ["./backend"]
