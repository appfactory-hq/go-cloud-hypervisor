// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cloudhypervisor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

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

func (m *Machine) DescribeInstanceInfo(ctx context.Context) (*client.VMInfoResponse, error) {
	m.logger.Debug("Called Machine.DescribeInstanceInfo()")

	info, err := m.client.VM().Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance info: %w", err)
	}

	return info, nil
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
	m.logger.Debug("Called Machine.GetVersion()")

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

// StopVMM stops the current VMM.
func (m *Machine) StopVMM() error {
	m.logger.Debug("Called Machine.StopVMM()")

	return m.stopVMM()
}

func (m *Machine) stopVMM() error {
	if m.cmd != nil && m.cmd.Process != nil {
		m.logger.Debug("stopVMM(): sending sigterm to cloud-hypervisor process")

		if err := m.cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fmt.Errorf("sending sigterm to process: %w", err)
		}

		return nil
	}

	m.logger.Debug("stopVMM(): no cloud-hypervisor process running, not sending a signal")

	return nil
}

// nolint: unused
// Set up a signal handler to pass through to cloud-hypervisor.
func (m *Machine) setupSignals() {
	// signals := m.Cfg.ForwardSignals
	signals := []os.Signal{}

	if len(signals) == 0 {
		return
	}

	m.logger.Debug("Setting up signal handler")

	sigchan := make(chan os.Signal, len(signals))

	signal.Notify(sigchan, signals...)

	go func() {
	ForLoop:
		for {
			select {
			case sig := <-sigchan:
				m.logger.Debug("Caught signal", "signal", sig)

				// Some signals kill the process, some of them are not.
				if err := m.cmd.Process.Signal(sig); err != nil {
					m.logger.Error("Failed to send signal to process", "err", err)
				}
			case <-m.exitCh:
				// And if a signal kills the process, we can stop this for loop and remove sigchan.
				break ForLoop
			}
		}

		signal.Stop(sigchan)

		close(sigchan)
	}()
}
