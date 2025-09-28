FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.work go.work.sum ./
COPY shared/go.mod shared/go.sum ./shared/
COPY payment/go.mod payment/go.sum ./payment/
COPY platform/go.mod platform/go.sum ./platform/

# Download dependencies
RUN go work sync

# Copy source code
COPY shared ./shared
COPY payment ./payment
COPY platform ./platform

# Build the application
WORKDIR /app/payment
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/payment/main .

# Copy environment file
COPY payment/.env* ./

CMD ["./main"]
