package logger

import (
	"context"
	"fmt"
	"os"
)

var _ ILogger = (*DefaultILogger)(nil)

type ILogger interface {
	Error(c context.Context, message string, kvs ...any)
	Warn(c context.Context, message string, kvs ...any)
	Info(c context.Context, message string, kvs ...any)
	Debug(c context.Context, message string, kvs ...any)
	Fatal(c context.Context, message string, kvs ...any)
}

type DefaultILogger struct {
}

func NewDefaultILogger() *DefaultILogger {
	return &DefaultILogger{}
}
func (l *DefaultILogger) Error(c context.Context, message string, kvs ...any) {
	fmt.Println(append([]any{"msg:", message}, kvs...)...)
}
func (l *DefaultILogger) Warn(c context.Context, message string, kvs ...any) {
	fmt.Println(append([]any{"msg:", message}, kvs...)...)
}
func (l *DefaultILogger) Debug(c context.Context, message string, kvs ...any) {
	fmt.Println(append([]any{"msg:", message}, kvs...)...)
}
func (l *DefaultILogger) Info(c context.Context, message string, kvs ...any) {
	fmt.Println(append([]any{"msg:", message}, kvs...)...)
}
func (l *DefaultILogger) Fatal(c context.Context, message string, kvs ...any) {
	fmt.Println(append([]any{"msg:", message}, kvs...)...)
	os.Exit(1)
}
