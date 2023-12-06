package api

import (
	"fmt"
	"net"
	"time"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/celestiaorg/bittwister/xdp/bandwidth"
	"github.com/celestiaorg/bittwister/xdp/latency"
	"github.com/celestiaorg/bittwister/xdp/packetloss"
)

const ServiceStopTimeout = 5 // Seconds

type netRestrictService struct {
	service xdp.XdpLoader
	cancel  xdp.CancelFunc
	ready   bool
}

func (n *netRestrictService) Start(networkInterfaceName string) error {
	if n.service == nil {
		return ErrServiceNotInitialized
	}

	if err := n.SetNetworkInterface(networkInterfaceName); err != nil {
		return fmt.Errorf("set network interface: %w", err)
	}

	var err error
	n.cancel, err = n.service.Start()
	if err != nil {
		return fmt.Errorf("start service: %w", err)
	}
	n.ready = true

	return nil
}

func (n *netRestrictService) Stop() error {
	if n.cancel == nil {
		return ErrServiceNotStarted
	}

	if err := n.cancel(); err != nil {
		return fmt.Errorf("stop service: %w", err)
	}
	n.ready = false
	return nil
}

func (n *netRestrictService) SetBandwidthLimit(limit int64) error {
	if s, ok := n.service.(*bandwidth.Bandwidth); ok {
		s.Limit = limit
		return nil
	}

	return fmt.Errorf("could not cast netRestrictService.service to *bandwidth.Bandwidth")
}

func (n *netRestrictService) SetLatencyParams(delay, jitter time.Duration) error {
	if s, ok := n.service.(*latency.Latency); ok {
		s.Latency = delay
		s.Jitter = jitter
		return nil
	}

	return fmt.Errorf("could not cast netRestrictService.service to *latency.Latency")
}

func (n *netRestrictService) SetPacketLossRate(rate int32) error {
	if s, ok := n.service.(*packetloss.PacketLoss); ok {
		s.PacketLossRate = rate
		return nil
	}

	return fmt.Errorf("could not cast netRestrictService.service to *packetloss.PacketLoss")
}

func (n *netRestrictService) SetNetworkInterface(networkInterfaceName string) error {
	iface, err := net.InterfaceByName(networkInterfaceName)
	if err != nil {
		return fmt.Errorf("lookup network device %q: %v", networkInterfaceName, err)
	}

	if s, ok := n.service.(*packetloss.PacketLoss); ok {
		s.NetworkInterface = iface
	} else if s, ok := n.service.(*bandwidth.Bandwidth); ok {
		s.NetworkInterface = iface
	} else if s, ok := n.service.(*latency.Latency); ok {
		s.NetworkInterface = iface
	} else {
		return fmt.Errorf("could not cast netRestrictService.service to *packetloss.PacketLoss, *bandwidth.Bandwidth or *latency.Latency")
	}
	return nil
}
