package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func RegisterService(
	consulAddress string,
	serviceID string,
	serviceName string,
	address string,
	port int,
) (*api.Client, error) {

	config := api.DefaultConfig()
	config.Address = consulAddress

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check: &api.AgentServiceCheck{
			GRPC: fmt.Sprintf("%s:%d", address, port), Interval: "5s",
			Timeout:                        "2s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return nil, err
	}

	return client, nil
}
