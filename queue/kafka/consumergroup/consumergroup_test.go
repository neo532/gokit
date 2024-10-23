package consumergroup

/*
 * @abstract consumer
 * @mail neo532@126.com
 * @date 2024-10-21
 */

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/neo532/gokit/queue"
)

func TraceID() queue.ConsumerMiddleware {
	return func(handler queue.ConsumerHandler) queue.ConsumerHandler {
		return func(c context.Context, message []byte) (err error) {

			if h, ok := queue.GetHeaderFromContext(c); ok {
				fmt.Println(fmt.Sprintf("h.Value(traceID):\t%+v", h.Value("traceID")))
				c = context.WithValue(c, "traceID", h.Value("traceID"))
			}
			err = handler(c, message)
			return
		}
	}
}

func TestConsumer(t *testing.T) {
	var err error
	c, cancel := context.WithCancel(context.Background())
	var csm queue.Consumer

	addr := []string{"127.0.0.1:9092"}
	if csm, err = NewGroup(
		"default",
		addr,
		"sender",
		WithTopics("message"),
		WithHandler(func(ctx context.Context, message []byte) (err error) {
			// do something...
			return
		}),
		WithMiddleware(TraceID()),
	); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}
	go func() {
		for {
			select {
			case <-time.After(9 * time.Minute):
				cancel()
			}
		}
	}()

	select {
	case <-c.Done():
		fmt.Println(t.Name())
		return
	default:
		if err = csm.Start(c); err != nil {
			t.Errorf("%s has error[%+v]", t.Name(), err)
			return
		}
	}
}
