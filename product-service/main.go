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

	repository := product.NewRepository(db)

	service := product.NewService(repository)

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

	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		log.Fatal("SERVICE_ID is required")
	}
	serviceAddress := os.Getenv("SERVICE_ADDRESS")
	if serviceAddress == "" {
		log.Fatal("SERVICE_ADDRESS is required")
	}

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
