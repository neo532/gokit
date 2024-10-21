# Go's toolkit


[![Go Report Card](https://goreportcard.com/badge/github.com/neo532/gokit)](https://goreportcard.com/report/github.com/neo532/gokit)
[![Sourcegraph](https://sourcegraph.com/github.com/neo532/gokit/-/badge.svg)](https://sourcegraph.com/github.com/neo532/gokit?badge)

Gokit is a toolkit written by Go (Golang).It aims to speed up the development.


## Contents

- [Gofr Web Framework](#gofr-web-framework)
    - [Installation](#installation)
    - [Usage](#Usage)
        - [Frequency controller](#Frequency-controller)
        - [Logger](#Logger)
            - [Slog](#Slog)
            - [Zap](#Zap)
        - [Database](#Database)
            - [Orm](#Orm)
            - [Redis](#Redis)
        - [Queue](#Queue)
            - [Kafka](Kafka)
        - [Guard panic](#Guard-panic)
        - [Distributed-lock](#Distributed-lock)
        - [Page Execute](#Page-Execute)
        - [Crypt](#Crypt)
            - [Openssl](#Openssl)
                - [Cbc](#Cbc)
                - [Ecb](#Ecb)
                - [Rsa](#Rsa)



## Installation

To install Gokit package, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.google.cn/dl) installed (**version 1.21+ is required**), then you can use the below Go command to install Gokit.

```sh
    $ go install github.com/neo532/gokit
```

2. Import it in your code:

```go
    import "github.com/neo532/gokit"
```

### Distributed lock

It is a distributed lock with signle instance by redis.

[example](https://github.com/neo532/gokit/blob/master/util/lockDistributed.go)

```go
    package main

    import (
        "github.com/go-redis/redis/v8"
        "github.com/neo532/gokit/util"
    )

    type RedisOne struct {
        cache *redis.Client
    }

    func (l *RedisOne) Eval(c context.Context, cmd string, keys []string, args []interface{}) (interface{}, error) {
        return l.cache.Eval(c, cmd, keys, args...).Result()
    }

    var Lock *util.Lock

    func init(){

        rdb := &RedisOne{
            redis.NewClient(&redis.Options{
                Addr:     "127.0.0.1:6379",
                Password: "password",
            })
        }

        Lock = util.NewLock(rdb)
    }

    func main() {

        c := context.Background()
        key := "IamAKey"
        expire := time.Duration(10) * time.Second
        wait := time.Duration(2) * time.Second

        code, err := Lock.Lock(c, key, expire, wait)
        Lock.UnLock(c, key, code)
    }
```

### Frequency controller

It is a frequency with signle instance by redis.

[example](https://github.com/neo532/gokit/blob/master/util/freq.go)

```go
    package main

    import (
        "github.com/go-redis/redis/v8"
        "github.com/neo532/gokit/util"
    )

    type RedisOne struct {
        cache *redis.Client
    }

    func (l *RedisOne) Eval(c context.Context, cmd string, keys []string, args []interface{}) (interface{}, error) {
        return l.cache.Eval(c, cmd, keys, args...).Result()
    }

    var Freq *util.Freq

    func init(){

        rdb := &RedisOne{
            redis.NewClient(&redis.Options{
                Addr:     "127.0.0.1:6379",
                Password: "password",
            })
        }

        Freq = util.NewFreq(rdb)
        Freq.Timezone("Local")
    }

    func main() {

        c := context.Background()
        preKey := "user.test"
        rule := []util.FreqRule{
            tool.FreqRule{Duri: "10000", Times: 80},
            tool.FreqRule{Duri: "today", Times: 5},
        }

        fmt.Println(Freq.IncrCheck(c, preKey, rule...))
        fmt.Println(Freq.Get(c, preKey, rule...))
    }
```

### Page Execute

It is a tool to page slice.

[example](https://github.com/neo532/gokit/blob/master/util/pageExec.go)

```go
    package main

    import (
        "github.com/neo532/gokit/util"
    )

    func main() {

        arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

        // [1 2 3] [4 5 6] [7 8 9] [10]
        err := PageExec(int64(len(arr)), 3, func(b, e int64, p int) (err error) {
            fmt.Println(arr[b:e])
            return
        })
    }
```

### Guard panic

It is a tool to exec goroutine safely.

[example](https://github.com/neo532/gokit/blob/master/gofunc/gofunc_test.go)

```go
    package main

    import (
        "github.com/neo532/gokit/gofunc"
    )

    func main() {

		fn := func(i int) (err error) {
			// do something...
			return
		}

		log := &gofunc.DefaultLogger{}
		gofn := gofunc.NewGoFunc(gofunc.WithLogger(log), gofunc.WithMaxGoroutine(20))

		l := 1000000
		fns := make([]func(i int) error, 0, l)
		for i := 0; i < l; i++ {
			fns = append(fns, fn)
		}

		gofn.WithTimeout(c, time.Second*2, fns...)
		err := log.Err()
    }
```

### Logger

It is a highly scalable logger.

[example](https://github.com/neo532/gokit/blob/master/logger/slog/log_test.go)

```go
    package main

    import (
        "github.com/neo532/gokit/logger"
    )

    func main() {

        c := context.Background()
        var l logger.Logger
        l = newSlog() // more detail in test file
        l = newZap()  // more detail in test file

        h.WithArgs(logger.KeyModule, "db").Error(c, "bug", "err", "panic")
        h.WithArgs(logger.KeyModule, "queue").WithLevel(logger.LevelFatal).Error(c, "b1", "err", "p1")
        h.WithArgs(logger.KeyModule, "redis").Errorf(c, "kkkk%s", "cc")
    }
```

### Orm

A well-encapsulated GORM that can support shadow databases, hot configuration updates, master-slave separation, high scalability, simplicity, and domain-layer transactions.

[example](https://github.com/neo532/gokit/blob/master/database/orm/orm_test.go)

```go
    package main

    import (
        "github.com/neo532/gokit/database/orm"
    )

    func main() {

        db, clean, err := initDB() // more detail in test file
        defer clean()
        if err != nil {
            return
        }

        c := context.Background()
        err = db.Transaction(c, func(c context.Context) (err error) {

            var databases []string

            if err = dbs.Write(c).Raw("show databases").Scan(&databases).Error; err != nil {
                return
            }

            if err = dbs.Read(c).Raw("show databases").Scan(&databases).Error; err != nil {
                return
            }
            return
        })
    }
```


### Redis

A well-encapsulated Redis client that can support shadow databases, hot configuration updates, gray environment, high scalability and simplicity.

[example](https://github.com/neo532/gokit/blob/master/database/redis/redis_test.go)

```go
    // more detail in test file
```

### Queue

A message queue client with high scalability that supports the full-link connection between producers and consumers and also supports customizable middleware.

[producer](https://github.com/neo532/gokit/blob/master/queue/kafka/producer/producer_test.go)
[consumer](https://github.com/neo532/gokit/blob/master/queue/kafka/consumergroup/consumergroup_test.go)

```go
    // more detail in test file
```
