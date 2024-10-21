package consumergroup

/*
 * @abstract consumer
 * @mail neo532@126.com
 * @date 2024-10-21
 */

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/neo532/gokit/queue"
)

func TestConsumer(t *testing.T) {
	var err error
	c, cancel := context.WithCancel(context.Background())
	var csm queue.Consumer

	addr := []string{"127.0.0.1:9092"}
	csm = NewGroup(
		"default",
		addr,
		"sender",
		WithTopics("message"),
		WithHandler(func(ctx context.Context, message []byte) (err error) {
			fmt.Println(runtime.Caller(0))
			panic(0)
			return
		}),
	)
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				cancel()
			}
		}
	}()

	for {
		time.Sleep(4 * time.Second)
		select {
		case <-c.Done():
			break
		default:
			if err = csm.Start(c); err != nil {
				t.Errorf("%s has error[%+v]", t.Name(), err)
				return
			}

		}
	}
	fmt.Println(runtime.Caller(0))

}