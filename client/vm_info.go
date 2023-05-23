package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMInfoResponse struct {
	Config           *VMConfig                `json:"config"`
	State            string                   `json:"state"`
	MemoryActuelSize int                      `json:"memory_actual_size"`
	DeviceTree       map[string]*VMDeviceItem `json:"device_tree"`
}

// Info returns general information about the cloud-hypervisor Virtual Machine (VM) instance.
func (c *VMClient) Info(ctx context.Context) (*VMInfoResponse, error) {
	info := &VMInfoResponse{}

	code, err := c.call(ctx, http.MethodGet, "info", nil, &info)
	if err != nil {
		return nil, fmt.Errorf("failed to call info: %w", err)
	}

	return info, c.expectCode(code, http.StatusOK)
}
