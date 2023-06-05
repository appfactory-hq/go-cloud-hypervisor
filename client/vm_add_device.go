// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMAddDeviceRequest struct {
	Path       string `json:"path"`
	IOMMU      bool   `json:"iommu"`
	PCISegment int    `json:"pci_segment"`
	ID         string `json:"id"`
}

type VMAddDeviceResponse struct {
	*VMPCIDeviceInfo
}

// AddDevice add a new device to the VM.
func (c *VMClient) AddDevice(ctx context.Context, req *VMAddDeviceRequest) (*VMAddDeviceResponse, error) {
	resp := &VMAddDeviceResponse{}

	code, err := c.call(ctx, http.MethodPut, "add-device", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call add-device: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
