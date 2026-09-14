package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mini-shop/gateway/api"
	"mini-shop/gateway/product"
	productpb "mini-shop/proto/productpb"

	_ "modernc.org/sqlite"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	grpcClient := productpb.NewProductServiceClient(conn)

	// Create gateway product client.
	productClient := product.NewClient(grpcClient)

	// Create HTTP handler.
	productHandler := product.NewHandler(productClient)

	// Create HTTP router.
	r := gin.Default()

	// Register OpenAPI-generated routes.
	api.RegisterHandlers(r, productHandler)

	// Start gateway.
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
