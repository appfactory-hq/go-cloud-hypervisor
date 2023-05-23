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

func TestClientVMAddDevice(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := NewClient(http.DefaultClient, "localhost:badport")

		_, err := c.VM().AddDevice(ctx, &VMAddDeviceRequest{})
		assert.EqualError(t, err, `failed to call add-device: do request: Put "localhost:badport/api/v1/vm.add-device": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.add-device", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"path":"/dev/sda","iommu":false,"pci_segment":2,"id":"id"}`, string(b))

			err = json.NewEncoder(w).Encode(&VMAddDeviceResponse{})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := NewClient(svr.Client(), svr.URL)

		resp, err := c.VM().AddDevice(ctx, &VMAddDeviceRequest{
			Path:       "/dev/sda",
			PCISegment: 2,
			ID:         "id",
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
