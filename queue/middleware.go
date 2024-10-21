package queue

import (
	"context"
)

// Handler defines the handler invoked by Middleware.
type Handler func(c context.Context, message interface{}) error

// Middleware is queue transport middleware.
type Middleware func(Handler) Handler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
