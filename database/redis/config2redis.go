package redis

import (
	"context"

	"github.com/neo532/gokit/database"
)

/*
 * @abstract Redis's Creator
 * @mail neo532@126.com
 * @date 2024-10-19
 */

func Connect(c context.Context, cfg *Config, rc *ConnectConfig, l database.Logger) *Redis {
	return New(
		rc.Name,
		rc.Addr,
		WithLogger(l),
		WithContext(c),
		WithDb(rc.DB),
		WithPassword(rc.Password),
	)
}

func NewRediss(c context.Context, d *Config, l database.Logger) (*Rediss, func(), error) {
	rdbs := News()
	With(c, rdbs, d, l)
	return rdbs, rdbs.Close(), rdbs.Error()
}

func With(c context.Context, rdbs *Rediss, d *Config, l database.Logger) (*Rediss, func(), error) {
	opts := make([]RedissOpt, 0, 4)
	if d.Default != nil {
		for _, v := range d.Default {
			opts = append(opts, WithDefault(Connect(c, d, v, l)))
		}
	}
	if d.Shadow != nil {
		for _, v := range d.Shadow {
			opts = append(opts, WithShadow(Connect(c, d, v, l)))
		}
	}
	if d.Gray != nil {
		for _, v := range d.Gray {
			opts = append(opts, WithGray(Connect(c, d, v, l)))
		}
	}
	rdbs.With(opts...)
	return rdbs, rdbs.Close(), rdbs.Error()
}
