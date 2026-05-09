package logger

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"default", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.level, "test")
			require.NotNil(t, l)
			require.NotNil(t, l.Logger)
		})
	}
}

func TestLogger_LogMethods(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(&buf, opts)
	l := &Logger{slog.New(handler)}

	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")

	output := buf.String()
	require.Contains(t, output, "debug msg")
	require.Contains(t, output, "info msg")
	require.Contains(t, output, "warn msg")
	require.Contains(t, output, "error msg")
}
