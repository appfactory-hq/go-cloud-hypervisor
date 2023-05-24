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

func TestClientVMAddDisk(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VM().AddDisk(ctx, &VMAddDiskRequest{})
		assert.EqualError(t, err, `failed to call add-disk: do request: Put "localhost:badport/api/v1/vm.add-disk": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.add-disk", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"path":"/dev/sda","readonly":false,"direct":false,"iommu":false,"vhost_user":false,"pci_segment":2,"id":"foo"}`, string(b))

			err = json.NewEncoder(w).Encode(&VMAddDiskResponse{})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VM().AddDisk(ctx, &VMAddDiskRequest{
			VMConfigDisk: &VMConfigDisk{
				Path:       "/dev/sda",
				ID:         "foo",
				PCISegment: 2,
				IOMMU:      false,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
