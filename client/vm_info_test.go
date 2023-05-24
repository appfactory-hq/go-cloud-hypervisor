package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientVMInfo(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		ctx := context.Background()

		c := New(WithHTTPEndpoint("localhost:badport"))

		_, err := c.VM().Info(ctx)
		assert.EqualError(t, err, `failed to call info: do request: Get "localhost:badport/api/v1/vm.info": unsupported protocol scheme "localhost"`)
	})

	t.Run("success", func(t *testing.T) {
		response := &VMInfoResponse{
			State: "Running",
		}

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/vm.info", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)

			err := json.NewEncoder(w).Encode(response)
			assert.NoError(t, err)
		}))
		defer svr.Close()

		ctx := context.Background()

		c := New(WithHTTPClient(svr.Client()), WithHTTPEndpoint(svr.URL))

		resp, err := c.VM().Info(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Running", resp.State)
	})
}
