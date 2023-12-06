package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	api "github.com/celestiaorg/bittwister/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TestLatencyStartStop() {
	t := s.T()

	jsonBody, err := json.Marshal(s.getDefaultLatencyStartRequest())
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/latency/start", bytes.NewReader(jsonBody))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	s.restAPI.LatencyStart(rr, req)

	slug, err := getServiceStatusSlug(s.restAPI.LatencyStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	rr = httptest.NewRecorder()
	s.restAPI.LatencyStop(rr, nil)
	require.Equal(t, http.StatusOK, rr.Code)

	slug, err = getServiceStatusSlug(s.restAPI.LatencyStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceNotReady, slug)
}

func (s *APITestSuite) TestLatencyStatus() {
	t := s.T()

	slug, err := getServiceStatusSlug(s.restAPI.LatencyStatus)
	require.NoError(t, err)
	if slug != api.SlugServiceNotReady && slug != api.SlugServiceNotInitialized {
		t.Fatalf("unexpected service status: %s", slug)
	}

	jsonBody, err := json.Marshal(s.getDefaultLatencyStartRequest())
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/latency/start", bytes.NewReader(jsonBody))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	s.restAPI.LatencyStart(rr, req)

	slug, err = getServiceStatusSlug(s.restAPI.LatencyStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	slug, err = getServiceStatusSlug(s.restAPI.LatencyStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	s.restAPI.LatencyStop(rr, nil)
	require.Equal(t, http.StatusOK, rr.Code)

	slug, err = getServiceStatusSlug(s.restAPI.LatencyStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceNotReady, slug)
}

func (s *APITestSuite) getDefaultLatencyStartRequest() api.LatencyStartRequest {
	return api.LatencyStartRequest{
		NetworkInterfaceName: s.ifaceName,
		Latency:              100,
		Jitter:               50,
	}
}
