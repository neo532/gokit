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
	d := []*ConnectConfig{&ConnectConfig{Name: "default", Addr: "127.0.0.1:6379"}}
	rdbs = &Config{
		MaxSlowtime: 3 * time.Second,
		Default:     d,
		Shadow:      d,
		Gray:        d,
	}
	return
}

func InitRDB() (rdbs *Rediss, clean func(), err error) {
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

	rdbs, clean, err := InitRDB()
	defer clean()
	if err != nil {
		t.Error(err)
		return
	}

	c := context.Background()
	key := "database.redis.testkey"
	var r string
	if r, err = rdbs.Rdb(c).SetEX(c, key, "aaaa", 10*time.Minute).Result(); err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}
	fmt.Println(fmt.Sprintf("dbs:%+v\t%+v", r, err))

	if r, err = rdbs.Rdb(c).Get(c, key).Result(); err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}
	fmt.Println(fmt.Sprintf("dbs:%+v\t%+v", r, err))
	time.Sleep(10 * time.Second)
}
