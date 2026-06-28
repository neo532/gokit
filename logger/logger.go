package logger

import (
	"context"
)

type Executor interface {
	Log(c context.Context, level Level, message string, kvs ...any) error
	Close() error
	Level() Level
}

type ContextArgs func(c context.Context) (k string, v any)

type Logger interface {
	WithArgs(kvs ...any) (n Logger)
	WithLevel(lv Level) (n Logger)
	Close() error

	Debugf(c context.Context, format string, kvs ...any)
	Warnf(c context.Context, format string, kvs ...any)
	Infof(c context.Context, format string, kvs ...any)
	Errorf(c context.Context, format string, kvs ...any)
	Fatalf(c context.Context, format string, kvs ...any)

	Debug(c context.Context, message string, kvs ...any)
	Warn(c context.Context, message string, kvs ...any)
	Info(c context.Context, message string, kvs ...any)
	Error(c context.Context, message string, kvs ...any)
	Fatal(c context.Context, message string, kvs ...any)
}
