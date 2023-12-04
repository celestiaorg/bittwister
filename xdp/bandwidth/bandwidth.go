package bandwidth

import (
	"context"
	"fmt"
	"net"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/cilium/ebpf"
)

type Bandwidth struct {
	NetworkInterface *net.Interface
	Limit            int64 // Bytes per second
}

var _ xdp.XdpLoader = (*Bandwidth)(nil)

func (b *Bandwidth) Start() (xdp.CancelFunc, error) {
	x, err := xdp.GetPreparedXdpObject(b.NetworkInterface.Index)
	if err != nil {
		return nil, fmt.Errorf("prepare XDP object: %w", err)
	}

	key := uint32(0)
	err = x.BpfObjs.BandwidthLimitMap.Update(key, b.Limit, ebpf.UpdateAny)
	if err != nil {
		if cErr := x.Close(); cErr != nil {
			return nil, fmt.Errorf("close XDP object: %w", cErr)
		}
		return nil, fmt.Errorf("update bandwidth limit rate: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
	}()

	cancelFunc := xdp.CancelFunc(func() error {
		// Update the map with a rate of 0 to disable the bandwidth limiter.
		zero := int64(0)
		err := x.BpfObjs.BandwidthLimitMap.Update(key, zero, ebpf.UpdateAny)
		if err != nil {
			return fmt.Errorf("update bandwidth limit rate to zero: %w", err)
		}

		if err := x.Close(); err != nil {
			return fmt.Errorf("close XDP object: %w", err)
		}
		cancel()
		return nil
	})

	return cancelFunc, nil
}
