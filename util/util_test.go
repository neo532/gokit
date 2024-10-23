package util

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestCompareVersion(t *testing.T) {

	cv := func(v1, v2 string, r int) {
		if v := CompareVersion(v1, v2); v != r {
			t.Errorf("%s: %s-%s is %d should %d", t.Name(), v1, v2, v, r)
		}
	}

	cv("1.2", "1.2.3", Smaller)
	cv("1.4", "1.2.3", Larger)
	cv("1.2.3", "1.2.3", Equal)

	fmt.Println(t.Name())
}

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

func TestPageExec(t *testing.T) {

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	m := map[int][]int{
		1: []int{1, 2, 3},
		2: []int{4, 5, 6},
		3: []int{7, 8, 9},
		4: []int{10},
	}

	// [1 2 3] [4 5 6] [7 8 9] [10]
	err := PageExec(int64(len(arr)), 3, func(b, e int64, p int) (err error) {
		a := arr[b:e]
		if v, ok := m[p]; ok && !CompareSlice(a, v) {
			t.Errorf("%s is %+v should %+v", t.Name(), a, v)
		}
		return
	})
	if err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	fmt.Println(t.Name())
}

func TestFreq(t *testing.T) {
	var err error
	c := context.Background()

	freq := NewFreq(newRDB())
	freq.Timezone("Local")

	preKey := "user.test"
	dstTimes := 2
	rule := []FreqRule{
		{Duri: "5", Times: int64(dstTimes)},
	}

	count := 10000
	var wg sync.WaitGroup
	wg.Add(count)
	num := 0
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			if ok, err := freq.IncrCheck(c, preKey, rule...); err == nil && ok {
				num++
			}
		}()
	}
	wg.Wait()
	if dstTimes != num {
		t.Errorf("%s has wrong %d should %d", t.Name(), num, dstTimes)
	}

	if _, err = freq.IncrCheck(c, preKey, rule...); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	var times int64
	if times, err = freq.Get(c, preKey, rule...); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}
	if dst := int64(count) + 1; dst != times {
		t.Errorf("%s has wrong %d should %d", t.Name(), times, dst)
	}

	fmt.Println(t.Name())
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
