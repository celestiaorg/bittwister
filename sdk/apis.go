package sdk

import (
	"encoding/json"

	"github.com/celestiaorg/bittwister/api/v1"
)

func (c *Client) PacketlossStart(req api.PacketLossStartRequest) error {
	_, err := c.postResource("/packetloss/start", req)
	return err
}

func (c *Client) PacketlossStop() error {
	_, err := c.postResource("/packetloss/stop", nil)
	return err
}

func (c *Client) PacketlossStatus() (*api.MetaMessage, error) {
	return c.getServiceStatus("/packetloss/status")
}

func (c *Client) BandwidthStart(req api.BandwidthStartRequest) error {
	_, err := c.postResource("/bandwidth/start", req)
	return err
}

func (c *Client) BandwidthStop() error {
	_, err := c.postResource("/bandwidth/stop", nil)
	return err
}

func (c *Client) BandwidthStatus() (*api.MetaMessage, error) {
	return c.getServiceStatus("/bandwidth/status")
}

func (c *Client) LatencyStart(req api.LatencyStartRequest) error {
	_, err := c.postResource("/latency/start", req)
	return err
}

func (c *Client) LatencyStop() error {
	_, err := c.postResource("/latency/stop", nil)
	return err
}

func (c *Client) LatencyStatus() (*api.MetaMessage, error) {
	return c.getServiceStatus("/latency/status")
}

func (c *Client) AllServicesStatus() ([]api.ServiceStatus, error) {
	resp, err := c.getResource("/services/status")
	msgs := []api.ServiceStatus{}

	if err := json.Unmarshal(resp, &msgs); err != nil {
		return nil, err
	}
	return msgs, err
}
