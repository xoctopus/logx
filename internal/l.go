package internal

import "context"

type Logger interface {
	WithGroup(string) Logger
	With(...any) Logger
	LogIfEnabled(ctx context.Context, lv LogLevel, msg string)
}

const (
	BadKey       = "!BADKEY"
	MissingValue = "!MISSINGVALUE"
)
