package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"mini-shop/gateway/api"
	"mini-shop/gateway/consul"
	"mini-shop/gateway/logging"
	"mini-shop/gateway/metrics"
	"mini-shop/gateway/middleware"
	"mini-shop/gateway/product"
	"mini-shop/gateway/tracing"

	_ "modernc.org/sqlite"
)

func run() error {
	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	logger := logging.New("gateway", hostname)

	shutdownTracing, err := tracing.Init(ctx, "gateway", hostname)
	if err != nil {
		return err
	}

	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			logger.Error(
				ctx,
				"failed to shutdown tracing",
				slog.Any("error", err),
			)
		}
	}()

	logger.Info(
		ctx,
		"gateway starting",
	)

	metrics.Init()

	consulClient, err := consul.NewClient("consul:8500")
	if err != nil {
		return err
	}

	var productServiceAddresses []string
	var lastIndex uint64

	for attempt := 1; attempt <= 30; attempt++ {
		addresses, index, err := consulClient.GetServiceAddresses(
			"product-service",
			0,
		)

		if err == nil {
			productServiceAddresses = addresses
			lastIndex = index
			break
		}

		logger.Warn(
			ctx,
			"waiting for product-service in Consul",
			slog.Int("attempt", attempt),
			slog.Int("max_attempts", 30),
			slog.Any("error", err),
		)

		if attempt == 30 {
			return fmt.Errorf("product-service never became available: %w", err)
		}

		time.Sleep(2 * time.Second)
	}

	logger.Info(
		ctx,
		"discovered product-service instances",
		slog.Any("instances", productServiceAddresses),
	)

	productClient, err := product.NewLoadBalancer(
		productServiceAddresses,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			addresses, index, err := consulClient.GetServiceAddresses(
				"product-service",
				lastIndex,
			)
			if err != nil {
				logger.Warn(
					ctx,
					"failed to watch product-service",
					slog.Any("error", err),
				)
				time.Sleep(2 * time.Second)
				continue
			}

			lastIndex = index

			logger.Info(
				ctx,
				"product-service instances changed",
				slog.Any("instances", addresses),
			)

			if err := productClient.Update(addresses); err != nil {
				logger.Error(
					ctx,
					"failed to update product-service load balancer",
					slog.Any("error", err),
				)
			}

		}
	}()

	// Create HTTP handler.
	productHandler := product.NewHandler(productClient, logger)

	// Create HTTP router.
	r := gin.Default()

	// MIDDLEWARE(s)
	r.Use(otelgin.Middleware("gateway"))
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Metrics())

	// Register OpenAPI-generated routes.
	api.RegisterHandlers(r, productHandler)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start gateway.
	return r.Run(":8080")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
