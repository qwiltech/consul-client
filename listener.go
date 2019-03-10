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
	cnf, err := getDefaultConfigWithHTTPHealthcheck(name, addr, port, hc, tags...)
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

func getDefaultConfigWithHTTPHealthcheck(name, addr string, port int, hc string, tags ...string) (*Config, error) {
	cnf, err := getDefaultConfig(name, addr, port, tags...)
	if err != nil {
		return nil, err
	}

	cnf.Check = &api.AgentServiceCheck{
		CheckID:  ksuid.New().String(),
		Name:     fmt.Sprintf("%s-healthcheck", name),
		Interval: "10s",
		Timeout:  "1s",
		Method:   "GET",
		HTTP:     fmt.Sprintf("http://%s:%d/health", addr, port),
		Status:   "warning",
		Header: map[string][]string{
			"Content-Type": []string{"application/json"},
			"Accept":       []string{"application/json"},
		},
		DeregisterCriticalServiceAfter: "60s",
	}

	if strings.Contains(hc, "://") {
		cnf.Check.HTTP = hc
	} else if len(hc) > 0 {
		cnf.Check.HTTP = fmt.Sprintf("http://%s:%d/%s", addr, port, strings.TrimLeft(hc, "/"))
	}

	return cnf, nil
}
