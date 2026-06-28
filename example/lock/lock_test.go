package lock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/neo532/gokit/lock"
)

func TestNoSpinLock(t *testing.T) {
	var num int
	lock := &lock.NoSpinLock{}

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

func TestDistributedLock(t *testing.T) {

	l := lock.NewDistributedLock(newRDB())

	var code string
	var err error
	c := context.Background()
	key := "IamAKey"

	count := 10000
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			if tc, err := l.Lock(c, key, 10*time.Second, 2*time.Second); err == nil {
				code = tc
				//fmt.Println(fmt.Sprintf("%s\t:Biz run!", t.Name()))
			}
		}()
	}
	wg.Wait()

	if err = l.UnLock(c, key, code); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	fmt.Println(t.Name())
}

type rDB struct {
	cache *redis.Client
}

func newRDB() (r *rDB) {
	r = &rDB{redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})}
	return
}

func (l *rDB) Eval(
	c context.Context,
	cmd string,
	keys []string,
	args []any) (any, error) {
	return l.cache.Eval(c, cmd, keys, args...).Result()
}

func (l *rDB) Get(c context.Context, key string) (string, error) {
	return l.cache.Get(c, key).Result()
}
