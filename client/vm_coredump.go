package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMCoreDumpRequest struct {
	DestinationURL string `json:"destination_url"`
}

// CoreDump create a VM coredump.
func (c *VMClient) CoreDump(ctx context.Context, req *VMCoreDumpRequest) error {
	code, err := c.call(ctx, http.MethodPut, "coredump", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call coredump: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
