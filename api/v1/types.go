package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RESTApiV1 struct {
	router        *mux.Router
	server        *http.Server
	logger        *zap.Logger
	loggerNoStack *zap.Logger

	pl, bw, lt *netRestrictService

	productionMode bool
}

type PacketLossStartRequest struct {
	NetworkInterfaceName string `json:"network_interface"`
	PacketLossRate       int32  `json:"packet_loss_rate"`
}

type BandwidthStartRequest struct {
	NetworkInterfaceName string `json:"network_interface"`
	Limit                int64  `json:"limit"`
}

type LatencyStartRequest struct {
	NetworkInterfaceName string `json:"network_interface"`
	Latency              int64  `json:"latency_ms"`
	Jitter               int64  `json:"jitter_ms"`
}
