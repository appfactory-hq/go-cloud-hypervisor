package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMResume(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		err := c.VM().Resume(ctx)
		assert.EqualError(t, err, `failed to call resume: do request: Put "localhost:badport/api/v1/vm.resume": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.resume", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		err := c.VM().Resume(ctx)
		assert.NoError(t, err)
	})
}
