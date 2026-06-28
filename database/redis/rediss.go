package redis

/*
 * @abstract Redis client
 * @mail neo532@126.com
 * @date 2024-10-19
 */

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/neo532/gokit/database"
	"github.com/neo532/gokit/errorx"
)

// ========== RedissOpt =========
type RedissOpt func(*Rediss)

func WithBenchmarker(fn database.Benchmarker) RedissOpt {
	return func(o *Rediss) {
		o.benchmarker = fn
	}
}
func WithGrayer(fn database.Grayer) RedissOpt {
	return func(o *Rediss) {
		o.grayer = fn
	}
}
func WithDefault(rs ...*Redis) RedissOpt {
	return func(o *Rediss) {
		if o.def == nil {
			o.def = &Rdbs{}
		}
		setRdb(o.def, o, rs...)
	}
}
func WithShadow(rs ...*Redis) RedissOpt {
	return func(o *Rediss) {
		if o.shadow == nil {
			o.shadow = &Rdbs{}
		}
		setRdb(o.shadow, o, rs...)
	}
}
func WithGray(rs ...*Redis) RedissOpt {
	return func(o *Rediss) {
		if o.gray == nil {
			o.gray = &Rdbs{}
		}
		setRdb(o.gray, o, rs...)
	}
}

func setRdb(rs *Rdbs, o *Rediss, rdbs ...*Redis) {

	dbOldM := make(map[string]*Redis, len(rs.rdbs))
	for _, v := range rs.rdbs {
		dbOldM[v.Key()] = v
	}

	dbNew := make([]*Redis, 0, len(rdbs))

	var isUpdate bool
	for _, db := range rdbs {

		if err := db.Error(); err != nil {
			o.err = err
			continue
		}

		dbNew = append(dbNew, db)
		delete(dbOldM, db.Key())
		isUpdate = true
	}
	if isUpdate {

		rs.lock.Lock()
		rs.rdbs = dbNew
		defer rs.lock.Unlock()

		for _, v := range dbOldM {
			cleanUp(v)
		}
	}
}
func cleanUp(os ...*Redis) (err error) {
	for _, o := range os {
		t := time.NewTimer(
			time.Duration(int(o.redisLogger.slowLogTime.Seconds())+1) * time.Second,
		)
		go func() {
			<-t.C
			o.Close()()
		}()
	}

	return
}

type Rediss struct {
	def    *Rdbs
	shadow *Rdbs
	gray   *Rdbs

	benchmarker database.Benchmarker
	grayer      database.Grayer
	pooler      Pooler

	err error
}
type Rdbs struct {
	rdbs []*Redis
	lock sync.RWMutex
}

func News(opts ...RedissOpt) (rdbs *Rediss) {
	rdbs = &Rediss{
		benchmarker: &database.DefaultBenchmarker{},
		grayer:      &database.DefaultGrayer{},
		pooler:      &RandomPolicy{},
	}
	rdbs.With(opts...)
	return
}

func (d *Rediss) With(opts ...RedissOpt) {
	for _, opt := range opts {
		opt(d)
	}
	if d.def == nil {
		d.err = errorx.New("Nil Redis")
	}
}

func (d *Rediss) Rdb(c context.Context) (rdb *redis.Client) {
	if d.gray != nil && d.grayer.Judge(c) {
		return d.pooler.Choose(c, d.gray)
	}
	if d.shadow != nil && d.benchmarker.Judge(c) {
		return d.pooler.Choose(c, d.shadow)
	}
	return d.pooler.Choose(c, d.def)
}

func (d *Rediss) Close() func() {
	return func() {
		if d.def != nil {
			for _, o := range d.def.rdbs {
				o.Close()()
			}
		}
		if d.shadow != nil {
			for _, o := range d.shadow.rdbs {
				o.Close()()
			}
		}
		if d.gray != nil {
			for _, o := range d.gray.rdbs {
				o.Close()()
			}
		}
	}
}

func (d *Rediss) Error() error {
	return d.err
}
