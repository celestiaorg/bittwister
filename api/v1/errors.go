package api

import (
	"errors"
)

const (
	APIMetaMessageTypeInfo    = "info"
	APIMetaMessageTypeWarning = "warning"
	APIMetaMessageTypeError   = "error"
)

const (
	SlugServiceAlreadyStarted = "service-already-started"
	SlugServiceStartFailed    = "service-start-failed"
	SlugServiceStopFailed     = "service-stop-failed"
	SlugServiceNotStarted     = "service-not-started"
	SlugServiceNotInitialized = "service-not-initialized"
	SlugServiceReady          = "service-ready"
	SlugServiceNotReady       = "service-not-ready"
	SlugServiceSetParamFailed = "service-set-param-failed"
	SlugJSONDecodeFailed      = "json-decode-failed"
	SlugTypeError             = "type-error"
)

type MetaMessage struct {
	Type    string `json:"type"` // info, warning, error
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

var (
	ErrServiceNotInitialized = errors.New(SlugServiceNotInitialized)
	ErrServiceAlreadyStarted = errors.New(SlugServiceAlreadyStarted)
	ErrServiceNotStarted     = errors.New(SlugServiceNotStarted)
	ErrServiceStopFailed     = errors.New(SlugServiceStopFailed)
	ErrServiceStartFailed    = errors.New(SlugServiceStartFailed)
)

// convert a ApiMetaMessage to map[string]interface{}
func (m MetaMessage) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":    m.Type,
		"slug":    m.Slug,
		"title":   m.Title,
		"message": m.Message,
	}
}
