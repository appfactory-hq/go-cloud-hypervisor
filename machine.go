package cloudhypervisor

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrMachineNotRunning = fmt.Errorf("machine is not running")
	ErrMachineExited     = fmt.Errorf("machine process has exited")
)

type MachineOption func(*Machine)

type Config struct {
	VMID string
}

type Machine struct {
	cmd       *exec.Cmd
	startOnce sync.Once
	exitCh    chan struct{}
	cleanupCh chan struct{}
	client    *Client
}

// NewMachine creates a new Machine.
func NewMachine(ctx context.Context, cfg Config, opts ...MachineOption) (*Machine, error) {
	m := &Machine{
		exitCh:    make(chan struct{}),
		cleanupCh: make(chan struct{}),
	}

	if cfg.VMID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to create random ID for VMID: %w", err)
		}

		cfg.VMID = id.String()
	}

	if m.client == nil {
		// m.client = NewClient(cfg.SocketPath, m.logger, false)
	}

	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

// PID returns the machine's running process PID or an error if not running
func (m *Machine) PID() (int, error) {
	if m.cmd == nil || m.cmd.Process == nil {
		return 0, ErrMachineNotRunning
	}

	select {
	case <-m.exitCh:
		return 0, ErrMachineExited
	default:
	}

	return m.cmd.Process.Pid, nil
}
