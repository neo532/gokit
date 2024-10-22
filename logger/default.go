package logger

import (
	"context"
	"fmt"
)

type DefaultLogger struct {
	log        Executor
	globalArgs []interface{}
	level      Level
}

func NewDefaultLogger(l Executor) Logger {
	return &DefaultLogger{
		log:        l,
		globalArgs: []interface{}{},
		level:      l.Level(),
	}
}

func (l *DefaultLogger) Close() (err error) {
	return l.log.Close()
}

func (l *DefaultLogger) WithArgs(kvs ...interface{}) (n Logger) {
	return &DefaultLogger{
		log:        l.log,
		globalArgs: kvs,
		level:      l.level,
	}
}

// WithLevel is only effective when the input parameter is higher than the original.
func (l *DefaultLogger) WithLevel(lv Level) (n Logger) {
	return &DefaultLogger{
		log:        l.log,
		globalArgs: l.globalArgs,
		level:      lv,
	}
}

func (l *DefaultLogger) disabled(lv Level) (b bool) {
	return lv < l.level
}

func (l *DefaultLogger) Debugf(c context.Context, format string, kvs ...interface{}) {
	if l.disabled(LevelDebug) {
		return
	}
	l.log.Log(c, LevelDebug, fmt.Sprintf(format, kvs...), l.globalArgs...)
}
func (l *DefaultLogger) Warnf(c context.Context, format string, kvs ...interface{}) {
	if l.disabled(LevelWarn) {
		return
	}
	l.log.Log(c, LevelWarn, fmt.Sprintf(format, kvs...), l.globalArgs...)
}
func (l *DefaultLogger) Infof(c context.Context, format string, kvs ...interface{}) {
	if l.disabled(LevelInfo) {
		return
	}
	l.log.Log(c, LevelInfo, fmt.Sprintf(format, kvs...), l.globalArgs...)
}
func (l *DefaultLogger) Errorf(c context.Context, format string, kvs ...interface{}) {
	if l.disabled(LevelError) {
		return
	}
	l.log.Log(c, LevelError, fmt.Sprintf(format, kvs...), l.globalArgs...)
}
func (l *DefaultLogger) Fatalf(c context.Context, format string, kvs ...interface{}) {
	if l.disabled(LevelFatal) {
		return
	}
	l.log.Log(c, LevelFatal, fmt.Sprintf(format, kvs...), l.globalArgs...)
}

func (l *DefaultLogger) Debug(c context.Context, message string, kvs ...interface{}) {
	if l.disabled(LevelDebug) {
		return
	}
	l.log.Log(c, LevelDebug, message, append(l.globalArgs, kvs...)...)
}
func (l *DefaultLogger) Warn(c context.Context, message string, kvs ...interface{}) {
	if l.disabled(LevelWarn) {
		return
	}
	l.log.Log(c, LevelWarn, message, append(l.globalArgs, kvs...)...)
}
func (l *DefaultLogger) Info(c context.Context, message string, kvs ...interface{}) {
	if l.disabled(LevelInfo) {
		return
	}
	l.log.Log(c, LevelInfo, message, append(l.globalArgs, kvs...)...)
}
func (l *DefaultLogger) Error(c context.Context, message string, kvs ...interface{}) {
	if l.disabled(LevelError) {
		return
	}
	l.log.Log(c, LevelError, message, append(l.globalArgs, kvs...)...)
}
func (l *DefaultLogger) Fatal(c context.Context, message string, kvs ...interface{}) {
	if l.disabled(LevelFatal) {
		return
	}
	l.log.Log(c, LevelFatal, message, append(l.globalArgs, kvs...)...)
}
