package producer

import (
	"context"
	"testing"
	"time"

	"github.com/neo532/gokit/queue"
)

func TraceID() queue.Middleware {
	return func(handler queue.Handler) queue.Handler {
		return func(c context.Context, message interface{}) (err error) {

			c = queue.AppendHeaderToContext(c, "traceID", "aaaaaaaa")
			err = handler(c, message)
			return
		}
	}
}

func TestProducer(t *testing.T) {
	var err error
	var pdc queue.Producer

	name := "default"
	addrs := []string{"127.0.0.1:9092"}
	pdc = New(
		name,
		addrs,
		WithTopic("message"),
		//WithAsync(true),
		WithMiddleware(TraceID()),
	)

	c := context.Background()
	msg := struct {
		MsgID string
		Tag   string
		Data  string
	}{"111", "aaa", "xxx"}

	for i := 0; i < 100000; i++ {
		time.Sleep(3 * time.Second)
		err = pdc.Send(c, msg)
		if err != nil {
			t.Errorf("%s has error[%+v]", t.Name(), err)
			return
		}
	}

}
