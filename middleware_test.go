package shogoa

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	// helper functions to verify that the context is being correctly passed.
	type contextKey string
	markContext := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextKey("key"), "value")
	}
	checkMarked := func(ctx context.Context, t *testing.T) {
		t.Helper()
		if ctx.Value(contextKey("key")) != "value" {
			t.Error("context not marked")
		}
	}

	t.Run("shogoa Middleware", func(t *testing.T) {
		// create a new middleware
		called := false
		myMiddleware := Middleware(func(h Handler) Handler {
			called = true
			return h
		})
		got, err := NewMiddleware(myMiddleware)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		got(nil)
		if !called {
			t.Fatal("middleware not called")
		}
	})

	t.Run("shogoa middleware func", func(t *testing.T) {
		// create a new middleware
		called := false
		myMiddleware := func(h Handler) Handler {
			called = true
			return h
		}
		got, err := NewMiddleware(myMiddleware)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		got(nil)
		if !called {
			t.Fatal("middleware not called")
		}
	})

	t.Run("using a shogoa handler", func(t *testing.T) {
		service := New("test")
		service.Encoder.Register(NewJSONEncoder, "*/*")

		// create a new middleware
		myHandler := Handler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			checkMarked(ctx, t)
			return service.Send(ctx, http.StatusOK, "ok")
		})
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := markContext(NewContext(rw, req, nil))
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error { return nil }
		if err := got(h)(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
		}
	})

	t.Run("using a shogoa handler func", func(t *testing.T) {
		service := New("test")
		service.Encoder.Register(NewJSONEncoder, "*/*")

		// create a new middleware
		myHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			checkMarked(ctx, t)
			return service.Send(ctx, http.StatusOK, "ok")
		}
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := markContext(NewContext(rw, req, nil))
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error { return nil }
		if err := got(h)(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
		}
	})

	t.Run("using a http middleware func", func(t *testing.T) {
		service := New("test")
		service.Encoder.Register(NewJSONEncoder, "*/*")

		// create a new middleware
		myHandler := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkMarked(r.Context(), t)
				h.ServeHTTP(w, r)
			})
		}
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := markContext(NewContext(rw, req, nil))
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			checkMarked(ctx, t)
			return service.Send(ctx, http.StatusOK, "ok")
		}
		if err := got(h)(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
		}
	})

	t.Run("using a http handler", func(t *testing.T) {
		service := New("test")
		service.Encoder.Register(NewJSONEncoder, "*/*")

		// create a new middleware
		myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			checkMarked(r.Context(), t)
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("ok")); err != nil {
				panic(err)
			}
		})
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := markContext(NewContext(rw, req, nil))
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			checkMarked(ctx, t)
			return nil
		}
		if err := got(h)(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
		}
	})

	t.Run("using a http handler func", func(t *testing.T) {
		service := New("test")
		service.Encoder.Register(NewJSONEncoder, "*/*")

		// create a new middleware
		myHandler := func(w http.ResponseWriter, r *http.Request) {
			checkMarked(r.Context(), t)
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("ok")); err != nil {
				panic(err)
			}
		}
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		// verify
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := markContext(NewContext(rw, req, nil))
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			checkMarked(ctx, t)
			return nil
		}
		if err := got(h)(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
		}
	})
}
