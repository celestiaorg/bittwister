package api

import (
	"encoding/json"
	"net/http"

	"github.com/celestiaorg/bittwister/xdp/packetloss"
	"go.uber.org/zap"
)

// PacketlossStart implements POST /packetloss/start
func (a *RESTApiV1) PacketlossStart(resp http.ResponseWriter, req *http.Request) {
	var body PacketLossStartRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		sendJSONError(resp,
			MetaMessage{
				Type:    APIMetaMessageTypeError,
				Slug:    SlugJSONDecodeFailed,
				Title:   "JSON decode failed",
				Message: err.Error(),
			},
			http.StatusBadRequest)
		return
	}

	if a.pl == nil {
		a.pl = &netRestrictService{
			service: &packetloss.PacketLoss{
				NetworkInterface: nil,
				PacketLossRate:   body.PacketLossRate,
			},
		}
	} else {
		pl, ok := a.pl.service.(*packetloss.PacketLoss)
		if !ok {
			sendJSONError(resp,
				MetaMessage{
					Type:    APIMetaMessageTypeError,
					Slug:    SlugTypeError,
					Title:   "Type cast error",
					Message: "could not cast netRestrictService.service to *packetloss.PacketLoss",
				},
				http.StatusInternalServerError)
			return
		}
		pl.PacketLossRate = body.PacketLossRate
	}

	err := netServiceStart(resp, a.pl, body.NetworkInterfaceName)
	if err != nil {
		a.loggerNoStack.Error("netServiceStart failed", zap.Error(err))
	}
}

// PacketlossStop implements POST /packetloss/stop
func (a *RESTApiV1) PacketlossStop(resp http.ResponseWriter, req *http.Request) {
	if err := netServiceStop(resp, a.pl); err != nil {
		a.loggerNoStack.Error("netServiceStop failed", zap.Error(err))
	}
}

// PacketlossStatus implements GET /packetloss/status
func (a *RESTApiV1) PacketlossStatus(resp http.ResponseWriter, _ *http.Request) {
	if err := netServiceStatus(resp, a.pl); err != nil {
		a.loggerNoStack.Error("netServiceStatus failed", zap.Error(err))
	}
}
