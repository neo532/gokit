package slog

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
	//"golang.org/x/exp/slog"

	"github.com/neo532/gokit/logger"
)

type PrettyHandler struct {
	l            *log.Logger
	contextParam []logger.ContextArgs
	writer       io.Writer
	opts         *slog.HandlerOptions
}

func NewPrettyHandler() *PrettyHandler {
	return &PrettyHandler{}
}

func (h *PrettyHandler) NewSlogHandler(
	writer io.Writer,
	opts *slog.HandlerOptions,
	contextParam []logger.ContextArgs,
) slog.Handler {
	h.writer = writer
	h.opts = opts
	h.contextParam = contextParam
	h.l = log.New(writer, "", 0)
	return h
}

func (h *PrettyHandler) Handle(c context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	for _, fn := range h.contextParam {
		r.AddAttrs(slog.Any(fn(c)))
	}
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func (h *PrettyHandler) Enabled(c context.Context, l slog.Level) (b bool) {
	// if !h.Handler.Enabled(c, l) {
	// 	return
	// }
	return true
}

func (h *PrettyHandler) WithAttrs(as []slog.Attr) slog.Handler {
	return &PrettyHandler{l: h.l, contextParam: h.contextParam}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{l: h.l, contextParam: h.contextParam}
}
