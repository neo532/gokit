package util

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	u "github.com/neo532/gokit/util"
)

func TestFreq(t *testing.T) {
	var err error
	c := context.Background()

	freq := u.NewFreq(newRDB())
	freq.Timezone("Local")

	preKey := "user.test"
	dstTimes := 2
	tz, _ := time.LoadLocation("UTC")
	rule := []u.FreqRule{
		//{Duri: "5", Times: int64(dstTimes)},
		{Duri: u.DurationToday, Times: int64(dstTimes), Timezone: tz, N: 2},
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
	return

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

	l := u.NewDistributedLock(newRDB())

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
	args []interface{}) (interface{}, error) {
	return l.cache.Eval(c, cmd, keys, args...).Result()
}

func (l *rDB) Get(c context.Context, key string) (string, error) {
	return l.cache.Get(c, key).Result()
}
