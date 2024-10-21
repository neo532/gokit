package kafka

/*
 * @abstract kafka的具体参数定义
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
func (l *logger) Print(v ...interface{}) {
	l.log.Info(l.ctx, fmt.Sprintf("%+v", v))
}
func (l *logger) Printf(format string, v ...interface{}) {
	l.log.Info(l.ctx, fmt.Sprintf(format, v...))
}
func (l *logger) Println(v ...interface{}) {
	l.log.Info(l.ctx, fmt.Sprintf("%+v", v))
}
