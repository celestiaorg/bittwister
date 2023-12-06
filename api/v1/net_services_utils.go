package api

import (
	"fmt"
	"net/http"
	"time"
)

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

	if ns.ready {
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

	if err := ns.Stop(); err != nil {
		sendJSONError(resp,
			MetaMessage{
				Type:    APIMetaMessageTypeError,
				Slug:    SlugServiceStopFailed,
				Title:   "Service stop failed",
				Message: err.Error(),
			},
			http.StatusInternalServerError)
		return err
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := ServiceStopTimeout * 1000 / 100
	for range ticker.C {
		timeout--
		if !ns.ready || timeout <= 0 {
			break
		}
	}

	if ns.ready {
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
	if ns.ready {
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

func ensureServiceInitialized(resp http.ResponseWriter, ns *netRestrictService) bool {
	if ns != nil {
		return true
	}
	sendJSONError(resp,
		MetaMessage{
			Type:    APIMetaMessageTypeError,
			Slug:    SlugServiceNotInitialized,
			Title:   "Service not initiated",
			Message: "a.(ns *netRestrictService) is nil",
		},
		http.StatusInternalServerError)

	return false
}
