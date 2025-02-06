FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o /app/main ./cmd/api-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
