package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shogo82148/shogoa"
)

func TestLogResponse(t *testing.T) {
	// create a new service with a JSON logger
	var buf bytes.Buffer
	service := shogoa.New("test")
	logHandler := slog.NewJSONHandler(&buf, nil)
	service.WithLogger(shogoa.NewLogger(logHandler))
	service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
	service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

	// create a new request context
	req := httptest.NewRequest(http.MethodPost, "/goo", strings.NewReader(`{"payload":42}`))
	rw := httptest.NewRecorder()
	params := url.Values{"query": []string{"value"}}
	ctx := shogoa.NewContext(rw, req.WithContext(service.Context), params)

	// call the handler
	handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		shogoa.ContextResponse(ctx).WriteHeader(200)
		if _, err := shogoa.ContextResponse(ctx).Write([]byte("some response data to be logged")); err != nil {
			return err
		}
		return nil
	}
	handler = LogResponse()(handler)
	if err := handler(ctx, rw, req); err != nil {
		t.Fatal(err)
	}

	// check the log entry
	var entry map[string]any
	decoder := json.NewDecoder(&buf)
	if err := decoder.Decode(&entry); err != nil {
		t.Fatal(err)
	}
	if entry["body"] != "some response data to be logged" {
		t.Errorf("unexpected body: %v", entry["body"])
	}
}
