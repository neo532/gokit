package queue

/*
 * @abstract queue's Producer define
 * @mail neo532@126.com
 * @date 2024-10-21
 */

import (
	"context"
)

type Producer interface {
	Err() error
	Close() func()
	Send(c context.Context, message interface{}) (err error)
}

type Consumer interface {
	Start(context.Context) error
	Stop(context.Context) error
	Name() string
}
