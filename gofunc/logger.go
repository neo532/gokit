package gofunc

/*
 * @abstract guard panic
 * @author liuxiaofeng
 * @mail neo532@126.com
 * @date 2020-10-06
 */

import (
	"context"
	"sync"

	"github.com/neo532/gokit/errorx"
)

type Logger interface {
	Error(c context.Context, err error)
	Err() error
}

type DefaultLogger struct {
	err  error
	lock sync.Mutex
}

func (l *DefaultLogger) Error(c context.Context, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.err != nil {
		l.err = errorx.WrapError(l.err, err)
		return
	}
	l.err = err
}

func (l *DefaultLogger) Err() error {
	return l.err
}
