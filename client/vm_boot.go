// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

// Boot the previously created VM instance.
func (c *VMClient) Boot(ctx context.Context) error {
	code, err := c.call(ctx, http.MethodPut, "boot", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to call boot: %w", err)
	}

	return c.expectCode(code, http.StatusNoContent)
}
