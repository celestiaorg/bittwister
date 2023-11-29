package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRESTApiV1_IndexPage(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	api := NewRESTApiV1(false, logger)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	api.IndexPage(rr, req)
	resp := rr.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status code")
}
