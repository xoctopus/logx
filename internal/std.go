package internal

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime"
	"strings"

	"github.com/xoctopus/x/reflectx"
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

	v := a.Value.Any()
	if x, ok := v.(SecurityStringer); ok {
		a.Value = slog.StringValue(x.SecurityString())
	} else {
		_, ok = sensitives[a.Key]
		if ok && reflectx.KindOf(v) == reflect.String {
			a.Key = a.Key + "*"
			a.Value = slog.StringValue(MASKED)
		}
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
		a.Value = slog.StringValue(fmt.Sprintf("%s:%03d", loc, s.Line))
	}
	return a
}

func _newstd(w io.Writer, skip int, level slog.Level) *slog.Logger {
	h := &handler{
		skip: skip,
	}
	if gLogFormat == LogFormatTEXT {
		h.h = slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: replacer,
		})
	} else {
		h.h = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
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

func StdLogger(w io.Writer, skip int, lv LogLevel) Logger {
	return &_std{l: _newstd(w, skip, lv)}
}

func StdDiscardLogger(skip int, lv LogLevel) Logger {
	return StdLogger(io.Discard, skip, lv)
}

func NawStdLogger(w io.Writer, skip int, lv slog.Level) *slog.Logger {
	return _newstd(w, skip, lv)
}
