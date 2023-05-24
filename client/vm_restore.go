package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMRestoreRequest struct {
	SourceURL string `json:"source_url"`
	Prefault  bool   `json:"prefault"`
}

// Restore a VM from a snapshot.
func (c *VMClient) Restore(ctx context.Context, req *VMRestoreRequest) error {
	code, err := c.call(ctx, http.MethodPut, "restore", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call restore: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
