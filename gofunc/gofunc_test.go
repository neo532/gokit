package gofunc

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {

	var num atomic.Int32
	fn := func(i int) (err error) {
		time.Sleep(time.Second * 2)
		num.Add(1)
		// do something...
		if i == 1 {
			//err = errors.New("aaaaaaa")
			//panic("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		}
		return
	}

	log := &DefaultLogger{}
	gofn := NewGoFunc(WithLogger(log), WithMaxGoroutine(20))

	l := 5
	fns := make([]func(i int) error, 0, l)
	for i := 0; i < l; i++ {
		fns = append(fns, fn)
	}

	c, closeFn := context.WithCancel(context.Background())
	go func() {
		select {
		case <-time.After(1 * time.Second):
			closeFn()
		}
	}()

	gofn.WithTimeout(
		c,
		time.Second*2,
		fns...,
	)
	if err := log.Err(); err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}
	if i := int(num.Load()); i != l {
		t.Errorf("%s has wrong %d should %d", t.Name(), i, l)
	}

	fmt.Println(t.Name())
}
