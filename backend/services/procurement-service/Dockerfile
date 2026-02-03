# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/procurement-service/go.mod services/procurement-service/go.sum ./services/procurement-service/

# Download dependencies
WORKDIR /app/services/procurement-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/procurement-service/ ./services/procurement-service/

# Build
WORKDIR /app/services/procurement-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/procurement-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/procurement-service .

EXPOSE 8005

CMD ["./procurement-service"]
