package middleware

import (
	"context"
	"errors"
	"testing"
)

// 创建测试用的中间件
func createMiddleware(id int) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
			return next(ctx, req, reply)
		}
	}
}

func TestChain(t *testing.T) {
	// 创建一个计数器来跟踪中间件的执行顺序
	var executionOrder []int

	// 创建测试用的中间件
	createTrackingMiddleware := func(id int) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
				executionOrder = append(executionOrder, id)
				return next(ctx, req, reply)
			}
		}
	}

	// 创建测试用的处理器
	handler := func(ctx context.Context, req, reply interface{}) (context.Context, error) {
		executionOrder = append(executionOrder, 0) // 0 表示处理器
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
			// 重置执行顺序
			executionOrder = nil

			// 创建中间件链
			chain := Chain(tt.middlewares...)

			// 执行处理器
			ctx := context.Background()
			_, err := chain(handler)(ctx, nil, nil)

			// 检查执行顺序
			if len(executionOrder) != len(tt.expectedOrder) {
				t.Errorf("execution order length = %v, want %v", len(executionOrder), len(tt.expectedOrder))
				return
			}

			for i, v := range tt.expectedOrder {
				if executionOrder[i] != v {
					t.Errorf("execution order[%d] = %v, want %v", i, executionOrder[i], v)
				}
			}

			// 检查错误
			if err != tt.expectedError {
				t.Errorf("error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestMiddlewareErrorHandling(t *testing.T) {
	// 创建一个返回错误的中间件
	errorMiddleware := func(err error) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
				return nil, err
			}
		}
	}

	// 创建测试用的处理器
	handler := func(ctx context.Context, req, reply interface{}) (context.Context, error) {
		return ctx, nil
	}

	// 测试错误
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
	// 创建一个修改上下文的中间件
	contextMiddleware := func(key, value interface{}) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, req, reply interface{}) (context.Context, error) {
				newCtx := context.WithValue(ctx, key, value)
				return next(newCtx, req, reply)
			}
		}
	}

	// 创建测试用的处理器
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
