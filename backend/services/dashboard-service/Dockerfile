FROM golang:1.24-alpine as builder

WORKDIR /app

# Copy shared module first
COPY shared/go.mod shared/go.sum ./shared/
COPY services/dashboard-service/go.mod services/dashboard-service/go.sum ./services/dashboard-service/

WORKDIR /app/services/dashboard-service
RUN go mod download

WORKDIR /app

COPY shared/ ./shared/
COPY services/dashboard-service/ ./services/dashboard-service/

WORKDIR /app/services/dashboard-service
RUN CGO_ENABLED=0 GOOS=linux go build -o dashboard-service cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/services/dashboard-service/dashboard-service .
# COPY --from=builder /app/services/dashboard-service/config/config.yaml ./config/ # If using yaml config

EXPOSE 8008

CMD ["./dashboard-service"]
