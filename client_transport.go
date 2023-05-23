package cloudhypervisor

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

			return net.DialUnix("unix", nil, addr)
		},
	}
}
