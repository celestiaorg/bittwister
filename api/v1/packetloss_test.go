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

func (s *APITestSuite) TestPacketlossStartStop() {
	t := s.T()

	reqBody := api.PacketLossStartRequest{
		NetworkInterfaceName: s.ifaceName,
		PacketLossRate:       10,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/packetloss/start", bytes.NewReader(jsonBody))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	s.restAPI.PacketlossStart(rr, req)

	slug, err := getServiceStatusSlug(s.restAPI.PacketlossStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	rr = httptest.NewRecorder()
	s.restAPI.PacketlossStop(rr, nil)
	require.Equal(t, http.StatusOK, rr.Code)

	slug, err = getServiceStatusSlug(s.restAPI.PacketlossStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceNotReady, slug)
}

func (s *APITestSuite) TestPacketlossStatus() {
	t := s.T()

	slug, err := getServiceStatusSlug(s.restAPI.PacketlossStatus)
	require.NoError(t, err)
	if slug != api.SlugServiceNotReady && slug != api.SlugServiceNotInitialized {
		t.Fatalf("unexpected service status: %s", slug)
	}

	reqBody := api.PacketLossStartRequest{
		NetworkInterfaceName: s.ifaceName,
		PacketLossRate:       10,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/packetloss/start", bytes.NewReader(jsonBody))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	s.restAPI.PacketlossStart(rr, req)

	slug, err = getServiceStatusSlug(s.restAPI.PacketlossStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	slug, err = getServiceStatusSlug(s.restAPI.PacketlossStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceReady, slug)

	s.restAPI.PacketlossStop(rr, nil)
	require.Equal(t, http.StatusOK, rr.Code)

	slug, err = getServiceStatusSlug(s.restAPI.PacketlossStatus)
	require.NoError(t, err)
	assert.Equal(t, api.SlugServiceNotReady, slug)
}
