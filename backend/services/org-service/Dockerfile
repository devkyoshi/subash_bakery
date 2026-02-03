# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/org-service/go.mod services/org-service/go.sum ./services/org-service/

# Download dependencies
WORKDIR /app/services/org-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/org-service/ ./services/org-service/

# Build
WORKDIR /app/services/org-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/org-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/org-service .

EXPOSE 8002

CMD ["./org-service"]
