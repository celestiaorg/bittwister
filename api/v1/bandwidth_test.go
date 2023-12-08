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

func (s *APITestSuite) TestBandwidthStartStop() {
	t := s.T()

	jsonBody, err := json.Marshal(s.getDefaultBandwidthStartRequest())
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, api.BandwidthPath.Start(), bytes.NewReader(jsonBody))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	s.restAPI.BandwidthStart(rr, req)

	slug, err := getServiceStatusSlug(s.restAPI.BandwidthStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	rr = httptest.NewRecorder()
	s.restAPI.BandwidthStop(rr, nil)
	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	slug, err = getServiceStatusSlug(s.restAPI.BandwidthStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceNotReady, slug)
}

func (s *APITestSuite) TestBandwidthStatus() {
	t := s.T()

	slug, err := getServiceStatusSlug(s.restAPI.BandwidthStatus)
	require.NoError(t, err)
	if slug != api.SlugServiceNotReady && slug != api.SlugServiceNotInitialized {
		t.Fatalf("unexpected service status: %s", slug)
	}

	jsonBody, err := json.Marshal(s.getDefaultBandwidthStartRequest())
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, api.BandwidthPath.Start(), bytes.NewReader(jsonBody))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	s.restAPI.BandwidthStart(rr, req)

	slug, err = getServiceStatusSlug(s.restAPI.BandwidthStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	s.restAPI.BandwidthStop(rr, nil)
	require.Equal(t, http.StatusOK, rr.Code)

	slug, err = getServiceStatusSlug(s.restAPI.BandwidthStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceNotReady, slug)
}

func (s *APITestSuite) getDefaultBandwidthStartRequest() api.BandwidthStartRequest {
	return api.BandwidthStartRequest{
		NetworkInterfaceName: s.ifaceName,
		Limit:                100,
	}
}
