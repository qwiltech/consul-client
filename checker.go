package consulclient

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/segmentio/ksuid"
)

// NewChecker register new health checker in consul
func (c *Client) NewChecker(serviceID, hcAddr string) (string, error) {
	id := fmt.Sprintf("healthcheck-%s", ksuid.New().String())

	checker := &api.AgentCheckRegistration{
		ID:        id,
		Name:      id,
		ServiceID: serviceID,
	}

	checker.CheckID = id
	checker.Name = id
	checker.Interval = "10s"
	checker.Timeout = "1s"
	checker.Method = "GET"
	checker.HTTP = hcAddr
	checker.Header = map[string][]string{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	checker.Status = "warning"

	return checker.ID, c.api.Agent().CheckRegister(checker)
}

// RemoveChecker removes a health checker from the consul pool
func (c *Client) RemoveChecker(id string) error {
	return c.api.Agent().CheckDeregister(id)
}
