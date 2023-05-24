package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddNetRequest struct {
	*VMConfigNet
}

type VMAddNetResponse struct {
	*VMPCIDeviceInfo
}

// AddNet add a new network device to the VM.
func (c *VMClient) AddNet(ctx context.Context, req *VMAddNetRequest) (*VMAddNetResponse, error) {
	resp := &VMAddNetResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-net", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-net: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
