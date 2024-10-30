package logger

import (
	"context"
	"fmt"
	"os"
)

var _ ILogger = (*DefaultILogger)(nil)

type ILogger interface {
	Error(c context.Context, message string, kvs ...interface{})
	Warn(c context.Context, message string, kvs ...interface{})
	Info(c context.Context, message string, kvs ...interface{})
	Fatal(c context.Context, message string, kvs ...interface{})
}

type DefaultILogger struct {
}

func NewDefaultILogger() *DefaultILogger {
	return &DefaultILogger{}
}
func (l *DefaultILogger) Error(c context.Context, message string, kvs ...interface{}) {
	fmt.Println(append([]interface{}{"msg", message}, kvs...)...)
}
func (l *DefaultILogger) Warn(c context.Context, message string, kvs ...interface{}) {
	fmt.Println(append([]interface{}{"msg", message}, kvs...)...)
}
func (l *DefaultILogger) Info(c context.Context, message string, kvs ...interface{}) {
	fmt.Println(append([]interface{}{"msg", message}, kvs...)...)
}
func (l *DefaultILogger) Fatal(c context.Context, message string, kvs ...interface{}) {
	fmt.Println(append([]interface{}{"msg", message}, kvs...)...)
	os.Exit(1)
}
