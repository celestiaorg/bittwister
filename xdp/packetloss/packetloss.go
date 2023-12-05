package packetloss

import (
	"context"
	"fmt"
	"net"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/cilium/ebpf"
)

type PacketLoss struct {
	NetworkInterface *net.Interface
	PacketLossRate   int32
}

var _ xdp.XdpLoader = (*PacketLoss)(nil)

func (p *PacketLoss) Start() (xdp.CancelFunc, error) {
	x, err := xdp.GetPreparedXdpObject(p.NetworkInterface.Index)
	if err != nil {
		return nil, fmt.Errorf("prepare XDP object: %w", err)
	}

	key := uint32(0)
	err = x.BpfObjs.PacketlossRateMap.Update(key, p.PacketLossRate, ebpf.UpdateAny)
	if err != nil {
		if cErr := x.Close(); cErr != nil {
			return nil, fmt.Errorf("close XDP object: %w", cErr)
		}
		return nil, fmt.Errorf("update packetloss drop rate: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
	}()

	cancelFunc := xdp.CancelFunc(func() error {
		// Update the map with a rate of 0 to disable the packetloss.
		zero := int32(0)
		err = x.BpfObjs.PacketlossRateMap.Update(key, zero, ebpf.UpdateAny)
		if err != nil {
			return fmt.Errorf("update packetloss drop rate to zero: %v", err)
		}

		if err := x.Close(); err != nil {
			return fmt.Errorf("close XDP object: %w", err)
		}
		cancel()
		return nil
	})

	return cancelFunc, nil
}
