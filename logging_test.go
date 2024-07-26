package shogoa

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

var optsRemoveTime = &slog.HandlerOptions{
	ReplaceAttr: removeTime,
}

// removeTime removes time attribute for stable testing.
func removeTime(groups []string, a slog.Attr) slog.Attr {
	// Remove time.
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

func TestLogInfo(t *testing.T) {
	LogInfo(context.Background(), "LogInfo with a nil log doesn't crash")

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, optsRemoveTime)
	ctx := WithLogger(context.Background(), NewLogger(handler))
	LogInfo(ctx, "message", "foo", "bar")
	if buf.String() != "level=INFO msg=message foo=bar\n" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestLogWarn(t *testing.T) {
	LogWarn(context.Background(), "LogWarn with a nil log doesn't crash")

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, optsRemoveTime)
	ctx := WithLogger(context.Background(), NewLogger(handler))
	LogWarn(ctx, "message", "foo", "bar")
	if buf.String() != "level=WARN msg=message foo=bar\n" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestLogError(t *testing.T) {
	LogError(context.Background(), "LogError with a nil log doesn't crash")

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, optsRemoveTime)
	ctx := WithLogger(context.Background(), NewLogger(handler))
	LogError(ctx, "message", "foo", "bar")
	if buf.String() != "level=ERROR msg=message foo=bar\n" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
