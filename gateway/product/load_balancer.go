package product

import (
	"context"
	"fmt"
	"sync/atomic"

	productpb "mini-shop/proto/productpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoadBalancer struct {
	clients []*Client
	counter uint64
}

func NewLoadBalancer(addresses []string) (*LoadBalancer, error) {
	clients := make([]*Client, 0, len(addresses))

	for _, address := range addresses {
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"connect to %s: %w",
				address,
				err,
			)
		}

		grpcClient := productpb.NewProductServiceClient(conn)

		clients = append(
			clients,
			NewClient(address, grpcClient),
		)
	}

	return &LoadBalancer{
		clients: clients,
	}, nil
}

func (lb *LoadBalancer) next() *Client {
	index := atomic.AddUint64(&lb.counter, 1)

	return lb.clients[(index-1)%uint64(len(lb.clients))]
}

func (lb *LoadBalancer) GetProduct(
	ctx context.Context,
	id int64,
) (*productpb.GetProductResponse, error) {

	client := lb.next()

	return client.GetProduct(ctx, id)
}

func (lb *LoadBalancer) CreateProduct(
	ctx context.Context,
	name string,
) (*productpb.CreateProductResponse, error) {

	client := lb.next()

	return client.CreateProduct(ctx, name)
}

func (lb *LoadBalancer) ListProducts(
	ctx context.Context,
) (*productpb.ListProductsResponse, error) {

	client := lb.next()

	return client.ListProducts(ctx)
}
