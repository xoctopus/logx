package logx

import (
	"context"

	"github.com/xoctopus/x/contextx"
)

type k struct{}

func WithLogger(ctx context.Context, l Logger) context.Context {
	return contextx.WithValue(ctx, k{}, l)
}

func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(k{}).(Logger); ok {
		return l
	}
	return DefaultStd()
}
