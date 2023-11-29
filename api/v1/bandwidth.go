package api

import (
	"encoding/json"
	"net/http"

	"github.com/celestiaorg/bittwister/xdp/bandwidth"
	"go.uber.org/zap"
)

// BandwidthStart implements POST /bandwidth/start
func (a *RESTApiV1) BandwidthStart(resp http.ResponseWriter, req *http.Request) {
	var body BandwidthStartRequest
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

	if a.bw == nil {
		a.bw = &netRestrictService{
			service: &bandwidth.Bandwidth{
				Limit: body.Limit,
			},
			logger: a.logger,
		}
	} else {
		bw, ok := a.bw.service.(*bandwidth.Bandwidth)
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
		bw.Limit = body.Limit
	}

	err := netServiceStart(resp, a.bw, body.NetworkInterfaceName)
	if err != nil {
		a.loggerNoStack.Error("netServiceStart failed", zap.Error(err))
	}
}

// BandwidthStop implements POST /bandwidth/stop
func (a *RESTApiV1) BandwidthStop(resp http.ResponseWriter, req *http.Request) {
	if err := netServiceStop(resp, a.bw); err != nil {
		a.loggerNoStack.Error("netServiceStop failed", zap.Error(err))
	}
}

// BandwidthStatus implements GET /bandwidth/status
func (a *RESTApiV1) BandwidthStatus(resp http.ResponseWriter, _ *http.Request) {
	if err := netServiceStatus(resp, a.bw); err != nil {
		a.loggerNoStack.Error("netServiceStatus failed", zap.Error(err))
	}
}
