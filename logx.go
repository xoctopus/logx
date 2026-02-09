package logx

import (
	"context"
	"runtime"
	"strings"
)

type Logger interface {
	Start(ctx context.Context, name string, kvs ...any) (context.Context, Logger)
	End()

	With(kvs ...any) Logger

	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(err error)
	Error(err error)
}

type Printer interface {
	Print(string)
	Printf(string, ...any)
}

func Start(ctx context.Context, name string, kvs ...any) (context.Context, Logger) {
	return From(ctx).Start(ctx, name, kvs...)
}

func Enter(ctx context.Context, kvs ...any) (context.Context, Logger) {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(name, "/"); idx != -1 {
		name = name[idx+1:]
	}
	return From(ctx).Start(ctx, name, kvs...)
}
