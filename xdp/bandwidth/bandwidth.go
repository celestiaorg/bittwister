package bandwidth

import (
	"context"
	"fmt"
	"net"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"go.uber.org/zap"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go bpf kerns/xdp_bandwidth.c -- -I../headers

type Bandwidth struct {
	NetworkInterface *net.Interface
	Limit            int64 // Bytes per second
	ready            bool
}

var _ xdp.XdpLoader = (*Bandwidth)(nil)

func (b *Bandwidth) Start(ctx context.Context, logger *zap.Logger) {
	// Load pre-compiled programs into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		logger.Error(fmt.Sprintf("loading objects: %v", err))
		return
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpBandwidthLimit,
		Interface: b.NetworkInterface.Index,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("could not attach XDP program: %v", err))
		return
	}
	defer l.Close()

	key := uint32(0)
	err = objs.BandwidthLimitMap.Update(key, b.Limit, ebpf.UpdateAny)
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

	b.ready = false
	logger.Info(fmt.Sprintf("Bandwidth limiter stopped on device %q", b.NetworkInterface.Name))
}

func (b *Bandwidth) Ready() bool {
	return b.ready
}
