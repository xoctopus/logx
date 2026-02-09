package handlers

import (
	"io"
	"log/slog"
	"os"
)

func Std(ws ...io.Writer) slog.Handler {
	var w io.Writer = os.Stderr
	if len(ws) > 0 {
		w = ws[0]
	}

	h := &handler{
		skip: 5,
	}
	if gLogFormat == LogFormatTEXT {
		h.h = slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       gLogLevel,
			ReplaceAttr: Replacer,
		})
	} else {
		h.h = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       gLogLevel,
			ReplaceAttr: Replacer,
		})
	}
	return h
}

func DiscardStd() slog.Handler {
	return Std(io.Discard)
}
