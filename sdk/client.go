package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/celestiaorg/bittwister/api/v1"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) getResource(resPath string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.baseURL + resPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var statusErr error
	if resp.StatusCode != http.StatusOK {
		statusErr = fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return body, statusErr
}

func (c *Client) postResource(resPath string, requestBody interface{}) ([]byte, error) {
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+resPath, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var statusErr error
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		statusErr = fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return body, statusErr
}

func (c *Client) getServiceStatus(resPath string) (*api.MetaMessage, error) {
	resp, resErr := c.getResource(resPath)
	if resErr != nil {
		return nil, resErr
	}
	msg := &api.MetaMessage{}

	if err := json.Unmarshal(resp, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) postServiceAction(resPath string, req interface{}) error {
	resp, resErr := c.postResource(resPath, req)
	// Since the response body can have more information about the error, we try to parse it
	if resErr != nil && resp == nil {
		return fmt.Errorf("postResource: %w", resErr)
	}
	msg := api.MetaMessage{}

	if err := json.Unmarshal(resp, &msg); err != nil {
		// if the response is not a MetaMessage, it's probably a success message
		// Therefore we just return the original error from the request
		return resErr
	}
	return Error{Message: msg}
}
