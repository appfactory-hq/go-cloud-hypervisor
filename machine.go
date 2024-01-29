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
	"time"

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

func WithVirtualMachineID(id string) MachineOption {
	return func(m *Machine) {
		m.id = id
	}
}

func WithSocketPath(path string) MachineOption {
	return func(m *Machine) {
		m.socketPath = path
	}
}

// forwardSignals
func WithForwardSignals(signals ...os.Signal) MachineOption {
	return func(m *Machine) {
		m.forwardSignals = signals
	}
}

func WithInitTimeout(timeout time.Duration) MachineOption {
	return func(m *Machine) {
		m.initTimeout = timeout
	}
}

type Machine struct {
	id             string
	socketPath     string
	cmd            *exec.Cmd
	startOnce      sync.Once
	cleanupOnce    sync.Once
	exitCh         chan struct{}
	cleanupCh      chan struct{}
	client         *client.Client
	logger         Logger
	forwardSignals []os.Signal
	fatalErr       error
	cleanupFuncs   []func() error
	initTimeout    time.Duration
}

// NewMachine creates a new Machine.
func NewMachine(opts ...MachineOption) (*Machine, error) {
	ctx := context.Background()

	m := &Machine{
		exitCh:    make(chan struct{}),
		cleanupCh: make(chan struct{}),
		logger:    &NoopLogger{},
		forwardSignals: []os.Signal{
			os.Interrupt,
			syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGABRT,
		},
		cleanupFuncs: []func() error{},
		initTimeout:  5 * time.Second,
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.cmd == nil {
		m.cmd = defaultCloudHypervisorVMMCommandBuilder.WithSocketPath(m.socketPath).Build(ctx)
	}

	if m.id == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to create random ID for VMID: %w", err)
		}

		m.id = id.String()
	}

	if m.client == nil {
		m.client = client.New(client.WithUnixSocket(m.socketPath))
	}

	return m, nil
}

// ID returns the machine's ID.
func (m *Machine) ID() string {
	return m.id
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

	var err error
	defer func() {
		if err != nil {
			if cleanupErr := m.doCleanup(); cleanupErr != nil {
				m.logger.ErrorCtx(ctx, "failed to cleanup VM after previous start failure", "err", cleanupErr)
			}
		}
	}()

	err = m.startInstance(ctx)
	if err != nil {
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

// startVMM starts the cloud-hypervisor vmm process and configures logging.
func (m *Machine) startVMM(ctx context.Context) error {
	m.logger.DebugCtx(ctx, "Called startVMM(), setting up a VMM", "socket-path", m.socketPath)

	startCmd := m.cmd.Start

	m.logger.DebugCtx(ctx, fmt.Sprintf("Starting %v", m.cmd.Args))

	// TODO: add CNI support

	if err := startCmd(); err != nil {
		m.logger.ErrorCtx(ctx, "Failed to start VMM", "err", err)

		m.fatalErr = err
		close(m.exitCh)

		return fmt.Errorf("failed to start VMM: %w", err)
	}

	m.logger.DebugCtx(ctx, "VMM started", "socket-path", m.socketPath)

	m.cleanupFuncs = append(m.cleanupFuncs,
		func() error {
			if err := os.Remove(m.socketPath); !os.IsNotExist(err) {
				return err
			}
			return nil
		},
	)

	errCh := make(chan error)
	go func() {
		var merr error

		if err := m.cmd.Wait(); err != nil {
			m.logger.WarnCtx(ctx, "cloud-hypervisor exited", "err", err.Error())
			merr = errors.Join(merr, err)
		} else {
			m.logger.InfoCtx(ctx, "cloud-hypervisor exited", "status", "0")
		}

		if err := m.doCleanup(); err != nil {
			m.logger.ErrorCtx(ctx, "failed to cleanup after VM exit", "err", err)

			merr = errors.Join(merr, err)
		}

		errCh <- merr

		// Notify subscribers that there will be no more values.
		// When err is nil, two reads are performed (waitForSocket and close exitCh goroutine),
		// second one never ends as it tries to read from empty channel.
		close(errCh)
		close(m.cleanupCh)
	}()

	m.setupSignals()

	// Wait for cloud-hypervisor to initialize:
	if err := m.waitForSocket(m.initTimeout, errCh); err != nil {
		err = fmt.Errorf("cloud-hypervisor did not create API socket %s: %w", m.socketPath, err)

		m.fatalErr = err

		close(m.exitCh)

		return err
	}

	// This goroutine is used to kill the process by context cancellation,
	// but doesn't tell anyone about that.
	go func() {
		select {
		case <-ctx.Done():
			break
		case <-m.exitCh:
			// VMM exited on its own; no need to stop it.
			return
		}

		if err := m.stopVMM(); err != nil {
			m.logger.ErrorCtx(ctx, "failed to stop vm", "err", err)
		}
	}()

	// This goroutine is used to tell clients that the process is stopped
	// (gracefully or not).
	go func() {
		m.fatalErr = <-errCh

		m.logger.DebugCtx(ctx, "closing the exitCh", "err", m.fatalErr)

		close(m.exitCh)
	}()

	m.logger.DebugCtx(ctx, "returning from startVMM()")

	return nil
}

func (m *Machine) doCleanup() error {
	var err error

	m.cleanupOnce.Do(func() {
		// run them in reverse order so changes are "unwound" (similar to defer statements)
		for i := range m.cleanupFuncs {
			cleanupFunc := m.cleanupFuncs[len(m.cleanupFuncs)-1-i]

			err = errors.Join(err, cleanupFunc())
		}
	})

	return err
}

// waitForSocket waits for the given file to exist
func (m *Machine) waitForSocket(timeout time.Duration, exitchan chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	ticker := time.NewTicker(10 * time.Millisecond)

	defer func() {
		cancel()

		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-exitchan:
			return err
		case <-ticker.C:
			if _, err := os.Stat(m.socketPath); err != nil {
				continue
			}

			// Send test HTTP request to make sure socket is available
			if _, err := m.client.VMM().Ping(ctx); err != nil {
				continue
			}

			return nil
		}
	}
}
