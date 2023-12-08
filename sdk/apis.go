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
	return c.postServiceAction(api.PacketlossPath.Start(), req)
}

func (c *Client) PacketlossStop() error {
	return c.postServiceAction(api.PacketlossPath.Stop(), nil)
}

func (c *Client) PacketlossStatus() (*MetaMessage, error) {
	return c.getServiceStatus(api.PacketlossPath.Status())
}

func (c *Client) BandwidthStart(req BandwidthStartRequest) error {
	return c.postServiceAction(api.BandwidthPath.Start(), req)
}

func (c *Client) BandwidthStop() error {
	return c.postServiceAction(api.BandwidthPath.Stop(), nil)
}

func (c *Client) BandwidthStatus() (*MetaMessage, error) {
	return c.getServiceStatus(api.BandwidthPath.Status())
}

func (c *Client) LatencyStart(req LatencyStartRequest) error {
	return c.postServiceAction(api.LatencyPath.Start(), req)
}

func (c *Client) LatencyStop() error {
	return c.postServiceAction(api.LatencyPath.Stop(), nil)
}

func (c *Client) LatencyStatus() (*MetaMessage, error) {
	return c.getServiceStatus(api.LatencyPath.Status())
}

func (c *Client) AllServicesStatus() ([]ServiceStatus, error) {
	resp, err := c.getResource(api.ServicesPath.Status())
	if err != nil {
		return nil, err
	}

	msgs := []api.ServiceStatus{}
	if err := json.Unmarshal(resp, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
