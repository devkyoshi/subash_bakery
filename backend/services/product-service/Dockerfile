# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/product-service/go.mod services/product-service/go.sum ./services/product-service/

# Download dependencies
WORKDIR /app/services/product-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/product-service/ ./services/product-service/

# Build
WORKDIR /app/services/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/product-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/product-service .

EXPOSE 8003

CMD ["./product-service"]