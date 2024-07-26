package shogoa

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMux(t *testing.T) {
	t.Run("with no handler", func(t *testing.T) {
		mux := NewMux()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, req)
		if rw.Code != http.StatusNotFound {
			t.Errorf("unexpected status code: %d", rw.Code)
		}
	})

	t.Run("with registered handlers", func(t *testing.T) {
		mux := NewMux()
		mux.Handle(http.MethodPost, "/foo", func(w http.ResponseWriter, r *http.Request, v url.Values) {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if r.URL.Path != "/foo" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Errorf("unexpected method: %s", r.Method)
			}
			if string(data) != "some body" {
				t.Errorf("unexpected body: %s", string(data))
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader("some body"))
		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, req)
		if rw.Code != http.StatusOK {
			t.Errorf("unexpected status code: %d", rw.Code)
		}
	})

	t.Run("with registered handlers and wrong method", func(t *testing.T) {
		mux := NewMux()
		mux.Handle(http.MethodPost, "/foo", func(w http.ResponseWriter, r *http.Request, v url.Values) {})
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, req)
		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("unexpected status code: %d", rw.Code)
		}
	})
}
