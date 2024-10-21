package message

/*
 * @abstract delimiter message
 * @mail neo532@126.com
 * @date 2023-08-13
 */

import "strings"

type Format struct {
	delimiter string
}

func NewFormat() *Delimiter {
	return &Format{}
}

func (f *Format) Delimiter(delimiter string) *Format {
	f.delimiter = delimiter
	return f
}

func (f *Format) Fmt(key, msg string) (b []byte) {
	return []byte(key + f.delimiter + string(msg))
}

// !!! key must be half-angle
func (f *Format) Parse(msg []byte) (key string, value []byte) {
	if f.delimiter == "" {
		return "", msg
	}

	s := string(msg)
	if i := strings.Index(s, f.delimiter); i > 0 {
		key = s[:i]
		value = []byte(string([]rune(s)[i+1:]))
		return
	}

	return key, msg
}
