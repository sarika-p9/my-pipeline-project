# Builder stage to compile both binaries (gRPC server and CLI)
FROM golang:1.23 AS builder
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the gRPC server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /grpc_server ./cmd/grpc_server/main.go

# Build the democtl (CLI) binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /democtl/democtl ./cmd/democtl/main.go

# Final stage using a lightweight Alpine image
FROM alpine:latest
WORKDIR /root/

# Copy the two binaries from the builder stage
COPY --from=builder /grpc_server /grpc_server
COPY --from=builder /democtl/democtl /democtl/democtl

# Default entrypoint (gRPC server or CLI)
# If you want the gRPC server to be the default:
ENTRYPOINT ["/grpc_server"]
