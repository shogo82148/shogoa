package shogoa

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	t.Run("shogoa Middleware", func(t *testing.T) {
		called := false
		myMiddleware := Middleware(func(h Handler) Handler {
			called = true
			return h
		})
		got, err := NewMiddleware(myMiddleware)
		if err != nil {
			t.Fatal(err)
		}
		got(nil)
		if !called {
			t.Fatal("middleware not called")
		}
	})

	t.Run("shogoa middleware func", func(t *testing.T) {
		called := false
		myMiddleware := func(h Handler) Handler {
			called = true
			return h
		}
		got, err := NewMiddleware(myMiddleware)
		if err != nil {
			t.Fatal(err)
		}
		got(nil)
		if !called {
			t.Fatal("middleware not called")
		}
	})

	t.Run("using a shogoa handler", func(t *testing.T) {
		service := New("test")
		service.Encoder.Register(NewJSONEncoder, "*/*")

		myHandler := Handler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return service.Send(ctx, http.StatusOK, "ok")
		})
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := NewContext(rw, req, nil)
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

		myHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return service.Send(ctx, http.StatusOK, "ok")
		}
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := NewContext(rw, req, nil)
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

		myHandler := func(h http.Handler) http.Handler { return h }
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := NewContext(rw, req, nil)
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
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

		myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("ok")); err != nil {
				panic(err)
			}
		})
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := NewContext(rw, req, nil)
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
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

		myHandler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("ok")); err != nil {
				panic(err)
			}
		}
		got, err := NewMiddleware(myHandler)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := NewContext(rw, req, nil)
		h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
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
