package client

import (
	"context"
	"fmt"
	"net/http"
)

// Shutdown the VM instance down.
func (c *VMClient) Shutdown(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "shutdown", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call shutdown: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
