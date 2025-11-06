package logx

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func DefaultHandler() slog.Handler {
	return &handler{Handler: slog.Default().Handler()}
}

type handler struct{ slog.Handler }

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	return h.Handler.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{Handler: h.Handler.WithGroup(name)}
}

func (h *handler) Enabled(_ context.Context, lv slog.Level) bool {
	return lv >= slog.LevelDebug
}

func NewHandler() slog.Handler {
	h := &Handler{
		skip: 5,
	}
	if gLogFormat == LogFormatTEXT {
		h.h = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource:   true,
			Level:       gLogLevel,
			ReplaceAttr: Replacer,
		})
	} else {
		h.h = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource:   true,
			Level:       gLogLevel,
			ReplaceAttr: Replacer,
		})
	}
	return h
}

type Handler struct {
	skip int
	h    slog.Handler
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var pcs [1]uintptr
	runtime.Callers(h.skip, pcs[:])
	r.PC = pcs[0]
	return h.h.Handle(ctx, r)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{skip: h.skip, h: h.h.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{skip: h.skip, h: h.h.WithGroup(name)}
}

func (h *Handler) Enabled(ctx context.Context, lv slog.Level) bool {
	return h.h.Enabled(ctx, lv)
}

type SecurityStringer interface {
	SecurityString() string
}

func Replacer(_ []string, a slog.Attr) slog.Attr {
	if a.Value.Kind() == slog.KindTime {
		a.Value = slog.StringValue(a.Value.Time().Format("20060102-150405.000"))
	}

	if x, ok := a.Value.Any().(SecurityStringer); ok {
		a.Value = slog.StringValue(x.SecurityString())
	}

	switch a.Key {
	case slog.TimeKey:
		a.Key = "@ts"
	case slog.LevelKey:
		a.Key = "@lv"
		a.Value = slog.StringValue(gLevelString[a.Value.Any().(slog.Level)])
	case slog.MessageKey:
		a.Key = "@msg"
	case slog.SourceKey:
		a.Key = "@src"
		s := a.Value.Any().(*slog.Source)
		dir, file := filepath.Split(s.Function)
		pkg := dir + file[0:strings.Index(file, ".")]
		a.Value = slog.StringValue(pkg + "/" + filepath.Base(s.File) + ":" + strconv.Itoa(s.Line))
	case "password":
		a.Value = slog.StringValue("****")
	}
	return a
}
