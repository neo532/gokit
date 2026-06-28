package slog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/neo532/gokit/logger"
)

var bufPool = sync.Pool{
	New: func() any { return &bytes.Buffer{} },
}

// textHandler is a custom text handler that writes log records
// using the specified separator between key=value pairs.
type textHandler struct {
	opts         *slog.HandlerOptions
	separator    string
	writer       io.Writer
	contextParam []logger.ContextArgs
	mu           sync.Mutex
	goa          []slog.Attr
}

func newTextHandler(writer io.Writer, opts *slog.HandlerOptions, separator string, contextParam []logger.ContextArgs) *textHandler {
	return &textHandler{
		opts:         opts,
		separator:    separator,
		writer:       writer,
		contextParam: contextParam,
	}
}

func (h *textHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.opts != nil && h.opts.Level != nil {
		return level >= h.opts.Level.Level()
	}
	return true
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	// time
	if !r.Time.IsZero() {
		buf.WriteString(r.Time.Format(time.RFC3339Nano))
		buf.WriteString(h.separator)
	}

	// level
	level := r.Level.String()
	if h.opts != nil && h.opts.ReplaceAttr != nil {
		if a := h.opts.ReplaceAttr(nil, slog.Any(slog.LevelKey, r.Level)); a.Key != "" {
			level = a.Value.String()
		}
	}
	buf.WriteString(level)
	buf.WriteString(h.separator)

	// message
	msg := r.Message
	if h.opts != nil && h.opts.ReplaceAttr != nil {
		if a := h.opts.ReplaceAttr(nil, slog.Any(slog.MessageKey, r.Message)); a.Key != "" {
			msg = a.Value.String()
		}
	}
	writeString(buf, msg)

	// attributes
	r.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		if h.opts != nil && h.opts.ReplaceAttr != nil {
			a = h.opts.ReplaceAttr(nil, a)
			if a.Equal(slog.Attr{}) {
				return true
			}
			a.Value = a.Value.Resolve()
		}
		buf.WriteString(h.separator)
		writeAttr(buf, a)
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	_, err := h.writer.Write(buf.Bytes())
	h.mu.Unlock()
	return err
}

func (h *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return &textHandler{
		opts:         h.opts,
		separator:    h.separator,
		writer:       h.writer,
		contextParam: h.contextParam,
		mu:           sync.Mutex{},
		goa:          append([]slog.Attr{}, h.goa...),
	}
}

func (h *textHandler) WithGroup(name string) slog.Handler {
	return h
}

func writeAttr(buf *bytes.Buffer, a slog.Attr) {
	if a.Value.Kind() == slog.KindGroup {
		for _, aa := range a.Value.Group() {
			writeAttr(buf, aa)
		}
		return
	}
	writeKey(buf, a.Key)
	buf.WriteByte('=')
	writeValue(buf, a.Value)
}

func writeKey(buf *bytes.Buffer, key string) {
	needsQuotes := false
	for _, r := range key {
		if r == '=' || r == '"' || r <= ' ' {
			needsQuotes = true
			break
		}
	}
	if needsQuotes {
		buf.WriteByte('"')
		buf.WriteString(key)
		buf.WriteByte('"')
		return
	}
	buf.WriteString(key)
}

func writeValue(buf *bytes.Buffer, v slog.Value) {
	switch v.Kind() {
	case slog.KindString:
		writeString(buf, v.String())
	case slog.KindInt64:
		buf.Write(strconv.AppendInt(nil, v.Int64(), 10))
	case slog.KindUint64:
		buf.Write(strconv.AppendUint(nil, v.Uint64(), 10))
	case slog.KindFloat64:
		buf.Write(strconv.AppendFloat(nil, v.Float64(), 'g', -1, 64))
	case slog.KindBool:
		buf.Write(strconv.AppendBool(nil, v.Bool()))
	case slog.KindTime:
		buf.WriteString(v.Time().Format(time.RFC3339Nano))
	case slog.KindDuration:
		buf.WriteString(v.Duration().String())
	case slog.KindAny:
		fmt.Fprint(buf, v.Any())
	}
}

func writeString(buf *bytes.Buffer, s string) {
	needsQuoting := false
	for _, r := range s {
		if r == '"' || r == '\\' || r == '=' || r <= ' ' {
			needsQuoting = true
			break
		}
	}
	if !needsQuoting && s != "" {
		buf.WriteString(s)
		return
	}
	buf.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 {
				writeHexByte(buf, byte(r))
			} else {
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteByte('"')
}

func writeHexByte(buf *bytes.Buffer, b byte) {
	const hex = "0123456789abcdef"
	buf.WriteString(`\x`)
	buf.WriteByte(hex[b>>4])
	buf.WriteByte(hex[b&0x0f])
}

// WithTextSeparator returns an Option that configures the text handler
// with a custom key=value separator (default is space, use "||" for pipe-separated).
func WithTextSeparator(sep string) Option {
	return WithHandler(textHandlerCreator{separator: sep})
}

type textHandlerCreator struct {
	separator string
}

func (t textHandlerCreator) NewSlogHandler(writer io.Writer, opts *slog.HandlerOptions, contextParam []logger.ContextArgs) slog.Handler {
	return newTextHandler(writer, opts, t.separator, contextParam)
}
