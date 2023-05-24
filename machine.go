package cloudhypervisor

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/appfactory-hq/go-cloud-hypervisor/client"
	"github.com/google/uuid"
)

var (
	ErrMachineNotRunning     = fmt.Errorf("machine is not running")
	ErrMachineExited         = fmt.Errorf("machine process has exited")
	ErrMachineAlreadyStarted = errors.New("machine already started")
)

type MachineOption func(*Machine)

func WithLogger(logger Logger) MachineOption {
	return func(m *Machine) {
		m.logger = logger
	}
}

type Config struct {
	VMID       string
	SocketPath string
}

type Machine struct {
	cmd       *exec.Cmd
	startOnce sync.Once
	exitCh    chan struct{}
	cleanupCh chan struct{}
	client    *client.Client
	logger    Logger
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
		m.client = client.New(client.WithUnixSocket(cfg.SocketPath))
	}

	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

// PID returns the machine's running process PID or an error if not running.
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

func (m *Machine) Start(ctx context.Context) error {
	m.logger.Debug("Called Machine.Start()")

	alreadyStarted := true

	m.startOnce.Do(func() {
		m.logger.DebugCtx(ctx, "Marking Machine as Started")

		alreadyStarted = false
	})

	if alreadyStarted {
		return ErrMachineAlreadyStarted
	}
	/*
		var err error
		defer func() {
			if err != nil {
				if cleanupErr := m.doCleanup(); cleanupErr != nil {
					m.logger.ErrorCtx(ctx, "failed to cleanup VM after previous start failure", "err", cleanupErr)
				}
			}
		}()
	*/
	if err := m.startInstance(ctx); err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	return nil
}

func (m *Machine) Shutdown(ctx context.Context) error {
	m.logger.Debug("Called Machine.Shutdown()")
	/*
		if runtime.GOARCH != "arm64" {
			return m.sendCtrlAltDel(ctx)
		} else {

			return m.StopVMM()
		}
	*/

	return nil
}

// GetVersion gets the machine's cloud-hypervisor version and returns it.
func (m *Machine) GetVersion(ctx context.Context) (string, error) {
	resp, err := m.client.VMM().Ping(ctx)
	if err != nil {
		m.logger.Error("Getting cloud-hypervisor version", "err", err)

		return "", fmt.Errorf("get vmm version: %w", err)
	}

	m.logger.Debug("GetVersion successful")

	return resp.Version, nil
}

func (m *Machine) startInstance(ctx context.Context) error {
	if err := m.client.VM().Boot(ctx); err != nil {
		m.logger.Error("Starting instance failed", "err", err)

		return fmt.Errorf("boot vm: %w", err)
	}

	m.logger.Info("Starting instance successful")

	return nil
}
