package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(level string, service string) *Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: l,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts).WithAttrs([]slog.Attr{
		slog.String("service", service),
	})
	return &Logger{slog.New(handler)}
}

func (l *Logger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.Logger.Error(msg)
}

func (l *Logger) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug(msg)
}
