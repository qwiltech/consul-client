package consulclient

import (
	"fmt"
	"strings"

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

// NewHTTPListenerWithHealthcheck register new listener in consul
// with default configuration and custom healthcheck
func (c *Client) NewHTTPListenerWithHealthcheck(name, addr string, port int, hc string, tags ...string) (string, error) {
	cnf, err := getDefaultConfig(name, addr, port, tags...)
	if err != nil {
		return "", err
	}
	if err := c.api.Agent().ServiceRegister(mapConfig(cnf)); err != nil {
		return "", err
	}

	var healthcheck string
	if strings.Contains(hc, "://") {
		healthcheck = hc
	} else if len(hc) > 0 {
		healthcheck = fmt.Sprintf("http://%s:%d/%s", addr, port, strings.TrimLeft(hc, "/"))
	} else {
		healthcheck = fmt.Sprintf("http://%s:%d/health", addr, port)
	}
	_, err = c.NewChecker(cnf.ID, healthcheck)

	return cnf.ID, err
}

// NewListenerWithConfig register new listener in consul with custom configuration
func (c *Client) NewListenerWithConfig(cnf *Config) (string, error) {
	return cnf.ID, c.api.Agent().ServiceRegister(mapConfig(cnf))
}

// RemoveListener removes a listener from the consul pool
func (c *Client) RemoveListener(id string) error {
	return c.api.Agent().ServiceDeregister(id)
}

// default listener configuration
func getDefaultConfig(name, addr string, port int, tags ...string) (*Config, error) {
	name = slug.Make(name)
	id := fmt.Sprintf("%s-%s", name, ksuid.New().String())
	tags = append(tags, id, name)
	return &Config{
		ID:      id,
		Name:    name,
		Address: addr,
		Port:    port,
		Tags:    tags,
	}, nil
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
