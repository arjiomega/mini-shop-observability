package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	_ "github.com/jackc/pgx/v5/stdlib"

	productpb "mini-shop/proto/productpb"

	"mini-shop/product-service/consul"
	"mini-shop/product-service/product"
	"mini-shop/product-service/tracing"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	serviceID := hostname
	serviceAddress := hostname

	ctx := context.Background()

	shutdownTracing, err := tracing.Init(ctx, "product-service", hostname)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			log.Printf("failed to shutdown tracing: %v", err)
		}
	}()

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

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	productpb.RegisterProductServiceServer(
		grpcServer,
		handler,
	)

	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(
		grpcServer,
		healthServer,
	)

	healthServer.SetServingStatus(
		"",
		healthpb.HealthCheckResponse_SERVING,
	)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
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
