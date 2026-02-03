# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/

# Download dependencies
WORKDIR /app/services/api-gateway
RUN go mod download

# Copy source code
WORKDIR /app
COPY services/api-gateway/ ./services/api-gateway/

# Build
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-gateway ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
