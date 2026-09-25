package product

import (
	"context"
	"log"

	productpb "mini-shop/proto/productpb"

	"google.golang.org/grpc"
)

type ProductClient interface {
	GetProduct(ctx context.Context, id int64) (*productpb.GetProductResponse, error)
	CreateProduct(ctx context.Context, name string) (*productpb.CreateProductResponse, error)
	ListProducts(ctx context.Context) (*productpb.ListProductsResponse, error)
}

type Client struct {
	address string
	client  productpb.ProductServiceClient
	conn    *grpc.ClientConn
}

func NewClient(
	address string,
	conn *grpc.ClientConn,
) *Client {
	return &Client{
		address: address,
		client:  productpb.NewProductServiceClient(conn),
		conn:    conn,
	}
}

func (c *Client) GetProduct(ctx context.Context, id int64) (*productpb.GetProductResponse, error) {
	log.Printf("gRPC → %s GetProduct(%d)", c.address, id)

	return c.client.GetProduct(
		ctx,
		&productpb.GetProductRequest{Id: id},
	)
}

func (c *Client) CreateProduct(ctx context.Context, name string) (*productpb.CreateProductResponse, error) {
	log.Printf("gRPC → %s CreateProduct(%q)", c.address, name)

	return c.client.CreateProduct(
		ctx,
		&productpb.CreateProductRequest{Name: name},
	)
}

func (c *Client) ListProducts(ctx context.Context) (*productpb.ListProductsResponse, error) {
	log.Printf("gRPC → %s ListProducts", c.address)

	return c.client.ListProducts(
		ctx,
		&productpb.ListProductsRequest{},
	)
}
