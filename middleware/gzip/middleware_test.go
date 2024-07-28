package gzip

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shogo82148/shogoa"
)

func TestGzip(t *testing.T) {
	verifyGzippedResponse := func(t *testing.T, resp *http.Response, body string) {
		t.Helper()
		if resp.Header.Get("Content-Encoding") != "gzip" {
			t.Errorf("unexpected Content-Encoding: %s", resp.Header.Get("Content-Encoding"))
		}

		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer gr.Close()
		actual, err := io.ReadAll(gr)
		if err != nil {
			t.Fatal(err)
		}
		if string(actual) != body {
			t.Errorf("unexpected body: %s", actual)
		}
	}

	t.Run("encodes response using gzip", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyGzippedResponse(t, result, "gzip me!")
	})

	t.Run("encodes response using gzip (custom status)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusBadRequest)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0), AddStatusCodes(http.StatusBadRequest))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyGzippedResponse(t, result, "gzip me!")
	})

	t.Run("encodes response using gzip (all status)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusBadRequest)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0), OnlyStatusCodes())(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyGzippedResponse(t, result, "gzip me!")
	})

	t.Run("encodes response using gzip (custom type)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.Header().Set("Content-Type", "custom/type")
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0), AddContentTypes("custom/type"))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyGzippedResponse(t, result, "gzip me!")
	})

	t.Run("encodes response using gzip (length check)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusOK)
			// Use multiple writes.
			for range 128 {
				if _, err := resp.Write([]byte("gzip me!")); err != nil {
					return err
				}
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(512))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyGzippedResponse(t, result, strings.Repeat("gzip me!", 128))
	})

	t.Run("removes Accept-Ranges header", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.Header().Set("Accept-Ranges", "some value")
			resp.WriteHeader(http.StatusOK)
			// Use multiple writes.
			for range 128 {
				if _, err := resp.Write([]byte("gzip me!")); err != nil {
					return err
				}
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(512))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Accept-Ranges") != "" {
			t.Errorf("unexpected Accept-Ranges header: %s", result.Header.Get("Accept-Ranges"))
		}
		verifyGzippedResponse(t, result, strings.Repeat("gzip me!", 128))
	})

	t.Run("should preserve status code", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusConflict)
			// Use multiple writes.
			for range 128 {
				if _, err := resp.Write([]byte("gzip me!")); err != nil {
					return err
				}
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(512), AddStatusCodes(http.StatusConflict))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusConflict {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Accept-Ranges") != "" {
			t.Errorf("unexpected Accept-Ranges header: %s", result.Header.Get("Accept-Ranges"))
		}
		verifyGzippedResponse(t, result, strings.Repeat("gzip me!", 128))
	})
}

func TestNotGzip(t *testing.T) {
	verifyNotGzippedResponse := func(t *testing.T, resp *http.Response, body string) {
		t.Helper()

		if resp.Header.Get("Content-Encoding") == "gzip" {
			t.Errorf("unexpected Content-Encoding: %s", resp.Header.Get("Content-Encoding"))
		}
		actual, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(actual) != body {
			t.Errorf("unexpected body: %s", actual)
		}
	}

	t.Run("does not encode response (already gzipped)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.Header().Set("Content-Type", "gzip")
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip data")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip data")
	})

	t.Run("does not encode response (too small)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression)(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("does not encode response (wrong status code)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusBadRequest)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("does not encode response (removed status code)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0), OnlyStatusCodes(http.StatusBadRequest))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("does not encode response (unknown content type)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.Header().Add("Content-Type", "unknown/content-type")
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("does not encode response (removed type)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0), OnlyContentTypes("some/type"))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("does not encode response (has Range header)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusOK)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0), IgnoreRange(false))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("should preserve status code", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusConflict)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression)(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusConflict {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("should preserve status code with no body", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			resp.WriteHeader(http.StatusConflict)
			return nil
		}
		handler = Middleware(gzip.BestCompression)(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusConflict {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "")
	})

	t.Run("should default to OK with no code set", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Range", "bytes=0-1023")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression)(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})

	t.Run("does not encode response (wrong Accept-Encoding)", func(t *testing.T) {
		// setup
		req := httptest.NewRequest(http.MethodPost, "/foo/bar", strings.NewReader(`{"payload":42}`))
		req.Header.Set("Accept-Encoding", "nothing")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)

		// test
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			resp := shogoa.ContextResponse(ctx)
			if _, err := resp.Write([]byte("gzip me!")); err != nil {
				return err
			}
			return nil
		}
		handler = Middleware(gzip.BestCompression, MinSize(0))(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// verify the result
		result := rw.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		verifyNotGzippedResponse(t, result, "gzip me!")
	})
}
