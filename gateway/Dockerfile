FROM golang:1.26 AS builder

WORKDIR /app

# Copy the shared proto module
COPY proto/ ./proto/

# Copy product-service
COPY gateway/ ./gateway/

WORKDIR /app/gateway

# Download dependencies
RUN go mod download

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .

# ---------------------------------

FROM alpine:latest

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/gateway/gateway .

EXPOSE 50051

CMD ["./gateway"]
