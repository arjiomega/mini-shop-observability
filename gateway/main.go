package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"mini-shop/gateway/api"
	"mini-shop/gateway/consul"
	"mini-shop/gateway/product"

	_ "modernc.org/sqlite"
)

func main() {
	consulClient, err := consul.NewClient("consul:8500")
	if err != nil {
		log.Fatal(err)
	}
	productServiceAddresses, err := consulClient.GetServiceAddresses(
		"product-service",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"discovered product-service instances: %v",
		productServiceAddresses,
	)

	productClient, err := product.NewLoadBalancer(
		productServiceAddresses,
	)
	if err != nil {
		log.Fatal(err)
	}

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
