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

func TestClientVMAddVDPA(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VM().AddVDPA(ctx, &VMAddVDPARequest{})
		assert.EqualError(t, err, `failed to call add-vdpa: do request: Put "localhost:badport/api/v1/vm.add-vdpa": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.add-vdpa", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"path":"/path/to/vdpa","num_queues":3,"iommu":true,"pci_segment":2,"id":"foo"}`, string(b))

			err = json.NewEncoder(w).Encode(&VMAddVDPAResponse{
				VMPCIDeviceInfo: &VMPCIDeviceInfo{
					ID: "foo",
				},
			})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VM().AddVDPA(ctx, &VMAddVDPARequest{
			VMConfigVDPA: &VMConfigVDPA{
				NumQueues:  3,
				Path:       "/path/to/vdpa",
				IOMMU:      true,
				ID:         "foo",
				PCISegment: 2,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "foo", resp.ID)
	})
}
