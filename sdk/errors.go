package sdk

import (
	"encoding/json"

	"github.com/celestiaorg/bittwister/api/v1"
)

type Error struct {
	Message MetaMessage
}

// Error returns the error message as a string
func (e Error) Error() string {
	errJSON, _ := json.MarshalIndent(e.Message, "", "  ")
	return string(errJSON)
}

func IsErrorServiceNotInitialized(err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Message.Slug == api.SlugServiceNotInitialized
}

func IsErrorServiceAlreadyStarted(err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Message.Slug == api.SlugServiceAlreadyStarted
}

func IsErrorServiceNotStarted(err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Message.Slug == api.SlugServiceNotStarted
}

func IsErrorServiceStopFailed(err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Message.Slug == api.SlugServiceStopFailed
}

func IsErrorServiceStartFailed(err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Message.Slug == api.SlugServiceStartFailed
}

func IsErrorServiceNotReady(err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Message.Slug == api.SlugServiceNotReady
}
