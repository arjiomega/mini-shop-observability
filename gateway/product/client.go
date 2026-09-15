package product

import (
	"context"
	"log"

	productpb "mini-shop/proto/productpb"
)

type ProductClient interface {
	GetProduct(ctx context.Context, id int64) (*productpb.GetProductResponse, error)
	CreateProduct(ctx context.Context, name string) (*productpb.CreateProductResponse, error)
	ListProducts(ctx context.Context) (*productpb.ListProductsResponse, error)
}

type Client struct {
	address    string
	grpcClient productpb.ProductServiceClient
}

func NewClient(address string, grpcClient productpb.ProductServiceClient) *Client {
	return &Client{
		address:    address,
		grpcClient: grpcClient,
	}
}

func (c *Client) GetProduct(ctx context.Context, id int64) (*productpb.GetProductResponse, error) {
	log.Printf("gRPC → %s GetProduct(%d)", c.address, id)

	return c.grpcClient.GetProduct(
		ctx,
		&productpb.GetProductRequest{Id: id},
	)
}

func (c *Client) CreateProduct(ctx context.Context, name string) (*productpb.CreateProductResponse, error) {
	log.Printf("gRPC → %s CreateProduct(%q)", c.address, name)

	return c.grpcClient.CreateProduct(
		ctx,
		&productpb.CreateProductRequest{Name: name},
	)
}

func (c *Client) ListProducts(ctx context.Context) (*productpb.ListProductsResponse, error) {
	log.Printf("gRPC → %s ListProducts", c.address)

	return c.grpcClient.ListProducts(
		ctx,
		&productpb.ListProductsRequest{},
	)
}
