package log

import (
	"context"
	"fmt"
	"testing"

	"github.com/neo532/gokit/logger"
)

func createLog() (h logger.Logger) {
	cp := func(c context.Context) (key string, value interface{}) {
		return "aa", "bbbbbbbbb"
	}
	l := New(
		WithGlobalParam("a", "b", "1", "2"),
		WithContextParam(cp),
	)
	if l.err != nil {
		fmt.Println(fmt.Sprintf("err:\t%+v", l.err))
	}
	return logger.NewDefaultLogger(l)
}
func TestLogger(t *testing.T) {

	c := context.Background()
	h := createLog()
	for i := 0; i < 1; i++ {
		// h.Error(c, "k")
		// time.Sleep(10 * time.Second)
		// return
		h.WithArgs(logger.KeyModule, "m1").WithLevel(logger.LevelDebug).Info(c, "kkkk", "vvvv", "cc")
		h.WithArgs(logger.KeyModule, "m2").Errorf(c, "kkkk%s", "cc")
	}

	a(c, h)
}

func a(c context.Context, h logger.Logger) {
	h.WithArgs(logger.KeyModule, "m3").Errorf(c, "kkkk%s", "cc")
}
