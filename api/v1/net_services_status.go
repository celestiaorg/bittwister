package api

import (
	"net/http"

	"github.com/celestiaorg/bittwister/xdp/bandwidth"
	"github.com/celestiaorg/bittwister/xdp/latency"
	"github.com/celestiaorg/bittwister/xdp/packetloss"
	"go.uber.org/zap"
)

type ServiceStatus struct {
	Name                 string                 `json:"name"`
	Ready                bool                   `json:"ready"`
	NetworkInterfaceName string                 `json:"network_interface_name"`
	Params               map[string]interface{} `json:"params"` // key:value
}

// NetServicesStatus implements GET /services/status
func (a *RESTApiV1) NetServicesStatus(resp http.ResponseWriter, req *http.Request) {
	out := make([]ServiceStatus, 0, 3)
	for _, ns := range []*netRestrictService{a.pl, a.bw, a.lt} {
		if ns == nil {
			continue
		}

		var (
			params       = make(map[string]interface{})
			name         string
			netIfaceName string
		)

		if s, ok := ns.service.(*packetloss.PacketLoss); ok {
			name = "packetloss"
			params["packet_loss_rate"] = s.PacketLossRate
			netIfaceName = s.NetworkInterface.Name

		} else if s, ok := ns.service.(*bandwidth.Bandwidth); ok {
			name = "bandwidth"
			params["limit"] = s.Limit
			netIfaceName = s.NetworkInterface.Name

		} else if s, ok := ns.service.(*latency.Latency); ok {
			name = "latency"
			params["latency_ms"] = s.Latency.Milliseconds()
			params["jitter_ms"] = s.Jitter.Milliseconds()
			netIfaceName = s.NetworkInterface.Name

		} else {
			sendJSONError(resp,
				MetaMessage{
					Type:    APIMetaMessageTypeError,
					Slug:    SlugTypeError,
					Title:   "Type cast error",
					Message: "could not cast netRestrictService.service to *packetloss.PacketLoss, *bandwidth.Bandwidth or *latency.Latency",
				},
				http.StatusInternalServerError)
			return
		}

		out = append(out, ServiceStatus{
			Name:                 name,
			Ready:                ns.ready,
			NetworkInterfaceName: netIfaceName,
			Params:               params,
		})
	}

	if err := sendJSON(resp, out); err != nil {
		a.loggerNoStack.Error("sendJSON failed", zap.Error(err))
	}
}
