package slicex

func OfType[T any](s []any) []T {
	r := make([]T, 0, len(s))
	for _, v := range s {
		if vt, ok := v.(T); ok {
			r = append(r, vt)
		}
	}
	return r
}

func OfTypeAny[T any](s []T) []any {
	r := make([]any, 0, len(s))
	for _, v := range s {
		r = append(r, v)
	}
	return r
}

// Range
func Range[T any](ss []T, fn func(T)) {
	for _, s := range ss {
		fn(s)
	}
}
