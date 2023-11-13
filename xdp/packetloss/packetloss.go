package packetloss

import (
	"context"
	"fmt"
	"net"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"go.uber.org/zap"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go bpf kerns/xdp_drop_percent.c -- -I../headers

type PacketLoss struct {
	NetworkInterface *net.Interface
	PacketLossRate   int32
}

func (p *PacketLoss) Start(ctx context.Context, logger *zap.Logger) {
	// Load pre-compiled programs into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		logger.Error(fmt.Sprintf("loading objects: %v", err))
		return
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpDropPercent,
		Interface: p.NetworkInterface.Index,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("could not attach XDP program: %v", err))
		return
	}
	defer l.Close()

	key := uint32(0)
	err = objs.DropRateMap.Update(key, p.PacketLossRate, ebpf.UpdateAny)
	if err != nil {
		logger.Error(fmt.Sprintf("could not update drop rate: %v", err))
		return
	}

	logger.Info(
		fmt.Sprintf("Packet loss started with rate %d%% on device %q",
			p.PacketLossRate,
			p.NetworkInterface.Name,
		),
	)

	<-ctx.Done()
}
