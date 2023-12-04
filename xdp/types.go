package xdp

import (
	"fmt"
	"sync"

	"github.com/cilium/ebpf/link"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go bpf kerns/main.c -- -I../headers

type XdpLoader interface {
	Start() (CancelFunc, error)
}

type CancelFunc func() error

type XdpObject struct {
	BpfObjs           bpfObjects
	Link              link.Link
	totalServices     int32
	mu                sync.Mutex
	netInterfaceIndex int
}

var xdpObject XdpObject

func GetPreparedXdpObject(netInterfaceIndex int) (*XdpObject, error) {
	xdpObject.mu.Lock()
	defer xdpObject.mu.Unlock()

	// We add this once, so we know how many services are using this object.
	xdpObject.totalServices++

	if xdpObject.Link != nil && xdpObject.netInterfaceIndex == netInterfaceIndex {
		return &xdpObject, nil
	}
	xdpObject.netInterfaceIndex = netInterfaceIndex

	// Load pre-compiled programs into the kernel.
	err := loadBpfObjects(&xdpObject.BpfObjs, nil)
	if err != nil {
		return nil, fmt.Errorf("could not load XDP program: %w", err)
	}

	xdpObject.Link, err = link.AttachXDP(link.XDPOptions{
		Program:   xdpObject.BpfObjs.XdpMain,
		Interface: netInterfaceIndex,
	})

	if err != nil {
		return nil, fmt.Errorf("could not attach XDP program: %w", err)
	}
	return &xdpObject, nil
}

func (x *XdpObject) Close() error {
	x.mu.Lock()
	defer x.mu.Unlock()

	// The object is actually closed when all services using it are closed.
	x.totalServices--
	if x.totalServices > 0 {
		return nil
	}

	if x.Link != nil {
		if err := x.Link.Close(); err != nil {
			return err
		}
	}
	return x.BpfObjs.Close()
}
