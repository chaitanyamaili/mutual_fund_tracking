package logger

import (
	"io"
	"log/slog"
)

type Handler struct {
	writer    io.Writer
	level     slog.Level
	addSource bool
}

func NewHandler(w io.Writer, addSource bool) *Handler {
	return &Handler{
		writer:    w,
		addSource: addSource,
	}
}

func WithFormatter(w io.Writer, addSource bool, format string) *slog.Logger {
	var log *slog.Logger
	l := NewHandler(w, addSource)

	switch format {
	case "json":
		log = l.JSONLogger()
	default:
		log = l.TextLogger()
	}

	return log
}

func (h *Handler) TextLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:       h.level,
		AddSource:   h.addSource,
		ReplaceAttr: nil,
	}

	return slog.New(
		slog.NewTextHandler(
			h.writer,
			opts,
		),
	)
}

func (h *Handler) JSONLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:       h.level,
		AddSource:   h.addSource,
		ReplaceAttr: nil,
	}

	return slog.New(
		slog.NewJSONHandler(
			h.writer,
			opts,
		),
	)
}
