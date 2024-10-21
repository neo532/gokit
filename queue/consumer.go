package queue

/*
 * @abstract queue define
 * @mail neo532@126.com
 * @date 2023-08-13
 */

import (
	"context"
)

type Consumer interface {
	Start(context.Context) error
	Stop(context.Context) error
	Name() string
}
