# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY shared/go.mod shared/go.sum ./shared/
COPY services/report-service/go.mod services/report-service/go.sum ./services/report-service/

# Download dependencies
WORKDIR /app/services/report-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY shared/ ./shared/
COPY services/report-service/ ./services/report-service/

# Build
WORKDIR /app/services/report-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/report-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/report-service .

EXPOSE 8009

CMD ["./report-service"]
