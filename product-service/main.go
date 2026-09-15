package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v5/stdlib"

	productpb "mini-shop/proto/productpb"

	"mini-shop/product-service/consul"
	"mini-shop/product-service/product"
)

func main() {
	db, err := sql.Open(
		"pgx",
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("connected to PostgreSQL")

	defer db.Close()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is required")
	}

	cache := product.NewRedisCache(redisAddr)

	repository := product.NewRepository(db)
	service := product.NewService(repository, cache)
	handler := product.NewHandler(service)

	grpcServer := grpc.NewServer()

	productpb.RegisterProductServiceServer(
		grpcServer,
		handler,
	)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	serviceID := hostname
	serviceAddress := hostname

	_, err = consul.RegisterService(
		"consul:8500",
		serviceID,
		"product-service",
		serviceAddress,
		50051,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("product service registered with Consul")
	log.Println("product service listening on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
