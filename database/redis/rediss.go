package redis

/*
 * @abstract Redis client
 * @mail neo532@126.com
 * @date 2024-10-19
 */

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/neo532/gokit/database"
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

	var ok bool
	for _, db := range rdbs {

		if err := db.Error(); err != nil {
			o.err = err
			continue
		}

		if _, ok := dbOldM[db.Key()]; ok {
			delete(dbOldM, db.Key())
		}

		dbNew = append(dbNew, db)

		if !ok {
			ok = true
		}
	}
	if ok {

		rs.lock.Lock()
		rs.rdbs = dbNew
		defer rs.lock.Unlock()

		for _, v := range dbOldM {
			cleanUp(v)
		}
	}
	return
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
	if len(opts) > 0 {
		rdbs.With(opts...)
	}
	return
}

func (d *Rediss) With(opts ...RedissOpt) {
	for _, opt := range opts {
		opt(d)
	}
	if d.def == nil {
		d.err = errors.New("Please input a instance at least")
	}
}

func (d *Rediss) Gray(c context.Context) (rdb *redis.Client) {
	if d.grayer.Judge(c) && d.gray != nil {
		return d.pooler.Choose(c, d.gray)
	}
	return d.Rdb(c)
}

func (d *Rediss) Rdb(c context.Context) (rdb *redis.Client) {
	if d.benchmarker.Judge(c) && d.shadow != nil {
		return d.pooler.Choose(c, d.shadow)
	}
	return d.pooler.Choose(c, d.def)
}

func (d *Rediss) Close() func() {
	return func() {
		if d.def != nil {
			for _, o := range d.def.rdbs {
				o.Close()
			}
		}
		if d.shadow != nil {
			for _, o := range d.shadow.rdbs {
				o.Close()
			}
		}
		if d.gray != nil {
			for _, o := range d.gray.rdbs {
				o.Close()
			}
		}
	}
}

func (d *Rediss) Error() error {
	return d.err
}
