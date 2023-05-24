package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMMPingResponse struct {
	BuildVersion string `json:"build_version"`
	Version      string `json:"version"`
	PID          int    `json:"pid"`
}

// Ping the VMM to check for API server availability.
func (c *VMMClient) Ping(ctx context.Context) (*VMMPingResponse, error) {
	resp := &VMMPingResponse{}

	code, err := c.call(ctx, http.MethodGet, "ping", nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call ping: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
