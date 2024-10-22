package zap

import (
	"context"
	"fmt"
	"testing"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/logger/writer/lumberjack"
)

func newZap() (h logger.Logger) {
	cp := func(c context.Context) (key string, value interface{}) {
		return "aa", "bbbbbbbbb"
	}
	cp(context.Background())

	l := New(
		WithWriter(
			lumberjack.New(
				lumberjack.WithFilename("./test.log"),
				lumberjack.WithMaxBackups(2),
				lumberjack.WithMaxSize(2),
			),
		),
		WithLevel("info"),
		// WithGlobalParam("a", "b", "1", "2"),
		// WithContextParam(cp),
		WithMessageKey("msg"),
		WithCallerSkip(2),
		WithPrettyLogger(nil),
	)
	if err := l.Error(); err != nil {
		fmt.Println(fmt.Sprintf("err:\t%+v", err))
	}
	return logger.NewDefaultLogger(l)
}
func TestLogger(t *testing.T) {

	c := context.Background()
	h := newZap()
	for i := 0; i < 1; i++ {
		h.WithArgs(logger.KeyModule, "m1").Info(c, "msg1", "err", "panic")
		h.WithArgs(logger.KeyModule, "m2").WithLevel(logger.LevelError).Info(c, "m2", "e1", "p1")
		h.WithArgs(logger.KeyModule, "m3").Errorf(c, "m%s", "3")
	}

	a(c, h)
}

func a(c context.Context, h logger.Logger) {
	h.WithArgs(logger.KeyModule, "m4").Errorf(c, "m%s", "4")
}
