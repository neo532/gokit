package redis

/*
 * @abstract Redis client
 * @mail neo532@126.com
 * @date 2024-10-19
 */

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/neo532/gokit/database"
	"github.com/neo532/gokit/logger"
)

// ========== Option ==========
type Option func(*Redis)

func WithMaxRetries(i int) Option {
	return func(o *Redis) {
		o.redisOpt.MaxRetries = i
		o.optsHash["MaxRetries"] = i
	}
}
func WithReadTimeout(t time.Duration) Option {
	return func(o *Redis) {
		o.redisOpt.ReadTimeout = t
		o.optsHash["ReadTimeout"] = t
	}
}
func WithIdleTimeout(t time.Duration) Option {
	return func(o *Redis) {
		o.redisOpt.IdleTimeout = t
		o.optsHash["IdleTimeout"] = t
	}
}
func WithPoolSize(i int) Option {
	return func(o *Redis) {
		o.redisOpt.PoolSize = i
		o.optsHash["PoolSize"] = i
	}
}
func WithPassword(s string) Option {
	return func(o *Redis) {
		o.redisOpt.Password = s
		o.optsHash["Password"] = s
	}
}
func WithDb(i int32) Option {
	return func(o *Redis) {
		o.redisOpt.DB = int(i)
		o.optsHash["DB"] = i
		o.redisLogger.Name += fmt.Sprintf("[%d]", i)
	}
}
func WithSlowTime(t time.Duration) Option {
	return func(o *Redis) {
		o.redisLogger.slowLogTime = t
	}
}
func WithLogger(l logger.ILogger) Option {
	return func(o *Redis) {
		o.redisLogger.logger = l
	}
}
func WithContext(c context.Context) Option {
	return func(o *Redis) {
		o.bootstrapContext = c
	}
}

// ========== /Option ==========

var (
	instanceLock sync.Mutex
	redisMap     = make(map[string]*Redis, 2)
)

type Redis struct {
	close            func()          `json:"-"`
	err              error           `json:"-"`
	redisLogger      *RedisLogger    `json:"-"`
	bootstrapContext context.Context `json:"-"`
	key              string          `json:"-"`

	client   *redis.Client          `json:"-"`
	redisOpt *redis.Options         `json:"-"`
	optsHash map[string]interface{} `json:"optsHash"`
}

func New(name string, addr string, opts ...Option) (rdb *Redis) {
	instanceLock.Lock()
	defer instanceLock.Unlock()

	rdb = &Redis{
		redisOpt: &redis.Options{
			Addr:        addr,
			PoolSize:    200,
			IdleTimeout: 240 * time.Second,
			ReadTimeout: 5 * time.Second,
			MaxRetries:  0,
		},
		bootstrapContext: context.Background(),
		optsHash:         make(map[string]interface{}, 2),
		redisLogger: &RedisLogger{
			Name:        name,
			logger:      logger.NewDefaultILogger(),
			slowLogTime: 10 * time.Second,
		},
		key: name + ":" + addr,
	}
	for _, o := range opts {
		o(rdb)
	}
	if b, e := json.Marshal(rdb.optsHash); e == nil {
		rdb.key += ":" + fmt.Sprintf("%+x", md5.Sum(b))
	}

	if r, ok := redisMap[rdb.key]; ok {
		rdb = r
		return
	}
	rdb.client = redis.NewClient(rdb.redisOpt)
	rdb.client.AddHook(rdb.redisLogger)

	if rdb.err = rdb.client.Ping(rdb.bootstrapContext).Err(); rdb.err != nil {
		rdb.LogError("New redis has err!")
		return
	}

	rdb.close = func() {
		if rdb.client == nil {
			rdb.LogWarn("Close redis is nil!")
			return
		}
		if rdb.err = rdb.client.Close(); rdb.err != nil {
			rdb.LogWarn("Close redis has error!")
			return
		}
	}

	redisMap[rdb.key] = rdb
	return
}
func (o *Redis) LogWarn(message string) {
	o.redisLogger.logger.Warn(
		o.bootstrapContext,
		message,
		database.KeyName, o.redisLogger.Name,
		database.KeyError, o.err,
	)
}
func (o *Redis) LogError(message string) {
	o.redisLogger.logger.Error(
		o.bootstrapContext,
		message,
		database.KeyName, o.redisLogger.Name,
		database.KeyError, o.err,
	)
}

func (o *Redis) Error() error {
	return o.err
}

func (o *Redis) Key() string {
	return o.key
}

func (o *Redis) Client() *redis.Client {
	return o.client
}

func (o *Redis) Close() func() {
	return o.close
}

type redisCtxBegintimeKey struct{}

type RedisLogger struct {
	Name        string         `json:"name"`
	slowLogTime time.Duration  `json:"slowTime"`
	logger      logger.ILogger `json:"-"`
}

func (h *RedisLogger) BeforeProcess(c context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(c, redisCtxBegintimeKey{}, time.Now()), nil
}

func (h *RedisLogger) AfterProcess(c context.Context, cmd redis.Cmder) (err error) {

	// slow
	begin := c.Value(redisCtxBegintimeKey{}).(time.Time)
	cost := time.Since(begin)

	p := []interface{}{
		database.KeyName, h.Name,
		database.KeyLimitTime, h.slowLogTime,
		database.KeyCostTime, cost.Seconds(),
	}

	if cost > h.slowLogTime {
		h.logger.Warn(c, database.FlagSlow+cmd.String(), p...)
		return
	}

	// error
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		p = append(p, database.KeyError, cmd.Err())
		h.logger.Error(c, cmd.String(), p...)
		return
	}

	// trace
	h.logger.Info(c, cmd.String(), p...)
	return
}

func (h *RedisLogger) BeforeProcessPipeline(c context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(c, redisCtxBegintimeKey{}, time.Now()), nil
}

func (h *RedisLogger) AfterProcessPipeline(c context.Context, cmds []redis.Cmder) (err error) {
	var b strings.Builder
	for _, s := range cmds {
		b.WriteString(s.String() + ",")
	}
	command := b.String()

	// slow
	begin := c.Value(redisCtxBegintimeKey{}).(time.Time)
	cost := time.Since(begin)

	p := []interface{}{
		database.KeyName, h.Name,
		database.KeyLimitTime, h.slowLogTime,
		database.KeyCostTime, cost.Seconds(),
	}

	if cost > h.slowLogTime {
		h.logger.Warn(c, database.FlagSlow+command, p...)
		return
	}

	// error
	for _, cmd := range cmds {
		if cmd.Err() != nil && cmd.Err() != redis.Nil {
			p = append(p, database.KeyError, err)
			h.logger.Error(c, cmd.String(), p...)
			return
		}
	}

	// trace
	h.logger.Info(c, command, p...)
	return
}
