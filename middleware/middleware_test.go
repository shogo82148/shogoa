package middleware_test

import (
	"context"
	"net/http"
	"net/url"

	"github.com/shogo82148/shogoa"
)

// Helper that sets up a "working" service
func newService(logger shogoa.LogAdapter) *shogoa.Service {
	service := shogoa.New("test")
	service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
	service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
	service.WithLogger(logger)
	return service
}

// Creates a test context
func newContext(_ *shogoa.Service, rw http.ResponseWriter, req *http.Request, params url.Values) context.Context {
	return shogoa.NewContext(rw, req, params)
}

type testResponseWriter struct {
	ParentHeader http.Header
	Body         []byte
	Status       int
}

func newTestResponseWriter() *testResponseWriter {
	h := make(http.Header)
	return &testResponseWriter{ParentHeader: h}
}

func (t *testResponseWriter) Header() http.Header {
	return t.ParentHeader
}

func (t *testResponseWriter) Write(b []byte) (int, error) {
	t.Body = append(t.Body, b...)
	return len(b), nil
}

func (t *testResponseWriter) WriteHeader(s int) {
	t.Status = s
}
