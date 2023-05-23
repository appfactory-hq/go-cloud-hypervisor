package client

import (
	"net/http"
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
