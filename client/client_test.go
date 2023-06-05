// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientEndpoint(t *testing.T) {
	c := &Client{}

	assert.Equal(t, "/api/v1/vmm.ping/", c.endpoint("vmm.ping/"))
}

func TestClientExpectCode(t *testing.T) {
	c := &Client{}

	assert.EqualError(t, c.expectCode(http.StatusCreated, http.StatusAccepted), "unexpected status code: 201")
}

func TestClientCall(t *testing.T) {
	t.Parallel()

	t.Run("With bad request body", func(t *testing.T) {
		ctx := context.Background()

		c := &Client{}

		var ch chan int

		_, err := c.call(ctx, http.MethodPut, "/api/v1/vmm.shutdown/", ch, nil)
		assert.EqualError(t, err, "marshal body: json: unsupported type: chan int")
	})

	t.Run("With bad request method", func(t *testing.T) {
		ctx := context.Background()

		c := &Client{}

		_, err := c.call(ctx, "bad,get", "/api/v1/vmm.shutdown/", nil, nil)
		assert.EqualError(t, err, `create request: net/http: invalid method "bad,get"`)
	})

	t.Run("With bad response body", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bad json"))
		}))
		defer svr.Close()
		ctx := context.Background()

		c := &Client{
			httpClient: svr.Client(),
			base:       svr.URL,
		}

		var resp map[string]interface{}
		_, err := c.call(ctx, http.MethodGet, "/", nil, &resp)
		assert.EqualError(t, err, "decode response: invalid character 'b' looking for beginning of value")
	})
}

func TestClientWithUnixSocket(t *testing.T) {
	sockPath := path.Join(os.TempDir(), "socket.sock")

	defer os.Remove(sockPath)

	serv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(&VMMPingResponse{
				Version: "1.0.0",
				PID:     1234,
			})
			assert.NoError(t, err)
		}),
	}

	unixListener, err := net.Listen("unix", sockPath)
	assert.NoError(t, err)

	go serv.Serve(unixListener)
	defer serv.Close()

	ctx := context.Background()

	c := New(WithUnixSocket(sockPath))

	_, err = c.VMM().Ping(ctx)
	assert.NoError(t, err)
}
