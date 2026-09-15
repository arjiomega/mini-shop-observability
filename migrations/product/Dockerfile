FROM golang:1.26

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY . /app/migrations

CMD ["goose", "-dir", "/app/migrations", "postgres", "postgres://product_user:product_password@postgres:5432/product_db?sslmode=disable", "up"]
