package cloudhypervisor

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMCommandBuilderImmutaility(t *testing.T) {
	b := VMCommandBuilder{}
	b.WithSocketPath("foo").
		WithArgs([]string{"baz", "qux"}).
		AddArgs("moo", "cow")

	assert.Equal(t, []string(nil), b.SocketPath())
	assert.Equal(t, []string(nil), b.Args())
	assert.Equal(t, "cloud-hypervisor", b.Bin())
}

func TestVMCommandBuilderChaining(t *testing.T) {
	b := VMCommandBuilder{}.
		WithSocketPath("socket-path").
		WithBin("bin")

	assert.Equal(t, []string{"--api-socket", "socket-path"}, b.SocketPath())
	assert.Equal(t, "bin", b.Bin())
}

func TestVMCommandBuilderBuild(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stdin := &bytes.Buffer{}

	ctx := context.Background()
	b := VMCommandBuilder{}.
		WithSocketPath("socket-path").
		WithBin("bin").
		WithStdout(stdout).
		WithStderr(stderr).
		WithStdin(stdin).
		WithArgs([]string{"foo"}).
		AddArgs("--bar", "baz")
	cmd := b.Build(ctx)

	expectedArgs := []string{
		"bin",
		"--api-socket",
		"socket-path",
		"foo",
		"--bar",
		"baz",
	}

	assert.Same(t, stdout, cmd.Stdout)
	assert.Same(t, stderr, cmd.Stderr)
	assert.Same(t, stdin, cmd.Stdin)
	assert.Equal(t, expectedArgs, cmd.Args)
}
