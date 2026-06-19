package lock

import (
	"fmt"
	"sync"
	"testing"
)

func TestNoSpinLock(t *testing.T) {
	var num int
	lock := &NoSpinLock{}

	l := 1000
	var wg sync.WaitGroup
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func() {
			defer wg.Done()
			if lock.Lock() {
				num++
			}
		}()
	}
	wg.Wait()
	lock.Unlock()

	if num != 1 {
		t.Errorf("%s: has wrong num %d should %d", t.Name(), num, 1)
	}

	fmt.Println(t.Name())
}
