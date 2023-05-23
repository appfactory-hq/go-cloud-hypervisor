package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMResizeRequest struct {
	DesiredVCPUs  int `json:"desired_vcpus,omitempty"`
	DesiredRAM    int `json:"desired_ram,omitempty"`
	DesiredBaloon int `json:"desired_balloon,omitempty"`
}

// Resize the VM.
func (c *VMClient) Resize(ctx context.Context, req *VMResizeRequest) error {
	code, err := c.call(ctx, http.MethodPut, "resize", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call resize: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
