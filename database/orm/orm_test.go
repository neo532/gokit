package orm

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

func getConfig() (d *Config) {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	d = &Config{
		MaxOpenConns:        2,
		MaxIdleConns:        2,
		ConnMaxLifetime:     3 * time.Second,
		MaxSlowtime:         3 * time.Second,
		TablePrefix:         "",
		RecordNotFoundError: false,
		Read:                []*DsnConfig{&DsnConfig{Name: "default_read", Dsn: dsn}},
		Write:               []*DsnConfig{&DsnConfig{Name: "default_write", Dsn: dsn}},
		ShadowRead:          []*DsnConfig{&DsnConfig{Name: "default_shadowread", Dsn: dsn}},
		ShadowWrite:         []*DsnConfig{&DsnConfig{Name: "default_shadowwrite", Dsn: dsn}},
	}
	return
}

func initDB() (dbs *Orms, clean func(), err error) {
	logger := logger.NewDefaultILogger()
	c := context.Background()
	var d *Config

	d = getConfig()
	dbs, clean, err = NewOrms(c, d, logger)

	d = getConfig()
	dbs, clean, err = With(c, dbs, d, logger)
	return
}

func TestOrms(t *testing.T) {

	dbs, clean, err := initDB()
	defer clean()
	if err != nil {
		t.Error(err)
		return
	}

	c := context.Background()
	var databases []string
	if err = dbs.Write(c).Raw("show databases").Scan(&databases).Error; err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}
	fmt.Println(fmt.Sprintf("dbs:%+v\t%+v", databases, err))

	if err = dbs.Read(c).Raw("show databases").Scan(&databases).Error; err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}
	fmt.Println(fmt.Sprintf("dbs:%+v\t%+v", databases, err))
	time.Sleep(10 * time.Second)
}

func TestTransaction(t *testing.T) {

	dbs, clean, err := initDB()
	defer clean()
	if err != nil {
		t.Error(err)
		return
	}

	c := context.Background()
	err = dbs.Transaction(c, func(c context.Context) (err error) {

		var databases []string

		if err = dbs.Write(c).Raw("show databases").Scan(&databases).Error; err != nil {
			return
		}

		if err = dbs.Read(c).Raw("show databases").Scan(&databases).Error; err != nil {
			return
		}
		fmt.Println(fmt.Sprintf("txdbs:%+v", databases))
		return
	})
	if err != nil {
		t.Errorf("%s has err[%+v]", t.Name(), err)
	}
}
