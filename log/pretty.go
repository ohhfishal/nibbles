package log

import (
	"context"
	"fmt"
	"github.com/ohhfishal/nibbles/color"
	"io"
	"log/slog"
	"strings"
	"time"
)

var _ slog.Handler = &PrettyHandler{}

const TimeFormat = time.DateTime

var (
	debug = color.String(color.Purple, "[DEBUG]")
	info  = color.String(color.Green, "[INFO] ")
	warn  = color.String(color.Yellow, "[WARN] ")
	_err  = color.String(color.Red, "[ERROR]")
)

var defaultColors = []string{
	debug, debug, debug, debug,
	info, info, info, info,
	warn, warn, warn, warn,
	_err,
}

type PrettyHandler struct {
	stdout io.Writer
	opts   slog.HandlerOptions
	attrs  []slog.Attr
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	options := slog.HandlerOptions{}
	if opts != nil {
		options = *opts
	}
	return &PrettyHandler{
		stdout: w,
		opts:   options,
	}
}

func (handler *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if handler.opts.Level != nil {
		minLevel = handler.opts.Level.Level()
	}
	return level >= minLevel
}
func (handler *PrettyHandler) Handle(ctx context.Context, record slog.Record) error {
	if _, err := fmt.Fprint(handler.stdout, strings.Join([]string{
		record.Time.Format(TimeFormat),
		handler.format(record.Level),
		record.Message,
	}, " ",
	)); err != nil {
		return err
	}

	record.AddAttrs(handler.attrs...)

	record.Attrs(func(attr slog.Attr) bool {
		_, err := fmt.Fprintf(handler.stdout, " %s=%s", attr.Key, attr.Value.String())
		return err == nil
	})

	if _, err := fmt.Fprintln(handler.stdout); err != nil {
		return err
	}
	return nil
}

func (handler PrettyHandler) clone() *PrettyHandler {
	return &PrettyHandler{
		stdout: handler.stdout,
		opts:   handler.opts,
		attrs:  handler.attrs,
	}
}

func (handler *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return handler
	}
	h := handler.clone()
	h.attrs = append(h.attrs, attrs...)
	return h
}

func (handler *PrettyHandler) WithGroup(name string) slog.Handler {
	panic("not implemented: groups")
	// return handler.textHandler.WithGroup(name)
}

func (handler *PrettyHandler) format(level slog.Level) string {
	return defaultColors[max(0, min(len(defaultColors)-1, (int)(level)+4))]
}
