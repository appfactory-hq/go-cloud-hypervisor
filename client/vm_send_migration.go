package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMSendMigrationRequest struct {
	DestinationURL string `json:"destination_url"`
	Local          bool   `json:"local"`
}

// SendMigration send a VM migration to URL.
func (c *VMClient) SendMigration(ctx context.Context, req *VMSendMigrationRequest) error {
	code, err := c.call(ctx, http.MethodPut, "send-migration", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call send-migration: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
