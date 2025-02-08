package logx

import (
	"context"
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
	return FromContext(ctx).Start(ctx, name, kvs...)
}
