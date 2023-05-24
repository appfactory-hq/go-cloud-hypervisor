package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMSendMigration(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		err := c.VM().SendMigration(ctx, &VMSendMigrationRequest{})
		assert.EqualError(t, err, `failed to call send-migration: do request: Put "localhost:badport/api/v1/vm.send-migration": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.send-migration", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			b, err := io.ReadAll(r.Body)

			assert.NoError(t, err)
			assert.Equal(t, `{"destination_url":"file:///path/to/state","local":true}`, string(b))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		err := c.VM().SendMigration(ctx, &VMSendMigrationRequest{
			DestinationURL: "file:///path/to/state",
			Local:          true,
		})
		assert.NoError(t, err)
	})
}
