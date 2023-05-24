package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddVSockRequest struct {
	*VMConfigVSOCK
}

type VMAddVSockResponse struct {
	*VMPCIDeviceInfo
}

// AddVSock add a new vsock to the VM.
func (c *VMClient) AddVSock(ctx context.Context, req *VMAddVSockRequest) (*VMAddVSockResponse, error) {
	resp := &VMAddVSockResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-vsock", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-vsock: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
