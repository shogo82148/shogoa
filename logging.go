package shogoa

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// ContextLogAdapter is the logger interface used by shogoa to log informational, warning and error messages.
// It allows to pass a context.Context to the logger.
type LogAdapter interface {
	// New appends to the logger context and returns the updated logger logger.
	New(keyvals ...any) LogAdapter
	// InfoContext logs an informational message.
	InfoContext(ctx context.Context, msg string, keyvals ...any)
	// ErrorContext logs an error.
	ErrorContext(ctx context.Context, msg string, keyvals ...any)
	// WarnContext logs a warning message.
	WarnContext(ctx context.Context, mgs string, keyvals ...any)
}

var _ LogAdapter = (*adapter)(nil)

// adapter is the slog shogoa logger adapter.
type adapter struct {
	handler slog.Handler
}

// New wraps a [log/slog.Handler] into a shogoa logger.
func NewLogger(handler slog.Handler) LogAdapter {
	return &adapter{handler: handler}
}

// Info logs messages using [log/slog].
func (a *adapter) Info(msg string, data ...any) {
	a.log(context.Background(), slog.LevelInfo, msg, data...)
}

// InfoContext logs messages using [log/slog].
func (a *adapter) InfoContext(ctx context.Context, msg string, data ...any) {
	a.log(ctx, slog.LevelInfo, msg, data...)
}

// Warn logs message using [log/slog].
func (a *adapter) Warn(msg string, data ...any) {
	a.log(context.Background(), slog.LevelWarn, msg, data...)
}

// WarnContext logs message using [log/slog].
func (a *adapter) WarnContext(ctx context.Context, msg string, data ...any) {
	a.log(ctx, slog.LevelWarn, msg, data...)
}

// Error logs errors using [log/slog].
func (a *adapter) Error(msg string, data ...any) {
	a.log(context.Background(), slog.LevelError, msg, data...)
}

// ErrorContext logs errors using [log/slog].
func (a *adapter) ErrorContext(ctx context.Context, msg string, data ...any) {
	a.log(ctx, slog.LevelError, msg, data...)
}

// New creates a new logger given a context.
func (a *adapter) New(data ...any) LogAdapter {
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "", 0)
	r.Add(data...)

	attrs := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	h := a.handler.WithAttrs(attrs)
	return &adapter{handler: h}
}

func (a *adapter) log(ctx context.Context, level slog.Level, msg string, data ...any) {
	if !a.handler.Enabled(ctx, level) {
		return
	}

	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this functions, this functions caller, the caller of the adapter]
	runtime.Callers(4, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(data...)
	_ = a.handler.Handle(ctx, r)
}

// LogInfo extracts the logger from the given context and calls Info on it.
// This is intended for code that needs portable logging such as the internal code of shogoa and
// middleware. User code should use the log adapters instead.
func LogInfo(ctx context.Context, msg string, keyvals ...any) {
	// This block should be synced with Service.LogInfo
	if l := ctx.Value(logKey); l != nil {
		switch logger := l.(type) {
		case LogAdapter:
			logger.InfoContext(ctx, msg, keyvals...)
		}
	}
}

// LogWarn extracts the logger from the given context and calls Warn on it.
// This is intended for code that needs portable logging such as the internal code of shogoa and
// middleware. User code should use the log adapters instead.
func LogWarn(ctx context.Context, msg string, keyvals ...any) {
	if l := ctx.Value(logKey); l != nil {
		switch logger := l.(type) {
		case LogAdapter:
			logger.WarnContext(ctx, msg, keyvals...)
		}
	}
}

// LogError extracts the logger from the given context and calls Error on it.
// This is intended for code that needs portable logging such as the internal code of shogoa and
// middleware. User code should use the log adapters instead.
func LogError(ctx context.Context, msg string, keyvals ...any) {
	// this block should be synced with Service.LogError
	if l := ctx.Value(logKey); l != nil {
		switch logger := l.(type) {
		case LogAdapter:
			logger.ErrorContext(ctx, msg, keyvals...)
		}
	}
}
