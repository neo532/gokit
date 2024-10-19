package writer

import "io"

type Writer interface {
	Close() (err error)
	Writer() (w io.Writer)
}
