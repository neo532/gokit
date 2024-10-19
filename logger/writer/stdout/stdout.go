package stdout

import (
	"io"
	"os"

	"github.com/neo532/gokit/logger/writer"
)

var _ writer.Writer = (*Writer)(nil)

type Writer struct {
	writer *os.File
}

func New() (w *Writer) {

	w = &Writer{
		writer: os.Stdout,
	}

	return
}

func (w *Writer) Close() (err error) {
	return
}

func (w *Writer) Writer() io.Writer {
	return w.writer
}
