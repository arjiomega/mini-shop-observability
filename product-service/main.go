package main

import (
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"

	_ "modernc.org/sqlite"

	productpb "mini-shop/proto/productpb"

	"mini-shop/product-service/product"
)

func main() {
	db, err := sql.Open("sqlite", "products.db")
	if err != nil {
		log.Fatal(err)
	}
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

	log.Println("product service listening on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
