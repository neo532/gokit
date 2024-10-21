package util

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
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

	fmt.Println(fmt.Sprintf("%s\t:num:\t%+v", t.Name(), num))
}

func TestPageExec(t *testing.T) {

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// [1 2 3] [4 5 6] [7 8 9] [10]
	err := PageExec(int64(len(arr)), 3, func(b, e int64, p int) (err error) {
		fmt.Println(fmt.Sprintf("%s\t>p,b,e,arr:%d,%d,%d,%+v", t.Name(), p, b, e, arr[b:e]))
		return
	})
	if err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}
}

func TestFreq(t *testing.T) {
	var err error
	c := context.Background()

	freq := NewFreq(newRDB())
	freq.Timezone("Local")

	preKey := "user.test"
	rule := []FreqRule{
		{Duri: "10", Times: 2},
	}

	count := 10000
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			if ok, err := freq.IncrCheck(c, preKey, rule...); err == nil && ok {
				fmt.Println(fmt.Sprintf("%s\t:Biz run!", t.Name()))
			}
		}()
	}
	wg.Wait()

	var b bool
	b, err = freq.IncrCheck(c, preKey, rule...)
	fmt.Println(fmt.Sprintf("%s\t:IncrCheck!,%v,%v", t.Name(), b, err))
	if err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	var times int64
	times, err = freq.Get(c, preKey, rule...)
	fmt.Println(fmt.Sprintf("%s\t:Get!,%d,%v", t.Name(), times, err))
	if err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}
}

func TestDistributedLock(t *testing.T) {

	l := NewDistributedLock(newRDB())

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
				fmt.Println(fmt.Sprintf("%s\t:Biz run!", t.Name()))
			}
		}()
	}
	wg.Wait()

	if err = l.UnLock(c, key, code); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}
}

type rDB struct {
	cache *redis.Client
}

func newRDB() *rDB {
	return &rDB{redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})}
}

func (l *rDB) Eval(
	c context.Context,
	cmd string,
	keys []string,
	args []interface{}) (interface{}, error) {
	return l.cache.Eval(c, cmd, keys, args...).Result()
}

func (l *rDB) Get(c context.Context, key string) (string, error) {
	return l.cache.Get(c, key).Result()
}
