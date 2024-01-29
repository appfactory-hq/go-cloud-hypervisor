// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

// Pause a previously booted VM instance.
func (c *VMClient) Pause(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "pause", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call pause: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
