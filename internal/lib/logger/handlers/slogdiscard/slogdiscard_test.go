package slogdiscard

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewDiscardLogger(t *testing.T) {
	logger := NewDiscardLogger()
	if logger == nil {
		t.Fatal("NewDiscardLogger returned nil")
	}

	// Test that logging doesn't panic and returns no error
	ctx := context.Background()
	logger.InfoContext(ctx, "test message", slog.String("key", "value"))
	logger.DebugContext(ctx, "debug message")
	logger.WarnContext(ctx, "warn message")
	logger.ErrorContext(ctx, "error message")
}

func TestDiscardHandler(t *testing.T) {
	handler := NewDiscardHandler()
	if handler == nil {
		t.Fatal("NewDiscardHandler returned nil")
	}

	ctx := context.Background()

	// Test Enabled - should always return false
	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("Enabled should return false for discard handler")
	}
	if handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("Enabled should return false for discard handler")
	}

	// Test Handle - should return nil without doing anything
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	if err := handler.Handle(ctx, record); err != nil {
		t.Errorf("Handle should return nil, got: %v", err)
	}

	// Test WithAttrs - should return a handler
	newHandler := handler.WithAttrs([]slog.Attr{slog.String("key", "value")})
	if newHandler == nil {
		t.Error("WithAttrs should return a non-nil handler")
	}

	// Test WithGroup - should return a handler
	groupHandler := handler.WithGroup("group")
	if groupHandler == nil {
		t.Error("WithGroup should return a non-nil handler")
	}
}
