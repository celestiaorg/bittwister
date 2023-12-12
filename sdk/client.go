package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/celestiaorg/bittwister/api/v1"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	proxy      *url.URL
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) SetProxy(proxyURL string) error {
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	c.proxy = proxy
	c.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(c.proxy)}
	return nil
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
	if err == nil {
		return nil
	}

	if len(resp) == 0 {
		return fmt.Errorf("postResource: %w", err)
	}

	msg := api.MetaMessage{}
	if err := json.Unmarshal(resp, &msg); err != nil {
		return fmt.Errorf("raw output: %s", string(resp))
	}
	return Error{Message: msg}
}
