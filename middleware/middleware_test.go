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
func newContext(service *shogoa.Service, rw http.ResponseWriter, req *http.Request, params url.Values) context.Context {
	ctrl := service.NewController("test")
	return shogoa.NewContext(ctrl.Context, rw, req, params)
}

type logEntry struct {
	Msg  string
	Data []interface{}
}

type testLogger struct {
	Context      []interface{}
	InfoEntries  []logEntry
	WarnEntries  []logEntry
	ErrorEntries []logEntry
}

func (t *testLogger) Info(msg string, data ...interface{}) {
	e := logEntry{msg, append(t.Context, data...)}
	t.InfoEntries = append(t.InfoEntries, e)
}

func (t *testLogger) Warn(msg string, data ...interface{}) {
	e := logEntry{msg, append(t.Context, data...)}
	t.WarnEntries = append(t.InfoEntries, e)
}

func (t *testLogger) Error(msg string, data ...interface{}) {
	e := logEntry{msg, append(t.Context, data...)}
	t.ErrorEntries = append(t.ErrorEntries, e)
}

func (t *testLogger) New(data ...interface{}) shogoa.LogAdapter {
	t.Context = append(t.Context, data...)
	return t
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
