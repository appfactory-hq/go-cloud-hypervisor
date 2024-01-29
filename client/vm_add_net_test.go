// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMAddNet(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VM().AddNet(ctx, &VMAddNetRequest{})
		assert.EqualError(t, err, `failed to call add-net: do request: Put "localhost:badport/api/v1/vm.add-net": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.add-net", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"tap":"tap0","ip":"10.0.20.15","iommu":false,"num_queues":1,"queue_size":128,"vhost_user":false,"id":"foo","pci_segment":2,"rate_limiter_config":{"bandwidth":{"size":1024,"refill_time":60}}}`, string(b))

			err = json.NewEncoder(w).Encode(&VMAddNetResponse{
				VMPCIDeviceInfo: &VMPCIDeviceInfo{
					ID: "foo",
				},
			})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VM().AddNet(ctx, &VMAddNetRequest{
			VMConfigNet: &VMConfigNet{
				Tap:        "tap0",
				IP:         "10.0.20.15",
				NumQueues:  1,
				QueueSize:  128,
				ID:         "foo",
				PCISegment: 2,
				RateLimiterConfig: &VMConfigRateLimiterConfig{
					Bandwidth: &VMConfigRateLimiterBandwidth{
						Size:       1024,
						RefillTime: 60,
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "foo", resp.ID)
	})
}
