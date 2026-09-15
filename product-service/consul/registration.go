package consul

import (
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
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return nil, err
	}

	return client, nil
}
