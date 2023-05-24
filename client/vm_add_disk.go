package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddDiskRequest struct {
	*VMConfigDisk
}

type VMAddDiskResponse struct {
	*VMPCIDeviceInfo
}

// AddDisk add a new disk to the VM.
func (c *VMClient) AddDisk(ctx context.Context, req *VMAddDiskRequest) (*VMAddDiskResponse, error) {
	resp := &VMAddDiskResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-disk", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-disk: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
