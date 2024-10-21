package queue

import (
	"context"
	"fmt"
	"strings"
)

type Header map[string]string
type headerKey struct{}

// New creates an Header from a given key-values map.
func InitHeaderToContext(c context.Context) context.Context {
	return context.WithValue(c, headerKey{}, Header{})
}

// GetHeaderFromContext returns the header in ctx if it exists.
func GetHeaderFromContext(ctx context.Context) (h Header, b bool) {
	h, b = ctx.Value(headerKey{}).(Header)
	return
}

func AppendHeaderToContext(c context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("header: AppendHeaderToContext got an odd number of input pairs for header: %d", len(kv)))
	}

	if h, ok := GetHeaderFromContext(c); ok {
		for i := 0; i < len(kv); i += 2 {
			h.Set(kv[i], kv[i+1])
		}
	}
	return c
}

// Values returns a slice of values associated with the passed key.
func (h Header) Value(key string) string {
	return h[strings.ToLower(key)]
}

// Set stores the key-value pair.
func (h Header) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	h[strings.ToLower(key)] = value
}

// Range iterate over element in header.
func (h Header) Range(f func(k string, v string) bool) {
	for k, v := range h {
		if !f(k, v) {
			break
		}
	}
}
