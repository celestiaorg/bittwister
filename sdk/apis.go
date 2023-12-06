package sdk

import (
	"encoding/json"

	"github.com/celestiaorg/bittwister/api/v1"
)

type PacketLossStartRequest = api.PacketLossStartRequest
type BandwidthStartRequest = api.BandwidthStartRequest
type LatencyStartRequest = api.LatencyStartRequest
type ServiceStatus = api.ServiceStatus
type MetaMessage = api.MetaMessage

func (c *Client) PacketlossStart(req PacketLossStartRequest) error {
	_, err := c.postResource("/packetloss/start", req)
	return err
}

func (c *Client) PacketlossStop() error {
	_, err := c.postResource("/packetloss/stop", nil)
	return err
}

func (c *Client) PacketlossStatus() (*MetaMessage, error) {
	return c.getServiceStatus("/packetloss/status")
}

func (c *Client) BandwidthStart(req BandwidthStartRequest) error {
	_, err := c.postResource("/bandwidth/start", req)
	return err
}

func (c *Client) BandwidthStop() error {
	_, err := c.postResource("/bandwidth/stop", nil)
	return err
}

func (c *Client) BandwidthStatus() (*MetaMessage, error) {
	return c.getServiceStatus("/bandwidth/status")
}

func (c *Client) LatencyStart(req LatencyStartRequest) error {
	_, err := c.postResource("/latency/start", req)
	return err
}

func (c *Client) LatencyStop() error {
	_, err := c.postResource("/latency/stop", nil)
	return err
}

func (c *Client) LatencyStatus() (*MetaMessage, error) {
	return c.getServiceStatus("/latency/status")
}

func (c *Client) AllServicesStatus() ([]ServiceStatus, error) {
	resp, err := c.getResource("/services/status")
	msgs := []api.ServiceStatus{}

	if err := json.Unmarshal(resp, &msgs); err != nil {
		return nil, err
	}
	return msgs, err
}
