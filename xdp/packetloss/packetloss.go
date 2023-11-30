package packetloss

import (
	"context"
	"fmt"
	"net"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/cilium/ebpf"
	"go.uber.org/zap"
)

type PacketLoss struct {
	NetworkInterface *net.Interface
	PacketLossRate   int32
	ready            bool
}

var _ xdp.XdpLoader = (*PacketLoss)(nil)

func (p *PacketLoss) Start(ctx context.Context, logger *zap.Logger) {
	x, err := xdp.GetPreparedXdpObject(p.NetworkInterface.Index)
	if err != nil {
		logger.Error(fmt.Sprintf("Preparing XDP objects: %v", err))
		return
	}
	defer x.Close()

	key := uint32(0)
	err = x.BpfObjs.PacketlossRateMap.Update(key, p.PacketLossRate, ebpf.UpdateAny)
	if err != nil {
		logger.Error(fmt.Sprintf("could not update packetloss drop rate: %v", err))
		return
	}

	logger.Info(
		fmt.Sprintf("Packetloss started with rate %d%% on device %q",
			p.PacketLossRate,
			p.NetworkInterface.Name,
		),
	)

	p.ready = true
	<-ctx.Done()

	// Update the map with a rate of 0 to disable the packetloss.
	zero := int32(0)
	err = x.BpfObjs.PacketlossRateMap.Update(key, zero, ebpf.UpdateAny)
	if err != nil {
		logger.Error(fmt.Sprintf("could not update packetloss drop rate to zero: %v", err))
		return
	}

	p.ready = false
	logger.Info(fmt.Sprintf("Packetloss stopped on device %q", p.NetworkInterface.Name))
}

func (p *PacketLoss) Ready() bool {
	return p.ready
}
