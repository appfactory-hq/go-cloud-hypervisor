package cloudhypervisor

import (
	"fmt"
	"net/http"
)

var (
	ErrNoEndpoint = fmt.Errorf("endpoint is required")
)

type CloudHypervisor struct {
	endpoint string
	client   *http.Client
}

func New(socketPath string, opts ...Option) (*CloudHypervisor, error) {
	ch := &CloudHypervisor{
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(ch)
	}

	if ch.endpoint == "" {
		return nil, ErrNoEndpoint
	}

	ch.client.Transport = newUnixSocketTransport(socketPath)

	return ch, nil
}
