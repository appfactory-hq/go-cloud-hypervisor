package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMCreate(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		err := c.VM().Create(ctx, &VMCreateRequest{})
		assert.EqualError(t, err, `failed to call create: do request: Put "localhost:badport/api/v1/vm.create": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.create", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"payload":{"firmware":"firmware","kernel":"kernel"},"iommu":false,"watchdog":false}`, string(b))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		err := c.VM().Create(ctx, &VMCreateRequest{
			VMConfig: &VMConfig{
				Payload: &VMConfigPayload{
					Kernel:   "kernel",
					Firmware: "firmware",
				},
			},
		})
		assert.NoError(t, err)
	})
}
