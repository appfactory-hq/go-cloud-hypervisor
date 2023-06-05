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

func TestClientVMAddFS(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VM().AddFS(ctx, &VMAddFSRequest{})
		assert.EqualError(t, err, `failed to call add-fs: do request: Put "localhost:badport/api/v1/vm.add-fs": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.add-fs", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"tag":"tag","socket":"/path/to/socket","num_queues":1,"queue_size":128,"pci_segment":2,"id":"foo"}`, string(b))

			err = json.NewEncoder(w).Encode(&VMAddFSResponse{
				VMPCIDeviceInfo: &VMPCIDeviceInfo{
					ID: "foo",
				},
			})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VM().AddFS(ctx, &VMAddFSRequest{
			VMConfigFS: &VMConfigFS{
				Tag:        "tag",
				Socket:     "/path/to/socket",
				NumQueues:  1,
				QueueSize:  128,
				ID:         "foo",
				PCISegment: 2,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "foo", resp.ID)
	})
}
