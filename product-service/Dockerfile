FROM golang:1.26 AS builder

WORKDIR /app

# Copy the shared proto module
COPY proto/ ./proto/

# Copy product-service
COPY product-service/ ./product-service/

WORKDIR /app/product-service

# Download dependencies
RUN go mod download

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service .

# ---------------------------------

FROM alpine:latest

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/product-service/product-service .

EXPOSE 50051

CMD ["./product-service"]
