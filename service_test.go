package shogoa

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// TErrorHandler is a test middleware that sets a witness to true if an error is returned.
func TErrorHandler(witness *bool) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			err := h(ctx, rw, req)
			if err != nil {
				*witness = true
			}
			return nil
		}
	}
}

// TMiddleware is a test middleware that sets a witness to true.
func TMiddleware(witness *bool) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			*witness = true
			return h(ctx, rw, req)
		}
	}
}

// CMiddleware is a test middleware that increments a witness.
func CMiddleware(witness *int) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			*witness++
			return h(ctx, rw, req)
		}
	}
}

// SecondMiddleware is a test middleware that sets a witness to true if the first witness is true.
func SecondMiddleware(witness1, witness2 *bool) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if !*witness1 {
				panic("middleware called in wrong order")
			}
			*witness2 = true
			return h(ctx, rw, req)
		}
	}
}

func TestService_New(t *testing.T) {
	const appName = "foo"
	s := New(appName)

	if s.Name != appName {
		t.Errorf("expected service name to be %s, got %s", appName, s.Name)
	}
	if s.Mux == nil {
		t.Errorf("expected service mux to be initialized, got nil")
	}
	if s.Server == nil {
		t.Errorf("expected service server to be initialized, got nil")
	}
}

func TestService_NotFound(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		req := httptest.NewRequest(http.MethodGet, "/404", nil)
		rw := httptest.NewRecorder()
		s.Mux.ServeHTTP(rw, req)

		result := rw.Result()
		body, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		ok, err := regexp.Match(`{"id":".*","code":"not_found","status":404,"detail":"/404"}`+"\n", body)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("expected 404 not found response, got %s", string(body))
		}
	})

	t.Run("with middleware", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		var called bool
		s.Use(TErrorHandler(&called))

		req := httptest.NewRequest(http.MethodGet, "/404", nil)
		rw := httptest.NewRecorder()
		s.Mux.ServeHTTP(rw, req)

		if !called {
			t.Error("expected error handler to be called, got false")
		}
	})

	t.Run("with middleware and multiple controllers", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		var calledCount int
		s.Use(CMiddleware(&calledCount))
		ctrl := s.NewController("test")
		ctrl.MuxHandler("/foo", nil, nil)
		ctrl.MuxHandler("/bar", nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		s.Mux.ServeHTTP(rw, req)

		if calledCount != 1 {
			t.Errorf("expected middleware to be called once, got %d", calledCount)
		}
	})
}

func TestService_MethodNotAllowed(t *testing.T) {
	s := New("foo")
	s.Decoder.Register(NewJSONDecoder, "*/*")
	s.Encoder.Register(NewJSONEncoder, "*/*")

	s.Mux.Handle(http.MethodPost, "/foo", func(rw http.ResponseWriter, req *http.Request, vals url.Values) {})
	s.Mux.Handle(http.MethodPut, "/foo", func(rw http.ResponseWriter, req *http.Request, vals url.Values) {})

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	s.Mux.ServeHTTP(rw, req)

	result := rw.Result()
	if result.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, result.StatusCode)
	}
	if got := result.Header.Get("Allow"); got != "POST, PUT" {
		t.Errorf("expected Allow header to be POST, PUT, got %s", got)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := regexp.Match(`{"id":".*","code":"method_not_allowed","status":405,"detail":".*","meta":{.*}}`+"\n", body)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("expected 404 not found response, got %s", string(body))
	}
}

func TestService_MaxRequestBodyLength(t *testing.T) {
	s := New("foo")
	s.Decoder.Register(NewJSONDecoder, "*/*")
	s.Encoder.Register(NewJSONEncoder, "*/*")

	ctrl := s.NewController("test")
	ctrl.MaxRequestBodyLength = 4
	unmarshaler := func(ctx context.Context, service *Service, req *http.Request) error {
		_, err := io.ReadAll(req.Body)
		return err
	}
	handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		rw.WriteHeader(http.StatusBadRequest)
		if _, err := rw.Write([]byte(ContextError(ctx).Error())); err != nil {
			return err
		}
		return nil
	}
	muxHandler := ctrl.MuxHandler("testMax", handler, unmarshaler)

	req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader(`"234"`))
	rw := httptest.NewRecorder()
	muxHandler(rw, req, nil)

	result := rw.Result()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := regexp.Match(`\[.*\] 413 request_too_large: request body length exceeds 4 bytes`, body)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("expected 404 not found response, got %s", string(body))
	}
}

