// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddPMEMRequest struct {
	*VMConfigPMEM
}

type VMAddPMEMResponse struct {
	*VMPCIDeviceInfo
}

// AddPMEM add a new pmem device to the VM.
func (c *VMClient) AddPMEM(ctx context.Context, req *VMAddPMEMRequest) (*VMAddPMEMResponse, error) {
	resp := &VMAddPMEMResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-pmem", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-pmem: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
