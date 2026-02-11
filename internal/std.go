package internal

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type handler struct {
	skip int
	h    slog.Handler
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	var pcs [1]uintptr
	runtime.Callers(h.skip, pcs[:])
	r.PC = pcs[0]
	return h.h.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{skip: h.skip, h: h.h.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{skip: h.skip, h: h.h.WithGroup(name)}
}

func (h *handler) Enabled(ctx context.Context, lv slog.Level) bool {
	return h.h.Enabled(ctx, lv)
}

func replacer(_ []string, a slog.Attr) slog.Attr {
	if a.Value.Kind() == slog.KindTime {
		a.Value = slog.StringValue(a.Value.Time().Format(TIME_FORMAT))
	}

	x, ok := a.Value.Any().(interface{ SecurityString() string })
	if ok {
		a.Value = slog.StringValue(x.SecurityString())
	}

	switch a.Key {
	case slog.TimeKey:
		a.Key = KEY_TIMESTAMP
	case slog.LevelKey:
		a.Key = KEY_LEVEL
		a.Value = slog.StringValue(gLevelString[a.Value.Any().(slog.Level)])
	case slog.MessageKey:
		a.Key = KEY_MESSAGE
	case slog.SourceKey:
		a.Key = KEY_SOURCE
		s := a.Value.Any().(*slog.Source)
		parts := strings.Split(s.File, "/")
		if l := len(parts); l >= 2 {
			parts = parts[l-2 : l]
		}
		loc := strings.Join(parts, "/")
		a.Value = slog.StringValue(loc + ":" + strconv.Itoa(s.Line))
	case "password":
		if !ok {
			a.Value = slog.StringValue("--------")
		}
	}
	return a
}

func _newstd(w io.Writer, skip int) *slog.Logger {
	h := &handler{
		skip: skip,
	}
	if gLogFormat == LogFormatTEXT {
		h.h = slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       gLogLevel,
			ReplaceAttr: replacer,
		})
	} else {
		h.h = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       gLogLevel,
			ReplaceAttr: replacer,
		})
	}
	return slog.New(h)
}

type _std struct {
	l *slog.Logger
}

func (s *_std) With(kvs ...any) Logger {
	l := s.l.With(kvs...)
	return &_std{l: l}
}

func (s *_std) WithGroup(name string) Logger {
	l := s.l.WithGroup(name)
	return &_std{l: l}
}

func (s *_std) LogIfEnabled(ctx context.Context, lv LogLevel, msg string) {
	if s.l.Enabled(ctx, lv) {
		s.l.Log(ctx, lv, msg)
	}
}

func StdLogger(skip int) Logger {
	return &_std{l: _newstd(os.Stderr, skip)}
}

func StdDiscardLogger(skip int) Logger {
	return &_std{l: _newstd(io.Discard, skip)}
}
