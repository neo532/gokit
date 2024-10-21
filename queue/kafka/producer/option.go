package producer

/*
 * @abstract producer
 * @mail neo532@126.com
 * @date 2023-08-13
 */

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
	"github.com/neo532/gokit/queue/kafka"
)

// ========== Option ==========
type Option func(*Producer)

func WithVersion(ver sarama.KafkaVersion) Option {
	return func(o *Producer) {
		o.Conf.Version = ver
	}
}

func WithLogger(l logger.ILogger, c context.Context) Option {
	return func(o *Producer) {
		o.logger = l
		sarama.DebugLogger = kafka.NewLogger(l).WithContext(c)
	}
}

func WithAsync(b bool) Option {
	return func(o *Producer) {
		o.IsAsync = b
	}
}

func WithRequiredAcks(r sarama.RequiredAcks) Option {
	return func(o *Producer) {
		o.Conf.Producer.RequiredAcks = r
	}
}

func WithReturnSucesses(b bool) Option {
	return func(o *Producer) {
		o.Conf.Producer.Return.Successes = b
	}
}

// sarama.NewHashPartitioner
func WithPartitioner(fn sarama.PartitionerConstructor) Option {
	return func(o *Producer) {
		o.Conf.Producer.Partitioner = fn
	}
}

func WithTopic(topic string) Option {
	return func(o *Producer) {
		o.Topic = topic
	}
}

func WithContext(c context.Context) Option {
	return func(o *Producer) {
		o.bootstrapContext = c
	}
}

func WithIdempotent(b bool) Option {
	return func(o *Producer) {
		o.Conf.Producer.Idempotent = b
	}
}

func WithNetMaxOpenRequest(i int) Option {
	return func(o *Producer) {
		o.Conf.Net.MaxOpenRequests = i
	}
}

func WithMiddleware(ms ...queue.Middleware) Option {
	return func(o *Producer) {
		o.middleware = append(o.middleware, ms...)
	}
}
