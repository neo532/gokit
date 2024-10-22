package zap

/*
 * @abstract zap's option
 * @mail neo532@126.com
 * @date 2023-08-13
 */
import (
	"fmt"
	"io"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/logger/writer"
)

type Option func(opt *Logger)

func WithLogger(log *zap.Logger) Option {
	return func(l *Logger) {
		// free old
		if l.logger != nil {
			l.Close()
		}

		l.logger = log
		return
	}
}

// WithPrettyLogger should be passed as a parameter at the end of the options.
func WithPrettyLogger(w io.Writer) Option {
	return func(l *Logger) {
		if w == nil {
			w = os.Stdout
		}

		// free old
		if l.logger != nil {
			l.Close()
		}

		l.logger = zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(l.core),
				zapcore.AddSync(w),
				l.levelEnabler,
			),
			l.opts...)
		l.Sync = l.logger.Sync
		return
	}
}

func WithCallerSkip(skip int) Option {
	return func(l *Logger) {
		l.opts = append(l.opts, zap.WithCaller(true))
		l.opts = append(l.opts, zap.AddCallerSkip(skip))
	}
}

func WithContextParam(fns ...logger.ContextArgs) Option {
	return func(l *Logger) {
		l.paramContext = fns
	}
}

func WithGlobalParam(kvs ...interface{}) Option {
	return func(l *Logger) {
		ls := len(kvs)
		ps := make([]zap.Field, 0, ls/2)
		for i := 0; i < ls; i += 2 {
			k, _ := kvs[i].(string)
			ps = append(ps, zap.Any(k, kvs[i+1]))
		}
		l.opts = append(l.opts, zap.Fields(ps...))
	}
}

func WithLevel(lv string) Option {
	return func(l *Logger) {

		l.level = logger.ParseLevel(lv)

		fmt.Println(fmt.Sprintf("option:\t%+v", l.level))
		var err error
		if l.levelEnabler, err = zapcore.ParseLevel(l.level.String()); err != nil {
			fmt.Println(runtime.Caller(0))
			l.err = err
			return
		}
	}
}

func WithWriter(w writer.Writer) Option {
	return func(l *Logger) {
		l.writer = w
	}
}
func WithMessageKey(s string) Option {
	return func(l *Logger) {
		l.core.MessageKey = s
	}
}
