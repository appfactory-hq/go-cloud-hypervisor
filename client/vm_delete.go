package client

import (
	"context"
	"fmt"
	"net/http"
)

// Delete the cloud-hypervisor Virtual Machine (VM) instance.
func (c *VMClient) Delete(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "delete", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call delete: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
