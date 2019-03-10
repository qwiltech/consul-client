package consulclient

import (
	"github.com/hashicorp/consul/api"
)

type (
	// Client strucure
	Client struct {
		api *api.Client
	}

	// Config of listener
	Config struct {
		ID                string            `json:",omitempty"`
		Name              string            `json:",omitempty"`
		Tags              []string          `json:",omitempty"`
		Port              int               `json:",omitempty"`
		Address           string            `json:",omitempty"`
		EnableTagOverride bool              `json:",omitempty"`
		Meta              map[string]string `json:",omitempty"`
		Weights           *api.AgentWeights `json:",omitempty"`
		Check             *api.AgentServiceCheck
		Checks            api.AgentServiceChecks
	}
)

// NewClient creates a new consul client
func NewClient(addr string) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = addr
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{api: c}, nil
}
