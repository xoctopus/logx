package logx

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func DefaultStd() Logger {
	return &std{
		ctx: context.Background(),
		l:   slog.New(DefaultHandler()),
	}
}

func Std(h slog.Handler) Logger {
	return &std{
		ctx: context.Background(),
		l:   slog.New(h),
	}
}

type std struct {
	ctx   context.Context
	l     *slog.Logger
	spans []string
}

func (s *std) Start(ctx context.Context, name string, kvs ...any) (context.Context, Logger) {
	spans := append(s.spans, name)
	if len(kvs) == 0 {
		return ctx, &std{
			ctx:   ctx,
			l:     s.l.WithGroup(strings.Join(spans, "/")),
			spans: spans,
		}
	}
	return ctx, &std{
		ctx:   ctx,
		spans: spans,
		l:     s.l.WithGroup(strings.Join(spans, "/")).With(kvs...),
	}
}

func (s *std) End() {
	if len(s.spans) != 0 {
		s.spans = s.spans[0 : len(s.spans)-1]
	}
}

func (s *std) With(kvs ...any) Logger {
	return &std{
		ctx:   s.ctx,
		spans: s.spans,
		l:     s.l.With(kvs...),
	}
}

func (s *std) Debug(msg string, args ...any) {
	if s.l.Enabled(s.ctx, slog.LevelDebug) {
		s.l.Log(s.ctx, slog.LevelDebug, fmt.Sprintf(msg, args...))
	}
}

func (s *std) Info(msg string, args ...any) {
	if s.l.Enabled(s.ctx, slog.LevelInfo) {
		s.l.Log(s.ctx, slog.LevelInfo, fmt.Sprintf(msg, args...))
	}
}

func (s *std) Warn(err error) {
	if s.l.Enabled(s.ctx, slog.LevelWarn) && err != nil {
		s.l.Log(s.ctx, slog.LevelWarn, err.Error())
	}
}

func (s *std) Error(err error) {
	if err != nil {
		s.l.Log(s.ctx, slog.LevelError, err.Error())
	}
}

func Discard() Logger {
	return discard{}
}

type discard struct{}

func (d discard) Start(ctx context.Context, _ string, _ ...any) (context.Context, Logger) {
	return ctx, d
}

func (d discard) End() {
}

func (d discard) With(kvs ...any) Logger {
	return d
}

func (d discard) Debug(msg string, args ...any) {
}

func (d discard) Info(msg string, args ...any) {
}

func (d discard) Warn(err error) {
}

func (d discard) Error(err error) {
}
