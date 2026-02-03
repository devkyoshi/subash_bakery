# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/inventory-service/go.mod services/inventory-service/go.sum ./services/inventory-service/

# Download dependencies
WORKDIR /app/services/inventory-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/inventory-service/ ./services/inventory-service/

# Build
WORKDIR /app/services/inventory-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/inventory-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/inventory-service .

EXPOSE 8006

CMD ["./inventory-service"]
