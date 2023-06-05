// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

// PowerButton trigger a power button in the VM.
func (c *VMClient) PowerButton(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "power-button", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call power-button: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
