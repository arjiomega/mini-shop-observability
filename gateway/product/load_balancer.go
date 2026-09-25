package product

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	productpb "mini-shop/proto/productpb"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoadBalancer struct {
	mu      sync.RWMutex
	clients []*Client
	counter uint64
}

func NewLoadBalancer(addresses []string) (*LoadBalancer, error) {
	lb := &LoadBalancer{}

	if err := lb.Update(addresses); err != nil {
		return nil, err
	}

	return lb, nil
}

func (lb *LoadBalancer) next() *Client {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.clients) == 0 {
		return nil
	}

	index := atomic.AddUint64(&lb.counter, 1)

	return lb.clients[(index-1)%uint64(len(lb.clients))]
}

func (lb *LoadBalancer) Update(addresses []string) error {
	clients := make([]*Client, 0, len(addresses))

	for _, address := range addresses {
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			return fmt.Errorf(
				"create connection to %s: %w",
				address,
				err,
			)
		}

		clients = append(
			clients,
			NewClient(address, conn),
		)
	}

	lb.mu.Lock()
	oldClients := lb.clients
	lb.clients = clients
	lb.mu.Unlock()

	// Close connections that are no longer used.
	for _, client := range oldClients {
		_ = client.conn.Close()
	}

	return nil
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
