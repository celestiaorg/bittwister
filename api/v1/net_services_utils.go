package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/celestiaorg/bittwister/xdp"
	"github.com/celestiaorg/bittwister/xdp/bandwidth"
	"github.com/celestiaorg/bittwister/xdp/latency"
	"github.com/celestiaorg/bittwister/xdp/packetloss"
	"go.uber.org/zap"
)

const ServiceStopTimeout = 5 // Seconds

type netRestrictService struct {
	service xdp.XdpLoader
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *zap.Logger
}

func (n *netRestrictService) Start(networkInterfaceName string) error {
	if n.service == nil {
		return ErrServiceNotInitialized
	}

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

	n.ctx, n.cancel = context.WithCancel(context.Background())
	go n.service.Start(n.ctx, n.logger)

	return nil
}

func (n *netRestrictService) Stop() error {
	if n.cancel == nil {
		return ErrServiceNotStarted
	}

	n.cancel()

	if n.service.Ready() {
		return ErrServiceStopFailed
	}
	return nil
}

func netServiceStart(resp http.ResponseWriter, ns *netRestrictService, ifaceName string) error {
	if ns == nil || ns.service == nil {
		sendJSONError(resp, MetaMessage{
			Type:    APIMetaMessageTypeError,
			Slug:    SlugServiceNotInitialized,
			Title:   "Service not initiated",
			Message: "To get the status of the service, it must be started first.",
		}, http.StatusOK)
		return ErrServiceNotInitialized
	}

	if ns.service.Ready() {
		sendJSONError(resp, MetaMessage{
			Type:    APIMetaMessageTypeError,
			Slug:    SlugServiceAlreadyStarted,
			Title:   "Service already started",
			Message: "To start the service again, it must be stopped first.",
		}, http.StatusBadRequest)
		return ErrServiceAlreadyStarted
	}

	if err := ns.Start(ifaceName); err != nil {
		sendJSONError(resp,
			MetaMessage{
				Type:    APIMetaMessageTypeError,
				Slug:    SlugServiceStartFailed,
				Title:   "Service start failed",
				Message: err.Error(),
			},
			http.StatusInternalServerError)
		return err
	}
	return nil
}

func netServiceStop(resp http.ResponseWriter, ns *netRestrictService) error {
	if ns == nil || ns.service == nil {
		sendJSONError(resp, MetaMessage{
			Type:    APIMetaMessageTypeError,
			Slug:    SlugServiceNotInitialized,
			Title:   "Service not initiated",
			Message: "To get the status of the service, it must be started first.",
		}, http.StatusOK)
		return ErrServiceNotInitialized
	}

	ns.cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := ServiceStopTimeout * 1000 / 100
	for range ticker.C {
		timeout--
		if !ns.service.Ready() || timeout <= 0 {
			break
		}
	}

	if ns.service.Ready() {
		sendJSONError(resp, MetaMessage{
			Type:    APIMetaMessageTypeError,
			Slug:    SlugServiceStopFailed,
			Title:   "Service stop failed",
			Message: "The service could not be stopped.",
		}, http.StatusInternalServerError)
		return ErrServiceStopFailed
	}

	err := sendJSON(resp, MetaMessage{
		Type:  APIMetaMessageTypeInfo,
		Slug:  SlugServiceNotReady,
		Title: "Service stopped",
	})

	if err != nil {
		return fmt.Errorf("sendJSON failed: %w", err)
	}

	return nil
}

func netServiceStatus(resp http.ResponseWriter, ns *netRestrictService) error {
	if ns == nil || ns.service == nil {
		sendJSONError(resp, MetaMessage{
			Type:    APIMetaMessageTypeError,
			Slug:    SlugServiceNotInitialized,
			Title:   "Service not initiated",
			Message: "To get the status of the service, it must be started first.",
		}, http.StatusOK)
		return ErrServiceNotInitialized
	}

	statusSlug := SlugServiceNotReady
	if ns.service.Ready() {
		statusSlug = SlugServiceReady
	}

	err := sendJSON(resp, MetaMessage{
		Type:  APIMetaMessageTypeInfo,
		Slug:  statusSlug,
		Title: "Service status",
	})

	if err != nil {
		return fmt.Errorf("sendJSON failed: %w", err)
	}
	return nil
}