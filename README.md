# Go's toolkit


[![Go Report Card](https://goreportcard.com/badge/github.com/neo532/gokit)](https://goreportcard.com/report/github.com/neo532/gokit)
[![Sourcegraph](https://sourcegraph.com/github.com/neo532/gokit/-/badge.svg)](https://sourcegraph.com/github.com/neo532/gokit?badge)

Gokit is a toolkit written by Go (Golang).It aims to speed up the development.


## Contents

- [Gofr Web Framework](#gofr-web-framework)
    - [Installation](#installation)
    - [Usage](#Usage)
        - [Crypt](#HTTP-request)
        - [Distributed-lock](#Distributed-lock)
        - [Database](#Database)
            - [Orm](#Orm)
            - [Redis](#Redis)
        - [Logger](#Logger)
        - [Queue](#Queue)
            - [Kafka](Kafka)
        - [Frequency controller](#Frequency-controller)
        - [Page Execute](#Page-Execute)
        - [Guard panic](#Guard-panic)



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
