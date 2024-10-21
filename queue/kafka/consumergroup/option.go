package consumergroup

/*
 * @abstract consumer's option
 * @mail neo532@126.com
 * @date 2023-08-13
 */

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
	"github.com/neo532/gokit/queue/kafka"
)

type Option func(o *ConsumerGroup)

func WithAutoCommit(b bool) Option {
	return func(o *ConsumerGroup) {
		o.conf.Consumer.Offsets.AutoCommit.Enable = b
	}
}
func WithBalanceStrategy(strategy sarama.BalanceStrategy) Option {
	return func(o *ConsumerGroup) {
		o.conf.Consumer.Group.Rebalance.Strategy = strategy
	}
}
func WithVersion(ver sarama.KafkaVersion) Option {
	return func(o *ConsumerGroup) {
		o.conf.Version = ver
	}
}

func WithLogger(l logger.ILogger, c context.Context) Option {
	return func(o *ConsumerGroup) {
		o.handler.logger = l
		sarama.DebugLogger = kafka.NewLogger(l).WithContext(c)
	}
}
func WithSlowLog(t time.Duration) Option {
	return func(o *ConsumerGroup) {
		o.handler.slowTime = t
	}
}
func WithGoCount(count int) Option {
	return func(o *ConsumerGroup) {
		o.goCount = count
	}
}
func WithEnv(env string) Option {
	return func(o *ConsumerGroup) {
		o.handler.env = env
	}
}
func WithTopics(s ...string) Option {
	return func(o *ConsumerGroup) {
		o.topics = s
	}
}
func WithHandler(fn func(ctx context.Context, message []byte) (err error)) Option {
	return func(o *ConsumerGroup) {
		o.handler.handler = fn
	}
}
func WithContext(c context.Context) Option {
	return func(o *ConsumerGroup) {
		o.bootstrapContext = c
	}
}
func WithMiddleware(ms ...queue.ConsumerMiddleware) Option {
	return func(o *ConsumerGroup) {
		o.handler.middleware = append(o.handler.middleware, ms...)
	}
}
