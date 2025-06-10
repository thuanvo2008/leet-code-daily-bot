package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// RequestSender defines an interface for sending HTTP requests
type RequestSender interface {
	SendRequest(req *http.Request) ([]byte, error)
}

// HTTPRequestClient implements RequestSender for HTTP requests
type HTTPRequestClient struct {
	client *http.Client
}

// NewHTTPRequestClient creates a new HTTP request client
func NewHTTPRequestClient(client *http.Client) *HTTPRequestClient {
	if client == nil {
		client = &http.Client{}
	}
	return &HTTPRequestClient{client: client}
}

// SendRequest sends an HTTP request and returns the response body
func (c *HTTPRequestClient) SendRequest(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

// CreateRequest creates a new HTTP request with JSON payload
func CreateRequest(method, url string, payload []byte) (*http.Request, error) {
	return http.NewRequest(method, url, bytes.NewBuffer(payload))
}
