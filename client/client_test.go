package client_test

import (
	"context"
	"testing"

	"github.com/shogo82148/shogoa/client"
)

func TestContextRequestID(t *testing.T) {
	ctx := context.Background()
	if reqID := client.ContextRequestID(ctx); reqID != "" {
		t.Errorf("expected empty request ID, got %s", reqID)
	}
}

func TestContextWithRequestID(t *testing.T) {
	ctx := context.Background()
	newCtx, reqID := client.ContextWithRequestID(ctx)
	if reqID == "" {
		t.Errorf("expected non-empty request ID")
	}
	if got := client.ContextRequestID(newCtx); got != reqID {
		t.Errorf("expected request ID %s, got %s", reqID, got)
	}
}

func TestSetContextRequestID(t *testing.T) {
	ctx := context.Background()
	const customID = "foo"
	newCtx := client.SetContextRequestID(ctx, customID)
	if newCtx == ctx {
		t.Errorf("expected new context, got the same context")
	}
	if got := client.ContextRequestID(newCtx); got != customID {
		t.Errorf("expected request ID %s, got %s", customID, got)
	}

	// request ID should not need to be generated again. the same context should be returned instead.
	newCtx2, reqID := client.ContextWithRequestID(newCtx)
	if newCtx != newCtx2 {
		t.Errorf("expected the same context, got a different one")
	}
	if reqID != customID {
		t.Errorf("expected request ID %s, got %s", customID, reqID)
	}
}
