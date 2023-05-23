package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMPowerButton(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := NewClient(http.DefaultClient, "localhost:badport")

		err := c.VM().PowerButton(ctx)
		assert.EqualError(t, err, `failed to call power-button: do request: Put "localhost:badport/api/v1/vm.power-button": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.power-button", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := NewClient(svr.Client(), svr.URL)

		err := c.VM().PowerButton(ctx)
		assert.NoError(t, err)
	})
}
