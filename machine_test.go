package cloudhypervisor

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func checkIfLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping test on non-linux system")
	}
}

func TestMachineID(t *testing.T) {
	t.Parallel()

	t.Run("with default id", func(t *testing.T) {
		m, err := NewMachine()
		assert.NoError(t, err)
		assert.NotEmpty(t, m.ID())

		_, err = uuid.Parse(m.ID())
		assert.NoError(t, err)
	})

	t.Run("with default id", func(t *testing.T) {
		id := uuid.New().String()

		m, err := NewMachine(WithVirtualMachineID(id))
		assert.NoError(t, err)
		assert.Equal(t, id, m.ID())
	})
}

func TestMachinePID(t *testing.T) {
	t.Parallel()

	t.Run("with process not running", func(t *testing.T) {
		m, err := NewMachine()
		assert.NoError(t, err)

		pid, err := m.PID()
		assert.EqualError(t, err, ErrMachineNotRunning.Error())
		assert.Equal(t, 0, pid)
	})
}

func TestMachine(t *testing.T) {
	checkIfLinux(t)

	ctx := context.Background()

	socketPath := path.Join(os.TempDir(), "cloud-hypervisor.sock")

	defer os.Remove(socketPath)

	m, err := NewMachine(WithSocketPath(socketPath))
	assert.NoError(t, err)

	go func() {
		err = m.Start(ctx)
		assert.NoError(t, err)
	}()
	/*
		defer func() {
			err = m.Stop(ctx)
			assert.NoError(t, err)
		}()
	*/

	time.Sleep(5 * time.Second)
}
