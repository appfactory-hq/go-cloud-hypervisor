// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net/http"
)

type VMCountersResponse map[string]map[string]int

// Counters returns counters from the VM.
func (c *VMClient) Counters(ctx context.Context) (VMCountersResponse, error) {
	resp := VMCountersResponse{}

	code, err := c.call(ctx, http.MethodGet, "counters", nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call counters: %w", err)
	}

	return resp, c.expectCode(code, http.StatusOK)
}
