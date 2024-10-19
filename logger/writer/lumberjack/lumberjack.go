package lumberjack

import (
	"io"

	"github.com/neo532/gokit/logger/writer"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ writer.Writer = (*Writer)(nil)

type Writer struct {
	writer *lumberjack.Logger
}

func New(opts ...Option) (w *Writer) {

	w = &Writer{
		writer: &lumberjack.Logger{},
	}
	for _, o := range opts {
		o(w)
	}

	return
}

func (w *Writer) Close() (err error) {
	return w.writer.Close()
}

func (w *Writer) Writer() io.Writer {
	return w.writer
}