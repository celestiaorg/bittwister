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
	return c.postServiceAction("/packetloss/start", req)
}

func (c *Client) PacketlossStop() error {
	return c.postServiceAction("/packetloss/stop", nil)
}

func (c *Client) PacketlossStatus() (*MetaMessage, error) {
	return c.getServiceStatus("/packetloss/status")
}

func (c *Client) BandwidthStart(req BandwidthStartRequest) error {
	return c.postServiceAction("/bandwidth/start", req)
}

func (c *Client) BandwidthStop() error {
	return c.postServiceAction("/bandwidth/stop", nil)
}

func (c *Client) BandwidthStatus() (*MetaMessage, error) {
	return c.getServiceStatus("/bandwidth/status")
}

func (c *Client) LatencyStart(req LatencyStartRequest) error {
	return c.postServiceAction("/latency/start", req)
}

func (c *Client) LatencyStop() error {
	return c.postServiceAction("/latency/stop", nil)
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
