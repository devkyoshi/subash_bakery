# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/notification-service/go.mod services/notification-service/go.sum ./services/notification-service/

# Download dependencies
WORKDIR /app/services/notification-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/notification-service/ ./services/notification-service/

# Build
WORKDIR /app/services/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/notification-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/notification-service .

EXPOSE 8007

CMD ["./notification-service"]
