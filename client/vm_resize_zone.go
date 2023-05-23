package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMResizeZoneRequest struct {
	ID         string `json:"id"`
	DesiredRAM int64  `json:"desired_ram"`
}

// ResizeZone resize a memory zone.
func (c *VMClient) ResizeZone(ctx context.Context, req *VMResizeZoneRequest) error {
	code, err := c.call(ctx, http.MethodPut, "resize-zone", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call resize-zone: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
