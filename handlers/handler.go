package handlers

import (
	"context"
	"log/slog"
	"runtime"
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
