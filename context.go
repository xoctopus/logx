package logx

import (
	"context"

	"github.com/xoctopus/x/contextx"

	"github.com/xoctopus/logx/handlers"
)

type tLogContext struct{}

func With(ctx context.Context, l Logger) context.Context {
	return contextx.WithValue(ctx, tLogContext{}, l)
}

func From(ctx context.Context) Logger {
	if l, ok := ctx.Value(tLogContext{}).(Logger); ok {
		return l
	}
	return New(handlers.Std())
}

func Carry(l Logger) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, l)
	}
}
