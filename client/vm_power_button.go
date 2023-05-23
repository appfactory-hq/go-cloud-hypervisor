package client

import (
	"context"
	"fmt"
	"net/http"
)

// PowerButton trigger a power button in the VM.
func (c *VMClient) PowerButton(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "power-button", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call power-button: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
