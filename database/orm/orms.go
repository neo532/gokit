package orm

/*
 * @abstract Orm client
 * @mail neo532@126.com
 * @date 2024-05-18
 */

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/neo532/gokit/database"
	"gorm.io/gorm"
)

type contextTransactionKey struct{}

type Orms struct {
	read        *DBs
	write       *DBs
	shadowRead  *DBs
	shadowWrite *DBs

	pooler      Pooler
	benchmarker database.Benchmarker

	err error
}

type DBs struct {
	dbs  []*Orm
	lock sync.RWMutex
}

// ========== OrmsOpt =========
type OrmsOpt func(*Orms)

func WithBenchmarker(fn database.Benchmarker) OrmsOpt {
	return func(o *Orms) {
		o.benchmarker = fn
	}
}

func WithPooler(fn Pooler) OrmsOpt {
	return func(o *Orms) {
		o.pooler = fn
	}
}

func WithRead(dbs ...*Orm) OrmsOpt {
	return func(o *Orms) {
		if o.read == nil {
			o.read = &DBs{}
		}
		setDB(o.read, o, dbs...)
	}
}

func WithWrite(dbs ...*Orm) OrmsOpt {
	return func(o *Orms) {
		if o.write == nil {
			o.write = &DBs{}
		}
		setDB(o.write, o, dbs...)
	}
}

func WithShadowRead(dbs ...*Orm) OrmsOpt {
	return func(o *Orms) {
		if o.shadowRead == nil {
			o.shadowRead = &DBs{}
		}
		setDB(o.shadowRead, o, dbs...)
	}
}

func WithShadowWrite(dbs ...*Orm) OrmsOpt {
	return func(o *Orms) {
		if o.shadowWrite == nil {
			o.shadowWrite = &DBs{}
		}
		setDB(o.shadowWrite, o, dbs...)
	}
}

func setDB(rs *DBs, o *Orms, dbs ...*Orm) {

	dbOldM := make(map[string]*Orm, len(rs.dbs))
	for _, v := range rs.dbs {
		dbOldM[v.Key()] = v
	}

	dbNew := make([]*Orm, 0, len(dbs))

	var ok bool
	for _, db := range dbs {

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
		rs.dbs = dbNew
		defer rs.lock.Unlock()

		for _, v := range dbOldM {
			cleanUp(v)
		}
	}
	return
}

func cleanUp(os ...*Orm) (err error) {
	for _, o := range os {
		t := time.NewTimer(
			time.Duration(int(o.ConnMaxLifetime.Seconds())+1) * time.Second,
		)
		go func() {
			<-t.C
			o.Close()()
		}()
	}

	return
}

// ========== /OrmsOpt =========

func News(opts ...OrmsOpt) (dbs *Orms) {
	dbs = &Orms{
		benchmarker: &database.DefaultBenchmarker{},
		pooler:      &RandomPolicy{},
	}
	if len(opts) > 0 {
		dbs.With(opts...)
	}
	return
}

func (d *Orms) With(opts ...OrmsOpt) {
	for _, opt := range opts {
		opt(d)
	}
	if d.write == nil && d.read == nil {
		d.err = errors.New("Please input a instance at least")
	}
}

func (d *Orms) get(c context.Context, dbs *DBs) (db *gorm.DB) {
	dbs.lock.RLock()
	defer dbs.lock.RUnlock()
	return d.pooler.Choose(c, dbs).WithContext(c)
}

func (d *Orms) Read(c context.Context) (db *gorm.DB) {
	if tx, ok := c.Value(contextTransactionKey{}).(*gorm.DB); ok {
		return tx
	}
	if d.benchmarker.Judge(c) {
		if d.shadowRead == nil {
			return d.get(c, d.shadowWrite)
		}
		return d.get(c, d.shadowRead)
	}
	if d.read == nil {
		return d.get(c, d.write)
	}
	return d.get(c, d.read)
}

func (d *Orms) Write(c context.Context) (db *gorm.DB) {
	if tx, ok := c.Value(contextTransactionKey{}).(*gorm.DB); ok {
		return tx
	}
	if d.benchmarker.Judge(c) {
		if d.shadowWrite == nil {
			return d.get(c, d.shadowRead)
		}
		return d.get(c, d.shadowWrite)
	}
	if d.write == nil {
		return d.get(c, d.read)
	}
	return d.get(c, d.write)
}

func (d *Orms) Transaction(c context.Context, fn func(c context.Context) (err error)) error {
	return d.Write(c).Transaction(func(tx *gorm.DB) error {
		if _, ok := c.Value(contextTransactionKey{}).(*gorm.DB); !ok {
			c = context.WithValue(c, contextTransactionKey{}, tx)
		}
		return fn(c)
	})
}

func (d *Orms) Close() func() {
	return func() {
		if d.read != nil {
			for _, o := range d.read.dbs {
				o.Close()()
			}
		}
		if d.write != nil {
			for _, o := range d.write.dbs {
				o.Close()()
			}
		}
		if d.shadowRead != nil {
			for _, o := range d.shadowRead.dbs {
				o.Close()()
			}
		}
		if d.shadowWrite != nil {
			for _, o := range d.shadowWrite.dbs {
				o.Close()()
			}
		}
	}
}

func (d *Orms) Error() error {
	return d.err
}
