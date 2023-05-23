package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMResize(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := NewClient(http.DefaultClient, "localhost:badport")

		err := c.VM().Resize(ctx, &VMResizeRequest{})
		assert.EqualError(t, err, `failed to call resize: do request: Put "localhost:badport/api/v1/vm.resize": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.resize", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"desired_vcpus":1,"desired_ram":1}`, string(b))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := NewClient(svr.Client(), svr.URL)

		err := c.VM().Resize(ctx, &VMResizeRequest{
			DesiredVCPUs: 1,
			DesiredRAM:   1,
		})
		assert.NoError(t, err)
	})
}
