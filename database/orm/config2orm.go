package orm

import (
	"context"

	"gorm.io/driver/mysql"

	"github.com/neo532/gokit/logger"
)

/*
 * @abstract Orm's Logger
 * @mail neo532@126.com
 * @date 2024-05-18
 */

func Connect(c context.Context, cfg *Config, dsn *DsnConfig, logger logger.ILogger) *Orm {
	return New(
		dsn.Name,
		mysql.Open(dsn.Dsn),
		WithTablePrefix(cfg.TablePrefix),
		WithConnMaxLifetime(cfg.ConnMaxLifetime),
		WithMaxIdleConns(cfg.MaxIdleConns),
		WithMaxOpenConns(cfg.MaxOpenConns),
		WithLogger(logger),
		WithSingularTable(),
		WithContext(c),
		WithSlowLog(cfg.MaxSlowtime),
		WithRecordNotFoundError(cfg.RecordNotFoundError),
	)
}

func NewOrms(c context.Context, d *Config, l logger.ILogger) (*Orms, func(), error) {
	dbs := News()
	With(c, dbs, d, l)
	return dbs, dbs.Close(), dbs.Error()
}

func With(c context.Context, dbs *Orms, d *Config, l logger.ILogger) (*Orms, func(), error) {
	opts := make([]OrmsOpt, 0, 4)
	if d.Read != nil {
		for _, dsn := range d.Read {
			opts = append(opts, WithRead(Connect(c, d, dsn, l)))
		}
	}
	if d.Write != nil {
		for _, dsn := range d.Write {
			opts = append(opts, WithWrite(Connect(c, d, dsn, l)))
		}
	}
	if d.ShadowRead != nil {
		for _, dsn := range d.ShadowRead {
			opts = append(opts, WithShadowRead(Connect(c, d, dsn, l)))
		}
	}
	if d.ShadowWrite != nil {
		for _, dsn := range d.ShadowWrite {
			opts = append(opts, WithShadowWrite(Connect(c, d, dsn, l)))
		}
	}
	dbs.With(opts...)
	return dbs, dbs.Close(), dbs.Error()
}
