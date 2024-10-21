package queue

import (
	"context"
)

// ProducerHandler defines the handler invoked by Middleware.
type ProducerHandler func(c context.Context, message interface{}) error

// ProducerMiddleware is queue transport middleware.
type ProducerMiddleware func(ProducerHandler) ProducerHandler

// Chain returns a ProducerMiddleware that specifies the chained handler for endpoint.
func ChainProducer(m ...ProducerMiddleware) ProducerMiddleware {
	return func(next ProducerHandler) ProducerHandler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// ConsumerHandler defines the handler invoked by Middleware.
type ConsumerHandler func(c context.Context, message []byte) error

// ConsumerMiddleware is queue transport middleware.
type ConsumerMiddleware func(ConsumerHandler) ConsumerHandler

// Chain returns a ConsumerMiddleware that specifies the chained handler for endpoint.
func ChainConsumer(m ...ConsumerMiddleware) ConsumerMiddleware {
	return func(next ConsumerHandler) ConsumerHandler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
