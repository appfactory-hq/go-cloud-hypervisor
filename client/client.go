package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	base       string
	prefix     string
}

func NewClient(httpClient *http.Client, base string) *Client {
	return &Client{
		httpClient: httpClient,
		base:       strings.TrimRight(base, "/"),
	}
}

func (c *Client) endpoint(path string) string {
	return c.base + "/api/v1/" + c.prefix + strings.TrimLeft(path, "/")
}

func (c *Client) call(ctx context.Context, method string, path string, body interface{}, v interface{}) (code int, err error) {
	var payload io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("marshal body: %w", err)
		}

		payload = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint(path), payload)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if rerr := resp.Body.Close(); rerr != nil {
			err = fmt.Errorf("close response body: %w", rerr)
		}
	}()

	if v == nil {
		return resp.StatusCode, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return resp.StatusCode, fmt.Errorf("decode response: %w", err)
	}

	return resp.StatusCode, nil
}

func (c *Client) expectCode(code, expected int) error {
	if code != expected {
		return fmt.Errorf("unexpected status code: %d", code)
	}

	return nil
}

func (c *Client) VMM() *VMMClient {
	return &VMMClient{
		Client: &Client{
			httpClient: c.httpClient,
			base:       c.base,
			prefix:     "vmm.",
		},
	}
}

func (c *Client) VM() *VMClient {
	return &VMClient{
		Client: &Client{
			httpClient: c.httpClient,
			base:       c.base,
			prefix:     "vm.",
		},
	}
}

type VMMClient struct {
	*Client
}

type VMClient struct {
	*Client
}
