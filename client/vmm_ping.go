package client

import (
	"context"
	"fmt"
	"net/http"
)

// Ping the VMM to check for API server availability.
func (c *VMMClient) Ping(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodGet, "ping", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call ping: %w", err)
	}

	return c.expectCode(code, http.StatusOK)
}
