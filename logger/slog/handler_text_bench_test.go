package slog

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/neo532/gokit/logger"
)

type nopWriter struct{}

func (nopWriter) Close() error      { return nil }
func (nopWriter) Writer() io.Writer { return io.Discard }

func BenchmarkCustomSeparator(b *testing.B) {
	h := New(
		WithLevel("info"),
		WithTextSeparator("||"),
		WithWriter(nopWriter{}),
	)
	if err := h.Error(); err != nil {
		b.Fatal(err)
	}
	l := logger.NewDefaultLogger(h)
	c := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithArgs("module", "m1").Info(c, "msg1", "err", "panic")
	}
}

func BenchmarkStdTextHandler(b *testing.B) {
	h := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo.Level()})
	l := slog.New(h)
	c := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.LogAttrs(c, slog.LevelInfo, "msg1", slog.String("err", "panic"), slog.String("module", "m1"))
	}
}

func BenchmarkDefaultHandler(b *testing.B) {
	// gokit slog with default handler (space separator) for comparison
	h := New(
		WithLevel("info"),
		WithWriter(nopWriter{}),
	)
	if err := h.Error(); err != nil {
		b.Fatal(err)
	}
	l := logger.NewDefaultLogger(h)
	c := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithArgs("module", "m1").Info(c, "msg1", "err", "panic")
	}
}
