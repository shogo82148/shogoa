package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestRecover(t *testing.T) {
	t.Run("panics with a string", func(t *testing.T) {
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("boom")
		}
		rg := Recover()(h)
		err := rg(nil, nil, nil)
		if !strings.HasPrefix(err.Error(), "panic: boom\n") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("panics with an error", func(t *testing.T) {
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic(errors.New("boom"))
		}
		rg := Recover()(h)
		err := rg(nil, nil, nil)
		if !strings.HasPrefix(err.Error(), "panic: boom\n") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("panics with something else", func(t *testing.T) {
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic(42)
		}
		rg := Recover()(h)
		err := rg(nil, nil, nil)
		if !strings.HasPrefix(err.Error(), "unknown panic\n") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
