package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/shogo82148/shogoa"
)

func TestRequiredHeader(t *testing.T) {
	t.Run("matches a header value", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return service.Send(ctx, http.StatusOK, "ok")
		}
		handler = RequireHeader(
			service,
			regexp.MustCompile("^/foo"),
			"Some-Header",
			regexp.MustCompile("^some value$"),
			http.StatusUnauthorized,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo/bar", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		req.Header.Set("Some-Header", "some value")
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		if shogoa.ContextResponse(ctx).Status != http.StatusOK {
			t.Errorf("unexpected status: %d", shogoa.ContextResponse(ctx).Status)
		}
	})

	t.Run("responds with failure on mismatch", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("unreachable")
		}
		handler = RequireHeader(
			service,
			regexp.MustCompile("^/foo"),
			"Some-Header",
			regexp.MustCompile("^some value$"),
			http.StatusUnauthorized,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo/bar", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		req.Header.Set("Some-Header", "some other value")
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		if shogoa.ContextResponse(ctx).Status != http.StatusUnauthorized {
			t.Errorf("unexpected status: %d", shogoa.ContextResponse(ctx).Status)
		}
	})

	t.Run("responds with failure when header is missing", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("unreachable")
		}
		handler = RequireHeader(
			service,
			regexp.MustCompile("^/foo"),
			"Some-Header",
			regexp.MustCompile("^some value$"),
			http.StatusUnauthorized,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo/bar", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		if shogoa.ContextResponse(ctx).Status != http.StatusUnauthorized {
			t.Errorf("unexpected status: %d", shogoa.ContextResponse(ctx).Status)
		}
	})

	t.Run("passes through for a non-matching path", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return service.Send(ctx, http.StatusOK, "ok")
		}
		handler = RequireHeader(
			service,
			regexp.MustCompile("^/baz"),
			"Some-Header",
			regexp.MustCompile("^some value$"),
			http.StatusUnauthorized,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo/bar", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		if shogoa.ContextResponse(ctx).Status != http.StatusOK {
			t.Errorf("unexpected status: %d", shogoa.ContextResponse(ctx).Status)
		}
	})

	t.Run("matches value for a nil path pattern", func(t *testing.T) {
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("unreachable")
		}
		handler = RequireHeader(
			service,
			nil,
			"Some-Header",
			regexp.MustCompile("^some value$"),
			http.StatusUnauthorized,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/foo/bar", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
		if shogoa.ContextResponse(ctx).Status != http.StatusUnauthorized {
			t.Errorf("unexpected status: %d", shogoa.ContextResponse(ctx).Status)
		}
	})
}
