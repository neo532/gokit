# Go's toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/neo532/gokit)](https://goreportcard.com/report/github.com/neo532/gokit)
[![Sourcegraph](https://sourcegraph.com/github.com/neo532/gokit/-/badge.svg)](https://sourcegraph.com/github.com/neo532/gokit?badge)

Gokit is a toolkit written by Go (Golang). It aims to speed up the development.

## Contents

- [Gokit](#Gokit)
  - [Installation](#installation)
  - [Usage](#Usage)
    - [Distributed lock](#Distributed-lock)
    - [Frequency limiter](#Frequency-limiter)
    - [Page execute](#Page-execute)
    - [Guard panic](#Guard-panic)
    - [Logger](#Logger)
      - [Slog](#Slog)
      - [Zap](#Zap)
    - [Database](#Database)
      - [Orm](#Orm)
      - [Redis](#Redis)
    - [Queue](#Queue)
      - [Kafka](#Kafka)
    - [File watcher](#File-watcher)
    - [Metadata](#Metadata)
    - [Middleware](#Middleware)
    - [HTTP client](#HTTP-client)
    - [Config generator](#Config-generator)
    - [Crypt](#Crypt)
      - [Converter](#Crypt-Converter)
      - [Encoding](#Crypt-Encoding)
      - [Marshaler](#Crypt-Marshaler)
      - [Openssl](#Openssl)

## Installation

To install Gokit package, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.google.cn/dl) installed (**version 1.21+ is required**), then you can use the below Go command to install Gokit.

```sh
    $ go get github.com/neo532/gokit
```

2. Import it in your code:

```go
    import "github.com/neo532/gokit"
```

### Distributed lock

It is a distributed lock with single instance by redis.

Also a `NoSpinLock` is provided for local lock without spinning.

[example](https://github.com/neo532/gokit/blob/master/lock/lockDistributed.go)

```go
    package main

    import (
        "github.com/go-redis/redis/v8"
        "github.com/neo532/gokit/lock"
    )

    type RedisOne struct {
        cache *redis.Client
    }

    func (l *RedisOne) Eval(c context.Context, cmd string, keys []string, args []any) (any, error) {
        return l.cache.Eval(c, cmd, keys, args...).Result()
    }

    var Lock *lock.DistributedLock

    func init() {
        rdb := &RedisOne{
            redis.NewClient(&redis.Options{
                Addr:     "127.0.0.1:6379",
                Password: "password",
            }),
        }
        Lock = lock.NewDistributedLock(rdb)
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

### Frequency limiter

It is a frequency with single instance by redis.

[example](https://github.com/neo532/gokit/blob/master/limiter/freq.go)

```go
    package main

    import (
        "github.com/go-redis/redis/v8"
        "github.com/neo532/gokit/limiter"
    )

    type RedisOne struct {
        cache *redis.Client
    }

    func (l *RedisOne) Eval(c context.Context, cmd string, keys []string, args []any) (any, error) {
        return l.cache.Eval(c, cmd, keys, args...).Result()
    }

    var Freq *limiter.Freq

    func init() {
        rdb := &RedisOne{
            redis.NewClient(&redis.Options{
                Addr:     "127.0.0.1:6379",
                Password: "password",
            }),
        }
        Freq = limiter.NewFreq(rdb)
        Freq.Timezone("Local")
    }

    func main() {
        c := context.Background()
        preKey := "user.test"
        rule := []limiter.FreqRule{
            {Duri: "10000", Times: 80},
            {Duri: "today", Times: 5},
        }

        fmt.Println(Freq.IncrCheck(c, preKey, rule...))
        fmt.Println(Freq.Get(c, preKey, rule...))
    }
```

### Page execute

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
        err := util.PageExec(int64(len(arr)), 3, func(b, e int64, p int) (err error) {
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

[example](https://github.com/neo532/gokit/blob/master/logger/slog/slog_test.go)

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

        l.WithArgs(logger.KeyModule, "db").Info(c, "msg1", "err", "panic")
        l.WithArgs(logger.KeyModule, "queue").WithLevel(logger.LevelError).Info(c, "m2", "e1", "p1")
        l.WithArgs(logger.KeyModule, "redis").Errorf(c, "m%s", "3")
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

### File watcher

A file watcher that monitors a directory for file changes, delivering both initial state and subsequent updates.

[example](https://github.com/neo532/gokit/blob/master/filepath/file_test.go)

```go
    package main

    import (
        "github.com/neo532/gokit/filepath"
    )

    func main() {
        w := filepath.New("/path/to/watch")

        // Watch delivers initial state for all existing files,
        // then async updates on writes/creates.
        err := w.Watch(ctx, func(fileName string, data []byte) error {
            fmt.Println(fileName, string(data))
            return nil
        })
    }
```

### Metadata

Metadata is a way of representing request headers internally, used at the RPC level to translate back and forth from transport headers.

[example](https://github.com/neo532/gokit/blob/master/metadata/metadata_test.go)

```go
    package main

    import (
        "github.com/neo532/gokit/metadata"
    )

    func main() {
        md := metadata.New(map[string][]string{
            "authorization": {"token123"},
        })

        ctx := metadata.NewClientContext(context.Background(), md)

        // Append more key-values later
        ctx = metadata.AppendToClientContext(ctx, "x-request-id", "req-123")

        md, ok := metadata.FromClientContext(ctx)
    }
```

### Middleware

HTTP/gRPC transport middleware with a chain helper.

[example](https://github.com/neo532/gokit/blob/master/middleware/middleware_test.go)

```go
    package main

    import (
        "github.com/neo532/gokit/middleware"
    )

    func logging() middleware.Middleware {
        return func(next middleware.Handler) middleware.Handler {
            return func(c context.Context, request, reply any) (context.Context, error) {
                log.Println("before")
                c, err := next(c, request, reply)
                log.Println("after")
                return c, err
            }
        }
    }

    func main() {
        chain := middleware.Chain(logging(), recovery())
        chain(func(c context.Context, request, reply any) (context.Context, error) {
            // handler logic
            return c, nil
        })(ctx, req, &reply)
    }
```

### HTTP client

A feature-rich HTTP client with middleware chain, connection pool management, TLS configuration, retry, and logging support.

[example](https://github.com/neo532/gokit/blob/master/transport/http/client/request.go)

```go
    package main

    import (
        "github.com/neo532/gokit/transport/http/client"
    )

    func main() {
        c := client.NewClient(
            client.WithLogger(myLogger),
            client.WithDefaultRetryTimes(3),
            client.WithDefaultTimeLimit(5*time.Second),
            client.WithMaxConnsPerHost(10),
            client.WithInsecureSkipVerify(true),
        )
        // Use c.HttpClient() to make requests
    }
```

### Config generator

A tool to generate Go struct definitions from configuration files (JSON, YAML, INI).

[example](https://github.com/neo532/gokit/blob/master/cmd/config-gen-go-struct/main.go)

```sh
    $ go run github.com/neo532/gokit/cmd/config-gen-go-struct -f config.yaml
```

### Crypt

Crypt provides a comprehensive cryptography suite including encryption/decryption, encoding, data conversion, and serialization.

#### Converter

Convert data between different formats and types.

[example](https://github.com/neo532/gokit/blob/master/crypt/converter/converter.go)

#### Encoding

Encoding and decoding support for standard and URL-safe formats.

[example](https://github.com/neo532/gokit/blob/master/crypt/encoding/encoding.go)

#### Marshaler

Marshal and unmarshal data with JSON and XML support.

[example](https://github.com/neo532/gokit/blob/master/crypt/marshaler/marshaler.go)

#### Openssl

Openssl-compatible encryption implementations.

- **CBC** — Cipher block chaining mode
- **ECB** — Electronic codebook mode
- **RSA** — RSA encryption and signing

[example](https://github.com/neo532/gokit/blob/master/crypt/crypt/openssl/cbc/cbc.go)

```go
    import (
        "github.com/neo532/gokit/crypt/crypt/openssl/cbc"
        "github.com/neo532/gokit/crypt/crypt/openssl/rsa"
    )
```
