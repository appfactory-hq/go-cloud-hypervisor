package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMCoreDump(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		err := c.VM().CoreDump(ctx, &VMCoreDumpRequest{})
		assert.EqualError(t, err, `failed to call coredump: do request: Put "localhost:badport/api/v1/vm.coredump": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.coredump", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"destination_url":"file:///path/to/coredump"}`, string(b))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		err := c.VM().CoreDump(ctx, &VMCoreDumpRequest{
			DestinationURL: "file:///path/to/coredump",
		})
		assert.NoError(t, err)
	})
}
