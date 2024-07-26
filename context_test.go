package shogoa

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestResponseData_SwitchWriter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rw := httptest.NewRecorder()
	params := url.Values{"query": []string{"value"}}
	ctx := NewContext(rw, req, params)
	data := ContextResponse(ctx)

	rwo := data.SwitchWriter(httptest.NewRecorder())
	if rwo != rw {
		t.Errorf("unexpected response writer: want %p, got %p", rw, rwo)
	}
	if data.ResponseWriter == rw {
		t.Error("response writer not switched")
	}
}

func TestResponseData_Write(t *testing.T) {
	t.Run("should call WriteHeader(http.StatusOK) if WriteHeader has not yet been called", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rw := httptest.NewRecorder()
		params := url.Values{"query": []string{"value"}}
		ctx := NewContext(rw, req, params)
		data := ContextResponse(ctx)

		_, err := data.Write(nil)
		if err != nil {
			t.Fatal(err)
		}
		if data.Status != http.StatusOK {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, data.Status)
		}
	})

	t.Run("should not affect Status if WriteHeader has been called", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rw := httptest.NewRecorder()
		params := url.Values{"query": []string{"value"}}
		ctx := NewContext(rw, req, params)
		data := ContextResponse(ctx)

		data.WriteHeader(http.StatusBadRequest)
		if _, err := data.Write(nil); err != nil {
			t.Fatal(err)
		}
		if data.Status != http.StatusBadRequest {
			t.Errorf("unexpected status code: want %d, got %d", http.StatusBadRequest, data.Status)
		}
	})
}
