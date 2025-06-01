package middleware

import (
	"context"
	"errors"
	"testing"
)

// Create test middleware
func createMiddleware(id int) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
			return next(ctx, req, reply)
		}
	}
}

func TestChain(t *testing.T) {
	// Create a counter to track middleware execution order
	var executionOrder []int

	// Create test middleware
	createTrackingMiddleware := func(id int) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
				executionOrder = append(executionOrder, id)
				return next(ctx, req, reply)
			}
		}
	}

	// Create test handler
	handler := func(ctx context.Context, req, reply interface{}) (context.Context, error) {
		executionOrder = append(executionOrder, 0) // 0 represents the handler
		return ctx, nil
	}

	tests := []struct {
		name          string
		middlewares   []Middleware
		expectedOrder []int
		expectedError error
	}{
		{
			name:          "no middleware",
			middlewares:   []Middleware{},
			expectedOrder: []int{0},
		},
		{
			name:          "single middleware",
			middlewares:   []Middleware{createTrackingMiddleware(1)},
			expectedOrder: []int{1, 0},
		},
		{
			name:          "multiple middleware",
			middlewares:   []Middleware{createTrackingMiddleware(1), createTrackingMiddleware(2), createTrackingMiddleware(3)},
			expectedOrder: []int{1, 2, 3, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset execution order
			executionOrder = nil

			// Create middleware chain
			chain := Chain(tt.middlewares...)

			// Execute handler
			ctx := context.Background()
			_, err := chain(handler)(ctx, nil, nil)

			// Check execution order
			if len(executionOrder) != len(tt.expectedOrder) {
				t.Errorf("execution order length = %v, want %v", len(executionOrder), len(tt.expectedOrder))
				return
			}

			for i, v := range tt.expectedOrder {
				if executionOrder[i] != v {
					t.Errorf("execution order[%d] = %v, want %v", i, executionOrder[i], v)
				}
			}

			// Check error
			if err != tt.expectedError {
				t.Errorf("error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestMiddlewareErrorHandling(t *testing.T) {
	// Create middleware that returns an error
	errorMiddleware := func(err error) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
				return nil, err
			}
		}
	}

	// Create test handler
	handler := func(ctx context.Context, req, reply interface{}) (context.Context, error) {
		return ctx, nil
	}

	// Test error
	testError := errors.New("test error")

	tests := []struct {
		name          string
		middlewares   []Middleware
		expectedError error
	}{
		{
			name:          "middleware returns error",
			middlewares:   []Middleware{errorMiddleware(testError)},
			expectedError: testError,
		},
		{
			name:          "error in chain",
			middlewares:   []Middleware{errorMiddleware(testError)},
			expectedError: testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := Chain(tt.middlewares...)
			_, err := chain(handler)(context.Background(), nil, nil)

			if err != tt.expectedError {
				t.Errorf("error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestMiddlewareContextHandling(t *testing.T) {
	// Create middleware that modifies context
	contextMiddleware := func(key, value interface{}) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
				newCtx := context.WithValue(ctx, key, value)
				return next(newCtx, req, reply)
			}
		}
	}

	// Create test handler
	handler := func(ctx context.Context, req, reply interface{}) (context.Context, error) {
		return ctx, nil
	}

	tests := []struct {
		name           string
		middlewares    []Middleware
		expectedValues map[interface{}]interface{}
	}{
		{
			name: "single context modification",
			middlewares: []Middleware{
				contextMiddleware("key1", "value1"),
			},
			expectedValues: map[interface{}]interface{}{
				"key1": "value1",
			},
		},
		{
			name: "multiple context modifications",
			middlewares: []Middleware{
				contextMiddleware("key1", "value1"),
				contextMiddleware("key2", "value2"),
			},
			expectedValues: map[interface{}]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := Chain(tt.middlewares...)
			ctx, err := chain(handler)(context.Background(), nil, nil)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			for key, expectedValue := range tt.expectedValues {
				if value := ctx.Value(key); value != expectedValue {
					t.Errorf("context value for key %v = %v, want %v", key, value, expectedValue)
				}
			}
		})
	}
}
