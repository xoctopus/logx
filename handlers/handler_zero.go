package handlers

import (
	"io"
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	slog0 "github.com/samber/slog-zerolog/v2"
)

func Zero(ws ...io.Writer) slog.Handler {
	w := io.Writer(os.Stderr)
	if len(ws) > 0 && ws[0] != nil {
		w = ws[0]
	}

	logger := zerolog.New(w)
	opt := slog0.Option{
		Level:     gLogLevel,
		Converter: Converter,
		Logger:    &logger,
	}
	slog0.SourceKey = KEY_SOURCE
	return &handler{
		h:    opt.NewZerologHandler(),
		skip: 5,
	}
}

func DiscardZero() slog.Handler {
	return Zero(io.Discard)
}
