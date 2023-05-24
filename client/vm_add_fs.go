package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddFSRequest struct {
	*VMConfigFS
}

type VMAddFSResponse struct {
	*VMPCIDeviceInfo
}

// AddFS add a new virtio-fs device to the VM.
func (c *VMClient) AddFS(ctx context.Context, req *VMAddFSRequest) (*VMAddFSResponse, error) {
	resp := &VMAddFSResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-fs", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-fs: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
