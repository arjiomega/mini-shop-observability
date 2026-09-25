package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	client *api.Client
}

func NewClient(address string) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetServiceAddresses(
	serviceName string,
	waitIndex uint64,
) ([]string, uint64, error) {

	queryOptions := &api.QueryOptions{
		WaitIndex: waitIndex,
		WaitTime:  5 * time.Minute,
	}

	services, meta, err := c.client.Health().Service(
		serviceName,
		"",
		true,
		queryOptions,
	)
	if err != nil {
		return nil, waitIndex, err
	}

	if len(services) == 0 {
		return nil, meta.LastIndex, fmt.Errorf(
			"service %q not found",
			serviceName,
		)
	}

	addresses := make([]string, 0, len(services))

	for _, service := range services {
		address := fmt.Sprintf(
			"%s:%d",
			service.Service.Address,
			service.Service.Port,
		)

		addresses = append(addresses, address)
	}

	return addresses, meta.LastIndex, nil
}
