package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LatencyStart implements POST /latency/start
func (a *RESTApiV1) LatencyStart(resp http.ResponseWriter, req *http.Request) {
	var body LatencyStartRequest
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

	if !ensureServiceInitialized(resp, a.lt) {
		return
	}

	err := a.lt.SetLatencyParams(
		time.Duration(body.Latency)*time.Millisecond,
		time.Duration(body.Jitter)*time.Millisecond)
	if err != nil {
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

	err = netServiceStart(resp, a.lt, body.NetworkInterfaceName)
	if err != nil {
		a.loggerNoStack.Error("netServiceStart failed", zap.Error(err))
	}
}

// LatencyStop implements POST /latency/stop
func (a *RESTApiV1) LatencyStop(resp http.ResponseWriter, req *http.Request) {
	if err := netServiceStop(resp, a.lt); err != nil {
		a.loggerNoStack.Error("netServiceStop failed", zap.Error(err))
	}
}

// LatencyStatus implements GET /latency/status
func (a *RESTApiV1) LatencyStatus(resp http.ResponseWriter, _ *http.Request) {
	if err := netServiceStatus(resp, a.lt); err != nil {
		a.loggerNoStack.Error("netServiceStatus failed", zap.Error(err))
	}
}
