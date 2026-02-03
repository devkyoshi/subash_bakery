FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy shared code
COPY shared/ ./shared/

# Copy service code
COPY services/license-service/ ./services/license-service/

# Build the application
WORKDIR /app/services/license-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /license-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /license-service .

EXPOSE 8004

CMD ["./license-service"]
