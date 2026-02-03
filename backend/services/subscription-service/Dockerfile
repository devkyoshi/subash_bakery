# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/subscription-service/go.mod services/subscription-service/go.sum ./services/subscription-service/

# Download dependencies
WORKDIR /app/services/subscription-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/subscription-service/ ./services/subscription-service/

# Build
WORKDIR /app/services/subscription-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/subscription-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/subscription-service .

EXPOSE 8003

CMD ["./subscription-service"]
