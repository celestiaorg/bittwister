package sdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SDK_Client_GetResource_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		respBody := []byte(`{"type": "info", "slug": "test", "title": "Test", "message": "Success"}`)
		_, err := w.Write(respBody)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	response, err := client.getResource("/test")

	require.NoError(t, err)
	assert.NotNil(t, response)

	var message MetaMessage
	err = json.Unmarshal(response, &message)
	require.NoError(t, err)
	assert.Equal(t, "info", message.Type)
	assert.Equal(t, "test", message.Slug)
	assert.Equal(t, "Test", message.Title)
	assert.Equal(t, "Success", message.Message)
}

func Test_SDK_Client_GetResource_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	res, err := client.getResource("/error")
	assert.Empty(t, res)
	assert.Error(t, err)
}

func Test_SDK_Client_PostResource_Success(t *testing.T) {
	testCases := []struct {
		RequestBody interface{}
		Expected    MetaMessage
	}{
		{
			RequestBody: struct {
				Type    string `json:"type"`
				Slug    string `json:"slug"`
				Title   string `json:"title"`
				Message string `json:"message"`
			}{
				Type:    "info",
				Slug:    "test",
				Title:   "Test",
				Message: "Success",
			},
			Expected: MetaMessage{
				Type:    "info",
				Slug:    "test",
				Title:   "Test",
				Message: "Success",
			},
		},
	}

	for _, tc := range testCases {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			requestBodyJSON, err := json.Marshal(tc.Expected)
			require.NoError(t, err)
			_, err = w.Write(requestBodyJSON)
			require.NoError(t, err)
		}))
		defer mockServer.Close()

		client := NewClient(mockServer.URL)
		response, err := client.postResource("/test", tc.RequestBody)

		require.NoError(t, err)
		assert.NotNil(t, response)

		var message MetaMessage
		err = json.Unmarshal(response, &message)
		require.NoError(t, err)
		assert.Equal(t, tc.Expected, message)
	}
}

func Test_SDK_Client_PostResource_HTTPError(t *testing.T) {
	testCases := []struct {
		StatusCode int
	}{
		{StatusCode: http.StatusNotFound},
		{StatusCode: http.StatusBadRequest},
		{StatusCode: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.StatusCode)
		}))
		defer mockServer.Close()

		client := NewClient(mockServer.URL)
		res, err := client.postResource("/error", nil)

		assert.Empty(t, res)
		assert.Error(t, err)
	}
}

func Test_SDK_Client_GetServiceStatus_Success(t *testing.T) {
	expectedStatus := MetaMessage{
		Type:    "info",
		Slug:    "service-ready",
		Title:   "Service Status",
		Message: "Service is ready",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		jsonBytes, err := json.Marshal(expectedStatus)
		require.NoError(t, err)

		_, err = w.Write(jsonBytes)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	status, err := client.getServiceStatus("/service/status")

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, expectedStatus, *status)
}

func Test_SDK_Client_GetServiceStatus_Error(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	status, err := client.getServiceStatus("/error/service/status")

	assert.Error(t, err)
	assert.Nil(t, status)
}

func Test_SDK_Client_PacketlossStart_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	err := client.PacketlossStart(PacketLossStartRequest{})

	assert.NoError(t, err)
}

func Test_SDK_Client_PacketlossStop_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()
	client := NewClient(mockServer.URL)
	err := client.PacketlossStop()
	assert.NoError(t, err)
}

func Test_SDK_Client_PacketlossStatus_Success(t *testing.T) {
	expectedStatus := MetaMessage{
		Type:    "info",
		Slug:    "service-ready",
		Title:   "Packetloss Service",
		Message: "Packetloss service is ready",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		jsonBytes, err := json.Marshal(expectedStatus)
		require.NoError(t, err)

		_, err = w.Write(jsonBytes)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	status, err := client.PacketlossStatus()

	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, *status)
}

func Test_SDK_Client_BandwidthStart_Success(t *testing.T) {
	expectedRequest := BandwidthStartRequest{
		NetworkInterfaceName: "eth0",
		Limit:                100,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	err := client.BandwidthStart(expectedRequest)

	assert.NoError(t, err)
}

func Test_SDK_Client_BandwidthStop_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	err := client.BandwidthStop()

	assert.NoError(t, err)
}

func Test_SDK_Client_BandwidthStatus_Success(t *testing.T) {
	expectedStatus := MetaMessage{
		Type:    "info",
		Slug:    "service-ready",
		Title:   "Bandwidth Service",
		Message: "Bandwidth service is ready",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		jsonBytes, err := json.Marshal(expectedStatus)
		require.NoError(t, err)

		_, err = w.Write(jsonBytes)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	status, err := client.BandwidthStatus()

	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, *status)
}

func Test_SDK_Client_LatencyStart_Success(t *testing.T) {
	expectedRequest := LatencyStartRequest{
		NetworkInterfaceName: "eth0",
		Latency:              50,
		Jitter:               10,
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	err := client.LatencyStart(expectedRequest)

	assert.NoError(t, err)
}

func Test_SDK_Client_LatencyStop_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	err := client.LatencyStop()

	assert.NoError(t, err)
}

func Test_SDK_Client_LatencyStatus_Success(t *testing.T) {
	expectedStatus := MetaMessage{
		Type:    "info",
		Slug:    "service-ready",
		Title:   "Latency Service",
		Message: "Latency service is ready",
	}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		jsonBytes, err := json.Marshal(expectedStatus)
		require.NoError(t, err)

		_, err = w.Write(jsonBytes)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	status, err := client.LatencyStatus()

	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, *status)
}

func Test_SDK_Client_AllServicesStatus_Success(t *testing.T) {
	expectedOutput := []ServiceStatus{{
		Name:                 "test-service",
		Ready:                true,
		NetworkInterfaceName: "eth0",
		Params:               map[string]interface{}{"key": "value"},
	}}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		jsonBytes, err := json.Marshal(expectedOutput)
		require.NoError(t, err)

		_, err = w.Write(jsonBytes)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL)
	statuses, err := client.AllServicesStatus()

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, statuses)
}

func Test_SDK_Client_ServiceStatus_Unmarshal(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected ServiceStatus
	}{
		{
			Name:  "Valid service status",
			Input: `{"name": "test-service", "ready": true, "network_interface_name": "eth0", "params": {"key": "value"}}`,
			Expected: ServiceStatus{
				Name:                 "test-service",
				Ready:                true,
				NetworkInterfaceName: "eth0",
				Params:               map[string]interface{}{"key": "value"},
			},
		},
		{
			Name:     "Empty service status",
			Input:    `{}`,
			Expected: ServiceStatus{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var result ServiceStatus
			err := json.Unmarshal([]byte(tc.Input), &result)

			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}