func TestService_MuxHandler(t *testing.T) {
	unmarshaler := func(ctx context.Context, service *Service, req *http.Request) error {
		var payload any
		if err := service.DecodeRequest(req, &payload); err != nil {
			return err
		}
		ContextRequest(ctx).Payload = payload
		return nil
	}

	t.Run("it should not race", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		nopHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return nil
		}
		ctrl := s.NewController("test")
		ctrl.Use(func(h Handler) Handler {
			return nopHandler
		})
		for range 5 {
			s.Use(func(h Handler) Handler {
				return nopHandler
			})
		}

		var handlers []MuxHandler
		for range 10 {
			handler := ctrl.MuxHandler("test", nopHandler, nil)
			handlers = append(handlers, handler)
		}

		// Run all handlers concurrently.
		// It should not race.
		var wg sync.WaitGroup
		wg.Add(len(handlers))
		for _, h := range handlers {
			go func() {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodGet, "/foo", nil)
				rw := httptest.NewRecorder()
				h(rw, req, nil)
			}()
		}
		wg.Wait()
	})

	t.Run("request parameters", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := ContextError(ctx); err != nil {
				t.Error(err)
			}
			params := ContextRequest(ctx).Params
			if got := params.Get("id"); got != "42" {
				t.Errorf("expected id to be 42, got %s", got)
			}
			if got := params.Get("sort"); got != "asc" {
				t.Errorf("expected sort to be asc, got %s", got)
			}
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		param := url.Values{
			"id":   []string{"42"},
			"sort": []string{"asc"},
		}
		muxHandler(rw, req, param)
	})

	t.Run("invalid payload", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			err := ContextError(ctx)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "failed to decode") {
				t.Errorf("expected error to contain 'failed to decode', got %s", err.Error())
			}
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader("not json"))
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)
	})

	t.Run("with middleware", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		var called bool
		s.Use(TMiddleware(&called))

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)

		if !called {
			t.Error("expected the middleware to be called, got false")
		}
	})

	t.Run("middleware chain", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		var first, second bool
		s.Use(TMiddleware(&first))
		s.Use(SecondMiddleware(&first, &second))

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)

		if !first {
			t.Error("expected the first middleware to be called, got false")
		}
		if !second {
			t.Error("expected the second middleware to be called, got false")
		}
	})

	t.Run("error in the handler", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		var called bool
		s.Use(TErrorHandler(&called))

		// return an error in handler.
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return errors.New("boom")
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)

		if !called {
			t.Error("expected the error handler to be called, got false")
		}
	})

	t.Run("decode JSON", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := ContextError(ctx); err != nil {
				t.Error(err)
			}
			payload := ContextRequest(ctx).Payload.(map[string]any)
			if len(payload) != 1 {
				t.Errorf("expected payload to have 1 key, got %d", len(payload))
			}
			if got := payload["foo"]; got != "bar" {
				t.Errorf("expected foo to be bar, got %s", got)
			}
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader(`{"foo":"bar"}`))
		req.Header.Set("Content-Type", "application/json")
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)
	})

	t.Run("empty Content-Type", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := ContextError(ctx); err != nil {
				t.Error(err)
			}
			payload := ContextRequest(ctx).Payload.(map[string]any)
			if len(payload) != 1 {
				t.Errorf("expected payload to have 1 key, got %d", len(payload))
			}
			if got := payload["foo"]; got != "bar" {
				t.Errorf("expected foo to be bar, got %s", got)
			}
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader(`{"foo":"bar"}`))
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)
	})

	t.Run("fallback to the default decoder", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "*/*")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := ContextError(ctx); err != nil {
				t.Error(err)
			}
			payload := ContextRequest(ctx).Payload.(map[string]any)
			if len(payload) != 1 {
				t.Errorf("expected payload to have 1 key, got %d", len(payload))
			}
			if got := payload["foo"]; got != "bar" {
				t.Errorf("expected foo to be bar, got %s", got)
			}
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader(`{"foo":"bar"}`))
		req.Header.Set("Content-Type", "application/octet-stream")
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)
	})

	t.Run("bypass decoding", func(t *testing.T) {
		s := New("foo")
		s.Decoder.Register(NewJSONDecoder, "application/json")
		s.Encoder.Register(NewJSONEncoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := ContextError(ctx); err != nil {
				t.Error(err)
			}
			payload := ContextRequest(ctx).Payload
			if payload != nil {
				t.Errorf("expected payload to be nil, got %v", payload)
			}
			return nil
		}

		ctrl := s.NewController("test")
		muxHandler := ctrl.MuxHandler("testAct", handler, unmarshaler)

		req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader(`{"foo":"bar"}`))
		req.Header.Set("Content-Type", "application/octet-stream")
		rw := httptest.NewRecorder()
		muxHandler(rw, req, nil)
	})
}

func TestService_FileHandler(t *testing.T) {
	s := New("foo")
	s.Decoder.Register(NewJSONDecoder, "*/*")
	s.Encoder.Register(NewJSONEncoder, "*/*")

	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "swagger.json")
	err := os.WriteFile(filename, []byte(`{"foo":"bar"}`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	ctrl := s.NewController("test")
	handler := ctrl.FileHandler("/swagger.json", filename)
	muxHandler := ctrl.MuxHandler("testAct", handler, nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger.json", nil)
	rw := httptest.NewRecorder()
	muxHandler(rw, req, nil)

	result := rw.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(body); got != `{"foo":"bar"}` {
		t.Errorf("expected body to be {\"foo\":\"bar\"}, got %s", got)
	}
}
