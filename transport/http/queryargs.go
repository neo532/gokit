package http

import (
	"context"
	"strings"
)

// QueryArgs is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type QueryArgs map[string][]string

// New creates an MD from a given key-values map.
func New(mds ...map[string][]string) QueryArgs {
	md := QueryArgs{}
	for _, m := range mds {
		for k, vList := range m {
			for _, v := range vList {
				md.Add(k, v)
			}
		}
	}
	return md
}

// Add adds the key, value pair to the header.
func (m QueryArgs) Add(key, value string) {
	if len(key) == 0 {
		return
	}

	m[strings.ToLower(key)] = append(m[strings.ToLower(key)], value)
}

// Get returns the value associated with the passed key.
func (m QueryArgs) Get(key string) string {
	v := m[strings.ToLower(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Set stores the key-value pair.
func (m QueryArgs) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	m[strings.ToLower(key)] = []string{value}
}

// Range iterate over element in queryArgs.
func (m QueryArgs) Range(f func(k string, v []string) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

// Values returns a slice of values associated with the passed key.
func (m QueryArgs) Values(key string) []string {
	return m[strings.ToLower(key)]
}

// Clone returns a deep copy of QueryArgs
func (m QueryArgs) Clone() QueryArgs {
	md := make(QueryArgs, len(m))
	for k, v := range m {
		md[k] = make([]string, len(v))
		copy(md[k], v)
	}
	return md
}

// AppendToUrl returns url with queryArgs
func (m QueryArgs) AppendToUrl(url string) string {

	m.Range(func(k string, v []string) (b bool) {

		pk := "{" + k + "}"
		var vv string
		if len(v) == 1 {
			vv = v[0]
		} else {
			vv = strings.Join(v, ",")
		}
		switch {
		case k == "":
		case strings.Index(url, pk) > -1:
			url = strings.Replace(url, pk, vv, -1)
		case strings.Index(url, "?") == -1:
			url += "?" + k + "=" + vv
		default:
			url += "&" + k + "=" + vv
		}

		return true
	})

	return url
}

type serverQueryArgsKey struct{}

// NewServerContext creates a new context with client md attached.
func NewServerContext(ctx context.Context, md QueryArgs) context.Context {
	return context.WithValue(ctx, serverQueryArgsKey{}, md)
}

// FromServerContext returns the server queryArgs in ctx if it exists.
func FromServerContext(ctx context.Context) (QueryArgs, bool) {
	md, ok := ctx.Value(serverQueryArgsKey{}).(QueryArgs)
	return md, ok
}

type clientQueryArgsKey struct{}

// NewClientContext creates a new context with client md attached.
func NewClientContext(ctx context.Context, md QueryArgs) context.Context {
	return context.WithValue(ctx, clientQueryArgsKey{}, md)
}

// FromClientContext returns the client queryArgs in ctx if it exists.
func FromClientContext(ctx context.Context) (QueryArgs, bool) {
	md, ok := ctx.Value(clientQueryArgsKey{}).(QueryArgs)
	return md, ok
}

// AppendToClientContext returns a new context with the provided kv merged
// with any existing queryArgs in the context.
func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		kv = append(kv, "")
	}
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return NewClientContext(ctx, md)
}

// MergeToClientContext merge new queryArgs into ctx.
func MergeToClientContext(ctx context.Context, cmd QueryArgs) context.Context {
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for k, v := range cmd {
		md[k] = v
	}
	return NewClientContext(ctx, md)
}
