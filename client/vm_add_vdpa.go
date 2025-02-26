// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddVDPARequest struct {
	*VMConfigVDPA
}

type VMAddVDPAResponse struct {
	*VMPCIDeviceInfo
}

// AddVDPA add a new vdpa device to the VM.
func (c *VMClient) AddVDPA(ctx context.Context, req *VMAddVDPARequest) (*VMAddVDPAResponse, error) {
	resp := &VMAddVDPAResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-vdpa", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-vdpa: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
