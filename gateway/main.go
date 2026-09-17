package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"mini-shop/gateway/api"
	"mini-shop/gateway/consul"
	"mini-shop/gateway/metrics"
	"mini-shop/gateway/middleware"
	"mini-shop/gateway/product"

	_ "modernc.org/sqlite"
)

func main() {
	metrics.Init()

	consulClient, err := consul.NewClient("consul:8500")
	if err != nil {
		log.Fatal(err)
	}

	var productServiceAddresses []string

	for attempt := 1; attempt <= 30; attempt++ {
		productServiceAddresses, err = consulClient.GetServiceAddresses("product-service")

		if err == nil {
			break
		}

		log.Printf(
			"waiting for product-service in Consul (attempt %d/30): %v",
			attempt,
			err,
		)

		if attempt == 30 {
			log.Fatal("product-service never became available")
		}

		time.Sleep(2 * time.Second)
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

	// MIDDLEWARE(s)
	r.Use(middleware.Metrics())

	// Register OpenAPI-generated routes.
	api.RegisterHandlers(r, productHandler)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start gateway.
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
