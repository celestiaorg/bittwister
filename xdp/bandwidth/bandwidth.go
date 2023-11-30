package bandwidth

import (
	"context"
	"fmt"
	"net"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/cilium/ebpf"
	"go.uber.org/zap"
)

type Bandwidth struct {
	NetworkInterface *net.Interface
	Limit            int64 // Bytes per second
	ready            bool
}

var _ xdp.XdpLoader = (*Bandwidth)(nil)

func (b *Bandwidth) Start(ctx context.Context, logger *zap.Logger) {
	x, err := xdp.GetPreparedXdpObject(b.NetworkInterface.Index)
	if err != nil {
		logger.Error(fmt.Sprintf("Preparing XDP objects: %v", err))
		return
	}
	defer x.Close()

	key := uint32(0)
	err = x.BpfObjs.BandwidthLimitMap.Update(key, b.Limit, ebpf.UpdateAny)
	if err != nil {
		logger.Error(fmt.Sprintf("could not update bandwidth limit rate: %v", err))
		return
	}

	logger.Info(
		fmt.Sprintf("Bandwidth limiter started with rate %d bps on device %q",
			b.Limit,
			b.NetworkInterface.Name,
		),
	)

	b.ready = true
	<-ctx.Done()

	// Update the map with a rate of 0 to disable the bandwidth limiter.
	zero := int64(0)
	err = x.BpfObjs.BandwidthLimitMap.Update(key, zero, ebpf.UpdateAny)
	if err != nil {
		logger.Error(fmt.Sprintf("could not update bandwidth limit rate to zero: %v", err))
		return
	}

	b.ready = false
	logger.Info(fmt.Sprintf("Bandwidth limiter stopped on device %q", b.NetworkInterface.Name))
}

func (b *Bandwidth) Ready() bool {
	return b.ready
}
