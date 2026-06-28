package limiter

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/neo532/gokit/limiter"
)

func TestFreq(t *testing.T) {
	var err error
	c := context.Background()

	freq := limiter.NewFreq(newRDB())
	freq.Timezone("Local")

	preKey := "user.test"
	dstTimes := 2
	tz, _ := time.LoadLocation("UTC")
	rule := []limiter.FreqRule{
		//{Duri: "5", Times: int64(dstTimes)},
		{Duri: limiter.DurationToday, Times: int64(dstTimes), Timezone: tz, N: 2},
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
