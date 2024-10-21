package slog

import (
	"os"

	"log/slog"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/logger/writer"
)

type Option func(opt *Logger)

func WithLogger(log *slog.Logger) Option {
	return func(l *Logger) {
		l.logger = log
		return
	}
}

func WithPrettyLogger(handler slog.Handler) Option {
	return func(l *Logger) {
		if handler == nil {
			l.logger = slog.New(
				NewPrettyHandler(os.Stdout, l.opts, l.paramContext),
			).With(l.paramGlobal...)
			return
		}
		l.logger = slog.New(handler).With(l.paramGlobal...)
		return
	}
}

func WithReplaceAttr(fns ...func() (k string, v interface{})) Option {
	return func(l *Logger) {
		l.opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			for _, fn := range fns {
				k, v := fn()
				if k == a.Key {
					if v == nil {
						a.Key = k
						break
					}
					a = slog.Any(k, v)
					break
				}
			}
			return a
		}
	}
}

func WithContextParam(fns ...logger.ContextArgs) Option {
	return func(l *Logger) {
		l.paramContext = fns
	}
}

func WithGlobalParam(vs ...interface{}) Option {
	return func(l *Logger) {
		l.paramGlobal = vs
	}
}

func WithLevel(lv string) Option {
	return func(l *Logger) {
		lvl := (&slog.LevelVar{})
		if err := lvl.UnmarshalText([]byte(lv)); err != nil && l.err == nil {
			l.err = err
			return
		}
		l.opts.Level = lvl
	}
}

func WithWriter(w writer.Writer) Option {
	return func(l *Logger) {
		l.writer = w
	}
}
