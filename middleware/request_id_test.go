package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shogo82148/shogoa"
)

func TestRequestID(t *testing.T) {
	t.Run("sets the request ID in the context", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			id := ContextRequestID(ctx)
			if id != "my-request-id" {
				t.Errorf("unexpected request id: %s", id)
			}
			return service.Send(ctx, http.StatusOK, "ok")
		}
		handler = RequestID()(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		req.Header.Set(RequestIDHeader, "my-request-id")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("truncates request ID when it exceeds a default limit", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			id := ContextRequestID(ctx)
			if id != strings.Repeat("x", DefaultRequestIDLengthLimit) {
				t.Errorf("unexpected request id: %s", id)
			}
			return service.Send(ctx, http.StatusOK, "ok")
		}
		handler = RequestID()(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		req.Header.Set(RequestIDHeader, strings.Repeat("x", 2*DefaultRequestIDLengthLimit))
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("truncates request ID when it exceeds a custom limit", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			id := ContextRequestID(ctx)
			if id != "1234567" {
				t.Errorf("unexpected request id: %s", id)
			}
			return service.Send(ctx, http.StatusOK, "ok")
		}
		handler = RequestIDWithHeaderAndLengthLimit("Foo", 7)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		req.Header.Set("Foo", "12345678")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("allows any request ID when length limit is negative", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			id := ContextRequestID(ctx)
			if id != strings.Repeat("x", 2*DefaultRequestIDLengthLimit) {
				t.Errorf("unexpected request id: %s", id)
			}
			return service.Send(ctx, http.StatusOK, "ok")
		}
		handler = RequestIDWithHeaderAndLengthLimit(RequestIDHeader, -1)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		req.Header.Set(RequestIDHeader, strings.Repeat("x", 2*DefaultRequestIDLengthLimit))
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})
}
