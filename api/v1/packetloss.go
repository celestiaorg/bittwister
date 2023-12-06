package api

import (
	"encoding/json"
	"net/http"

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

	if !ensureServiceInitialized(resp, a.pl) {
		return
	}

	if err := a.pl.SetPacketLossRate(body.PacketLossRate); err != nil {
		sendJSONError(resp,
			MetaMessage{
				Type:    APIMetaMessageTypeError,
				Slug:    SlugServiceSetParamFailed,
				Title:   "Service set param failed",
				Message: err.Error(),
			},
			http.StatusInternalServerError)
		return
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
