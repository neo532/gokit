package log

/*
 * @abstract log's option
 * @mail neo532@126.com
 * @date 2023-08-13
 */
import (
	"errors"
	"log/syslog"

	"github.com/neo532/gokit/logger"
)

type Option func(opt *Logger)

func WithLogger(log *syslog.Writer) Option {
	return func(l *Logger) {
		l.logger = log
	}
}

func WithContextParam(fns ...logger.ContextArgs) Option {
	return func(l *Logger) {
		l.paramContext = fns
	}
}

func WithGlobalParam(kvs ...interface{}) Option {
	return func(l *Logger) {
		l.paramGlobal = kvs
	}
}

var levelMap = map[syslog.Priority]logger.Level{
	syslog.LOG_DEBUG:   logger.LevelDebug,
	syslog.LOG_INFO:    logger.LevelInfo,
	syslog.LOG_WARNING: logger.LevelWarn,
	syslog.LOG_ERR:     logger.LevelError,

	syslog.LOG_NOTICE: logger.LevelInfo,
	syslog.LOG_EMERG:  logger.LevelFatal,
	syslog.LOG_ALERT:  logger.LevelFatal,
	syslog.LOG_CRIT:   logger.LevelFatal,
}

func WithLevel(lv syslog.Priority) Option {
	return func(l *Logger) {
		var ok bool
		if l.level, ok = levelMap[lv]; !ok {
			l.err = errors.New("Wrong level")
		}
	}
}
