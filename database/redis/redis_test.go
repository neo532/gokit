package redis

/*
 * @abstract Orm client
 * @mail neo532@126.com
 * @date 2024-05-18
 */

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/neo532/gokit/logger"
)

func getConfig() (rdbs *Config) {
	addr := "127.0.0.1:6379"
	rdbs = &Config{
		MaxSlowtime: 3 * time.Second,
		Default:     []*ConnectConfig{{Name: "default", Addr: addr}},
		Shadow:      []*ConnectConfig{{Name: "shadow", Addr: addr}},
		Gray:        []*ConnectConfig{{Name: "gray", Addr: addr}},
	}
	return
}

func initRDB() (rdbs *Rediss, clean func(), err error) {
	logger := logger.NewDefaultILogger()
	c := context.Background()
	var d *Config

	d = getConfig()
	rdbs, clean, err = NewRediss(c, d, logger)

	d = getConfig()
	rdbs, clean, err = With(c, rdbs, d, logger)
	return
}

func TestRediss(t *testing.T) {

	rdbs, clean, err := initRDB()
	defer func() {
		clean()
	}()
	if err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
		return
	}

	c := context.Background()
	key := "database.redis.testkey"
	value := "aaaa"
	if _, err := rdbs.Rdb(c).SetEX(c, key, value, 10*time.Minute).Result(); err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}

	if _, err := rdbs.Rdb(c).Get(c, key).Result(); err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}

	fmt.Println(t.Name())
}
