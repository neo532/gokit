package kafka

/*
 * @abstract kafka's inner logger
 * @mail neo532@126.com
 * @date 2023-08-13
 */

import (
	"context"
	"fmt"

	ilogger "github.com/neo532/gokit/logger"
)

type logger struct {
	log ilogger.ILogger
	ctx context.Context
}

func NewLogger(l ilogger.ILogger) *logger {
	return &logger{
		log: l,
		ctx: context.Background(),
	}
}
func (l *logger) WithContext(c context.Context) *logger {
	l.ctx = c
	return l
}
func (l *logger) Print(v ...any) {
	l.log.Info(l.ctx, fmt.Sprintf("%+v", v))
}
func (l *logger) Printf(format string, v ...any) {
	l.log.Info(l.ctx, fmt.Sprintf(format, v...))
}
func (l *logger) Println(v ...any) {
	l.log.Info(l.ctx, fmt.Sprintf("%+v", v))
}
