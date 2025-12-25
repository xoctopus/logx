package logx

import (
	"context"

	"github.com/xoctopus/x/contextx"
)

type k struct{}

func With(ctx context.Context, l Logger) context.Context {
	return contextx.WithValue(ctx, k{}, l)
}

func From(ctx context.Context) Logger {
	if l, ok := ctx.Value(k{}).(Logger); ok {
		return l
	}
	return DefaultStd()
}

func Carry(l Logger) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, l)
	}
}
