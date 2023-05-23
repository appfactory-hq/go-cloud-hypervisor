package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMCreateRequest struct {
	*VMConfig
}

// Create the cloud-hypervisor Virtual Machine (VM) instance.
func (c *VMClient) Create(ctx context.Context, req *VMCreateRequest) error {
	code, err := c.call(ctx, http.MethodPut, "create", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call create: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
