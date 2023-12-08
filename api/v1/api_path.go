package api

const endpointPrefix = "/api/v1"

type serviceEndpointPath struct {
	basePath string
}

func (e *serviceEndpointPath) Status() string {
	return endpointPrefix + e.basePath + "/status"
}

func (e *serviceEndpointPath) Start() string {
	return endpointPrefix + e.basePath + "/start"
}

func (e *serviceEndpointPath) Stop() string {
	return endpointPrefix + e.basePath + "/stop"
}

var (
	PacketlossPath = &serviceEndpointPath{basePath: "/packetloss"}
	BandwidthPath  = &serviceEndpointPath{basePath: "/bandwidth"}
	LatencyPath    = &serviceEndpointPath{basePath: "/latency"}
	ServicesPath   = &serviceEndpointPath{basePath: "/services"}
)
