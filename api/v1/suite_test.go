package api_test

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/celestiaorg/bittwister/api/v1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type APITestSuite struct {
	suite.Suite

	logger    *zap.Logger
	restAPI   *api.RESTApiV1
	ifaceName string
}

func (s *APITestSuite) SetupSuite() {
	t := s.T()
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	s.logger = logger

	s.restAPI = api.NewRESTApiV1(false, s.logger)

	ifaceName, err := getLoopbackInterfaceName()
	require.NoError(t, err)
	s.ifaceName = ifaceName
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func getServiceStatusSlug(statusFunc func(http.ResponseWriter, *http.Request)) (string, error) {
	rr := httptest.NewRecorder()
	statusFunc(rr, nil)
	if rr.Code != http.StatusOK {
		return "", errors.New("failed to get service status")
	}

	var msg api.MetaMessage
	err := json.NewDecoder(rr.Body).Decode(&msg)
	if err != nil {
		return "", err
	}
	return msg.Slug, nil
}

func getLoopbackInterfaceName() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			return iface.Name, nil
		}
	}

	return "", errors.New("loopback interface not found")
}
