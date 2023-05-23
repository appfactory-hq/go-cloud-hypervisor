package client

import (
	"context"
	"fmt"
	"net/http"
)

// Reboot the VM instance.
func (c *VMClient) Reboot(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "reboot", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call reboot: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
