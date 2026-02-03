# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/

# Download dependencies
WORKDIR /app/services/auth-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/auth-service/ ./services/auth-service/

# Build
WORKDIR /app/services/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/auth-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/auth-service .

EXPOSE 8001

CMD ["./auth-service"]
