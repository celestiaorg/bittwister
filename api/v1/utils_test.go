package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendJson(t *testing.T) {
	testCases := []struct {
		name     string
		resp     *httptest.ResponseRecorder
		obj      interface{}
		expected string
		hasError bool
	}{
		{
			name: "valid input",
			resp: httptest.NewRecorder(),
			obj: map[string]string{
				"name": "Gholi Sibil",
				"age":  "30",
			},
			expected: "{\n  \"age\": \"30\",\n  \"name\": \"Gholi Sibil\"\n}",
			hasError: false,
		},
		{
			name:     "invalid input",
			resp:     httptest.NewRecorder(),
			obj:      make(chan int),
			expected: "",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sendJSON(tc.resp, tc.obj)
			if tc.hasError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.JSONEq(t, tc.expected, tc.resp.Body.String(), "response body should match the expected JSON")
			assert.Equal(t, "application/json", tc.resp.Header().Get("Content-Type"))
		})
	}
}

func TestSendJsonError(t *testing.T) {
	resp := httptest.NewRecorder()
	obj := MetaMessage{
		Type:    APIMetaMessageTypeInfo,
		Slug:    "test",
		Title:   "Test Message",
		Message: "This is a test error message",
	}
	code := http.StatusBadRequest

	sendJSONError(resp, obj, code)

	assert.Equal(t, code, resp.Code, "response code should be equal to the passed code")

	expectedBody := `{
      "type": "info",
      "slug": "test",
      "title": "Test Message",
      "message": "This is a test error message"
    }`
	assert.JSONEq(t, expectedBody, resp.Body.String(), "response body should match the expected JSON")
}
