package database

/*
 * @abstract Determine whether it is a shadow database
 * @mail neo532@126.com
 * @date 2024-05-18
 */

import (
	"context"
)

var _ Benchmarker = (*DefaultBenchmarker)(nil)

type Benchmarker interface {
	Judge(c context.Context) (b bool)
}

type DefaultBenchmarker struct {
}

func (j *DefaultBenchmarker) Judge(c context.Context) (b bool) {
	return
}
