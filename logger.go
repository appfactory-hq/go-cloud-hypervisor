// Copyright 2023 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cloudhypervisor

import "context"

type Logger interface {
	Debug(msg string, args ...any)
	DebugCtx(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoCtx(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnCtx(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorCtx(ctx context.Context, msg string, args ...any)
}

type NoopLogger struct{}

func (l NoopLogger) Debug(msg string, args ...any) {}

func (l NoopLogger) DebugCtx(ctx context.Context, msg string, args ...any) {}

func (l NoopLogger) Info(msg string, args ...any) {}

func (l NoopLogger) InfoCtx(ctx context.Context, msg string, args ...any) {}

func (l NoopLogger) Warn(msg string, args ...any) {}

func (l NoopLogger) WarnCtx(ctx context.Context, msg string, args ...any) {}

func (l NoopLogger) Error(msg string, args ...any) {}

func (l NoopLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {}
