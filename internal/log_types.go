package internal

import "context"

type Logger interface {
	WithGroup(string) Logger
	With(...any) Logger
	LogIfEnabled(ctx context.Context, lv LogLevel, msg string)
}

type SecurityStringer interface {
	SecurityString() string
}
