// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

// Shutdown the cloud-hypervisor VMM.
func (c *VMMClient) Shutdown(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "shutdown", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call shutdown: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
