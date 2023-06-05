// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

func newUnixSocketTransport(socketPath string) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, path string) (net.Conn, error) {
			addr, err := net.ResolveUnixAddr("unix", socketPath)
			if err != nil {
				return nil, fmt.Errorf("resolve unix addr: %w", err)
			}

			// nolint: wrapcheck // no need to wrap this error
			return net.DialUnix("unix", nil, addr)
		},
	}
}
