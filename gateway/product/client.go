package product

import (
	"context"

	productpb "mini-shop/proto/productpb"
)

type Client struct {
	grpcClient productpb.ProductServiceClient
}

func NewClient(grpcClient productpb.ProductServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

func (c *Client) GetProduct(
	ctx context.Context,
	id int64,
) (*productpb.GetProductResponse, error) {
	return c.grpcClient.GetProduct(ctx, &productpb.GetProductRequest{
		Id: id,
	})
}

func (c *Client) CreateProduct(
	ctx context.Context,
	name string,
) (*productpb.CreateProductResponse, error) {
	return c.grpcClient.CreateProduct(ctx, &productpb.CreateProductRequest{
		Name: name,
	})
}

func (c *Client) ListProducts(
	ctx context.Context,
) (*productpb.ListProductsResponse, error) {
	return c.grpcClient.ListProducts(
		ctx,
		&productpb.ListProductsRequest{},
	)
}
