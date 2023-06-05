// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMSnapshotRequest struct {
	DestinationURL string `json:"destination_url"`
}

// Snapshot create a VM snapshot.
func (c *VMClient) Snapshot(ctx context.Context, req *VMSnapshotRequest) error {
	code, err := c.call(ctx, http.MethodPut, "snapshot", req, nil)
	if err != nil {
		return fmt.Errorf("failed to call snapshot: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
