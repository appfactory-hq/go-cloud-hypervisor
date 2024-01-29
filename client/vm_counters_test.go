// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMCounters(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VM().Counters(ctx)
		assert.EqualError(t, err, `failed to call counters: do request: Get "localhost:badport/api/v1/vm.counters": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.counters", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)

			err := json.NewEncoder(w).Encode(map[string]map[string]int{
				"foo": {
					"bar": 1,
				},
			})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VM().Counters(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp["foo"]["bar"])
	})
}
