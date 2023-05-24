package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMMPing(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VMM().Ping(ctx)
		assert.EqualError(t, err, `failed to call ping: do request: Get "localhost:badport/api/v1/vmm.ping": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vmm.ping", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			err := json.NewEncoder(w).Encode(&VMMPingResponse{
				Version: "1.0.0",
				PID:     1234,
			})
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VMM().Ping(ctx)
		assert.NoError(t, err)

		assert.Equal(t, "1.0.0", resp.Version)
		assert.Equal(t, 1234, resp.PID)
	})
}
