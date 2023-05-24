// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMReceiveMigrationRequest struct {
	ReceiverURL string `json:"receiver_url"`
}

// ReceiveMigration receive a VM migration from URL.
func (c *VMClient) ReceiveMigration(ctx context.Context, req *VMReceiveMigrationRequest) error {
	code, err := c.call(ctx, http.MethodPut, "receive-migration", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call receive-migration: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
