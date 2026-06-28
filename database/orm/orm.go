package orm

/*
 * @abstract Orm client
 * @mail neo532@126.com
 * @date 2024-10-19
 */

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	//"gorm.io/gorm/hints"
	"github.com/neo532/gokit/database"
	ilogger "github.com/neo532/gokit/logger"
)

var (
	instanceLock sync.Mutex
	ormMap       = make(map[string]*Orm, 2)
)

// ========== Option ==========
type gormOpt struct {
	schema schema.NamingStrategy
}

type Option func(*Orm)

func WithMaxIdleConns(i int32) Option {
	return func(o *Orm) {
		o.opts = append(o.opts, func(db *sql.DB) {
			db.SetMaxIdleConns(int(i))
		})
		o.OptsHash["SetMaxIdleConns"] = i
	}
}
func WithMaxOpenConns(i int32) Option {
	return func(o *Orm) {
		o.opts = append(o.opts, func(db *sql.DB) {
			db.SetMaxOpenConns(int(i))
		})
		o.OptsHash["SetMaxOpenConns"] = i
	}
}
func WithConnMaxLifetime(t time.Duration) Option {
	return func(o *Orm) {
		o.opts = append(o.opts, func(db *sql.DB) {
			db.SetConnMaxLifetime(t)
			o.ConnMaxLifetime = t
		})
		o.OptsHash["ConnMaxLifetime"] = t
	}
}
func WithSlowLog(t time.Duration) Option {
	return func(o *Orm) {
		o.GormLogger.slowLogTime = t
		o.OptsHash["slowLogTime"] = t

	}
}
func WithTablePrefix(s string) Option {
	return func(o *Orm) {
		o.GormOpt.schema.TablePrefix = s
		o.OptsHash["TablePrefix"] = s

	}
}
func WithLogger(l ilogger.ILogger) Option {
	return func(o *Orm) {
		o.GormLogger.logger = l
		o.OptsHash["logger"] = l

	}
}
func WithSingularTable() Option {
	return func(o *Orm) {
		o.GormOpt.schema.SingularTable = true
		o.OptsHash["SingularTable"] = true
	}
}
func WithContext(c context.Context) Option {
	return func(o *Orm) {
		o.bootstrapContext = c
	}
}
func WithGormProcessor(fns ...func(*gorm.DB)) Option {
	return func(o *Orm) {
		o.processors = fns
	}
}

// ========== /Option ==========
type Orm struct {
	orm              *gorm.DB
	close            func()
	err              error
	bootstrapContext context.Context
	ConnMaxLifetime  time.Duration

	GormLogger *gormLogger
	GormOpt    *gormOpt
	opts       []func(db *sql.DB)
	processors []func(*gorm.DB)
	OptsHash   map[string]any

	key string
}

// New returns a instance of Orm.
// this Name must be unique to special instance.
func New(name string, dsn gorm.Dialector, opts ...Option) (db *Orm) {
	instanceLock.Lock()
	defer instanceLock.Unlock()

	db = &Orm{
		bootstrapContext: context.Background(),
		GormOpt: &gormOpt{
			schema: schema.NamingStrategy{},
		},
		GormLogger:      &gormLogger{Name: name, logger: ilogger.NewDefaultILogger()},
		opts:            make([]func(db *sql.DB), 0),
		ConnMaxLifetime: 3 * time.Second,
		OptsHash:        make(map[string]any, 3),
		key:             name + dsn.Name(),
	}
	for _, o := range opts {
		o(db)
	}

	if b, e := json.Marshal(dsn); e == nil {
		db.key += ":" + fmt.Sprintf("%x", md5.Sum(b))
	}

	if odb, ok := ormMap[db.key]; ok {
		db = odb
		return
	}

	db.orm, db.err = gorm.Open(
		dsn,
		&gorm.Config{
			Logger:         db.GormLogger,
			NamingStrategy: db.GormOpt.schema,
			ClauseBuilders: map[string]clause.ClauseBuilder{
				//hints.Comment("select", "master"),
			},
			//ClauseBuilders: map[string]hints.Comment("select", "master")clause.ClauseBuilder{},
		},
	)

	if db.err != nil {
		db.LogError("Gorm open client error")
		return
	}

	var sqlDB *sql.DB
	if sqlDB, db.err = db.orm.DB(); db.err != nil {
		db.LogError("Orm DB has error")
		return
	}
	for _, o := range db.opts {
		o(sqlDB)
	}

	if db.err = sqlDB.Ping(); db.err != nil {
		db.LogError("Orm DB has error")
		return
	}

	for _, fn := range db.processors {
		fn(db.orm)
	}

	db.close = func() {
		if sqlDB == nil {
			db.LogWarn("Close db is nil!")
			return
		}
		if db.err = sqlDB.Close(); db.err != nil {
			db.LogWarn("Close db has error!")
			return
		}
	}
	ormMap[name] = db
	return
}

func (o *Orm) LogError(message string) {
	o.GormLogger.logger.Error(
		o.bootstrapContext,
		message,
		database.KeyName, o.GormLogger.Name,
		database.KeyError, o.err,
	)
}
func (o *Orm) LogWarn(message string) {
	o.GormLogger.logger.Warn(
		o.bootstrapContext,
		message,
		database.KeyName, o.GormLogger.Name,
		database.KeyError, o.err,
	)
}

func (o *Orm) Error() error {
	return o.err
}

func (o *Orm) Key() string {
	return o.key
}

func (o *Orm) Orm() *gorm.DB {
	return o.orm
}

func (o *Orm) Close() func() {
	return o.close
}

type gormLogger struct {
	gorm.Config

	Name        string
	slowLogTime time.Duration
	logger      ilogger.ILogger

	LogLevel logger.LogLevel
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.LogLevel = level
	return g
}

func (g *gormLogger) Info(c context.Context, s string, i ...any) {
	g.logger.Info(c, s, i...)
}

func (g *gormLogger) Warn(c context.Context, s string, i ...any) {
	g.logger.Warn(c, s, i...)
}

func (g *gormLogger) Error(c context.Context, s string, i ...any) {
	g.logger.Error(c, s, i...)
}

func (g *gormLogger) Trace(c context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	cost := time.Since(begin)

	p := []any{
		"name", g.Name,
		"limit", g.slowLogTime.Seconds(),
		"cost", cost.Seconds(),
		"rows", rows,
	}

	if err != nil {
		p = append(p, "err", err)
		g.logger.Error(c, sql, p...)
		return
	}

	if cost > g.slowLogTime {
		g.logger.Warn(c, "[slow]"+sql, p...)
		return
	}

	g.logger.Info(c, sql, p...)
}
