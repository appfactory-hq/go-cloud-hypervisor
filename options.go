package cloudhypervisor

import "net/http"

type Option func(*CloudHypervisor)

// WithEndpoint sets the endpoint of the CloudHypervisor instance.
func WithEndpoint(endpoint string) Option {
	return func(c *CloudHypervisor) {
		c.endpoint = endpoint
	}
}

// WithHTTPClient sets the HTTP client of the CloudHypervisor instance.
func WithHTTPClient(client *http.Client) Option {
	return func(c *CloudHypervisor) {
		c.client = client
	}
}

// WithClientTransport sets the HTTP transport of the CloudHypervisor instance's HTTP client.
func WithClientTransport(transport http.RoundTripper) Option {
	return func(c *CloudHypervisor) {
		c.client.Transport = transport
	}
}
