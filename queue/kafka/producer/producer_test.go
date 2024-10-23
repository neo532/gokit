package producer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/neo532/gokit/queue"
)

func TraceID() queue.ProducerMiddleware {
	return func(handler queue.ProducerHandler) queue.ProducerHandler {
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
	c := context.Background()

	name := "default"
	addrs := []string{"127.0.0.1:9092"}
	pdc = New(
		name,
		addrs,
		WithTopic("message"),
		//WithAsync(true),
		WithMiddleware(TraceID()),
	)
	defer pdc.Close()
	if err = pdc.Error(); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	msg := struct {
		MsgID string
		Tag   string
		Data  string
	}{"111", "aaa", "xxx"}

	for i := 0; i < 10; i++ {

		time.Sleep(1 * time.Second)

		if err = pdc.Send(c, msg); err != nil {
			t.Errorf("%s has error[%+v]", t.Name(), err)
			return
		}
	}
	fmt.Println(t.Name())
}
