package consulclient

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/hashicorp/consul/api"
	"github.com/segmentio/ksuid"
)

// NewListener register new listener in consul with default configuration
func (c *Client) NewListener(name, addr string, port int, tags ...string) (string, error) {
	cnf, err := getDefaultConfig(name, addr, port, tags...)
	if err != nil {
		return "", err
	}
	return cnf.ID, c.api.Agent().ServiceRegister(mapConfig(cnf))
}

// NewListenerWithConfig register new listener in consul with custom configuration
func (c *Client) NewListenerWithConfig(cnf *Config) (string, error) {
	return cnf.ID, c.api.Agent().ServiceRegister(mapConfig(cnf))
}

// RemoveListener removes a listener from the consul pool
func (c *Client) RemoveListener(id string) error {
	return c.api.Agent().ServiceDeregister(id)
}

// map custom config structure to agent service registration structure
func mapConfig(cnf *Config) *api.AgentServiceRegistration {
	return &api.AgentServiceRegistration{
		ID:                cnf.ID,
		Name:              cnf.Name,
		Address:           cnf.Address,
		Port:              cnf.Port,
		Tags:              cnf.Tags,
		EnableTagOverride: cnf.EnableTagOverride,
		Meta:              cnf.Meta,
		Weights:           cnf.Weights,
		Check:             cnf.Check,
		Checks:            cnf.Checks,
	}
}

func getDefaultConfig(name, addr string, port int, tags ...string) (*Config, error) {
	id := ksuid.New().String()
	name = slug.Make(name)
	return &Config{
		ID:      fmt.Sprintf("%s-%s", name, id),
		Name:    name,
		Address: addr,
		Port:    port,
		Tags:    tags,
	}, nil
}
