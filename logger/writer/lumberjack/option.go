package lumberjack

type Option func(w *Writer)

func WithFilename(s string) Option {
	return func(w *Writer) {
		w.writer.Filename = s
	}
}

func WithMaxSize(i int) Option {
	return func(w *Writer) {
		w.writer.MaxSize = i
	}
}

func WithMaxAge(i int) Option {
	return func(w *Writer) {
		w.writer.MaxAge = i
	}
}

func WithMaxBackups(i int) Option {
	return func(w *Writer) {
		w.writer.MaxBackups = i
	}
}

func WithCompress(b bool) Option {
	return func(w *Writer) {
		w.writer.Compress = b
	}
}
