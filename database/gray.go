package database

/*
 * @abstract Determine whether it is a shadow database
 * @mail neo532@126.com
 * @date 2024-05-18
 */

import (
	"context"
)

var _ Grayer = (*DefaultGrayer)(nil)

type Grayer interface {
	Judge(c context.Context) (b bool)
}

type DefaultGrayer struct {
}

func (j *DefaultGrayer) Judge(c context.Context) (b bool) {
	return
}
