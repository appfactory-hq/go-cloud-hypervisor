package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMRemoveDeviceRequest struct {
	ID string `json:"id"`
}

// RemoveDevice from the VM.
func (c *VMClient) RemoveDevice(ctx context.Context, req *VMRemoveDeviceRequest) error {
	code, err := c.call(ctx, http.MethodPut, "remove-device", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call remove-device: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
