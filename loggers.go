package logx

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/xoctopus/logx/internal"
)

// NewStd creates a new logger with the slog.Logger.
func NewStd() Logger {
	return &std{
		ctx: context.Background(),
		l:   internal.StdLogger(os.Stderr, 6, internal.GetLogLevel()),
	}
}

// NewZap creates a new logger with the zap.Logger.
func NewZap() Logger {
	return &std{
		ctx: context.Background(),
		l:   internal.ZapLogger(os.Stderr, 2, internal.GetLogLevel()),
	}
}

// NewWithInstance creates a new logger with the LoggerInstance
func NewWithInstance(l LoggerInstance) Logger {
	return &std{
		ctx: context.Background(),
		l:   l,
	}
}

var NewDefault = NewStd

type std struct {
	ctx   context.Context
	l     internal.Logger
	spans []string
}

func (s *std) Start(ctx context.Context, name string, kvs ...any) (context.Context, Logger) {
	spans := append(s.spans, name)
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

func (s *std) WithAttr(attrs ...slog.Attr) Logger {
	kvs := make([]any, len(attrs)*2)
	for i := range attrs {
		kvs = append(kvs, attrs[i].Key, attrs[i].Value.Any())
	}
	return &std{
		ctx:   s.ctx,
		spans: s.spans,
		l:     s.l.With(kvs...),
	}
}

func (s *std) Debug(msg string, args ...any) {
	s.l.LogIfEnabled(s.ctx, LogLevelDebug, fmt.Sprintf(msg, args...))
}

func (s *std) Info(msg string, args ...any) {
	s.l.LogIfEnabled(s.ctx, LogLevelInfo, fmt.Sprintf(msg, args...))
}

func (s *std) Warn(err error) {
	s.l.LogIfEnabled(s.ctx, LogLevelWarn, err.Error())
}

func (s *std) Error(err error) {
	s.l.LogIfEnabled(s.ctx, LogLevelError, err.Error())
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

func (d discard) WithAttr(attrs ...slog.Attr) Logger {
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
