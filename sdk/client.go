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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d output: %q", resp.StatusCode, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d output: %q", resp.StatusCode, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) getServiceStatus(resPath string) (*api.MetaMessage, error) {
	resp, err := c.getResource(resPath)
	if err != nil {
		return nil, err
	}
	msg := &api.MetaMessage{}

	if err := json.Unmarshal(resp, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) postServiceAction(resPath string, req interface{}) error {
	resp, err := c.postResource(resPath, req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	msg := api.MetaMessage{}

	if err := json.Unmarshal(resp, &msg); err != nil {
		return nil // if the response is not a MetaMessage, it's probably a success message
	}
	return Error{Message: msg}
}
