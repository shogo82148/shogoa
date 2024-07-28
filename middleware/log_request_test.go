package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/shogo82148/shogoa"
)

func TestLogRequest(t *testing.T) {
	// newService creates a new shogoa.Service with a JSON logger
	newService := func(w io.Writer) *shogoa.Service {
		service := shogoa.New("test")
		logHandler := slog.NewJSONHandler(w, nil)
		service.WithLogger(shogoa.NewLogger(logHandler))
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
		return service
	}

	newRequest := func(service *shogoa.Service) (context.Context, http.ResponseWriter, *http.Request) {
		ctrl := service.NewController("test")
		req := httptest.NewRequest(http.MethodPost, "/goo?param=value", strings.NewReader(`{"payload":42}`))
		rw := httptest.NewRecorder()
		params := url.Values{"query": []string{"value"}}
		ctx := shogoa.NewContext(rw, req.WithContext(service.Context), params)
		ctx = ctrl.BaseContext(req.WithContext(ctx))
		ctx = shogoa.WithAction(ctx, "goo")
		shogoa.ContextRequest(ctx).Payload = map[string]interface{}{"payload": 42}
		return ctx, rw, req
	}

	t.Run("logs normal request", func(t *testing.T) {
		var buf bytes.Buffer
		service := newService(&buf)

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return service.Send(ctx, 200, "ok")
		}
		logRequest := LogRequest(true)(handler)
		if err := logRequest(newRequest(service)); err != nil {
			t.Fatal(err)
		}

		var entry map[string]any
		decoder := json.NewDecoder(&buf)

		// check the 1st log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["req_id"] == "" {
			t.Error("req_id is empty")
		}
		if entry["POST"] != "/goo?param=value" {
			t.Error("POST is not /goo?param=value")
		}

		// check the 2nd log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["req_id"] == "" {
			t.Error("req_id is empty")
		}
		if entry["query"] != "value" {
			t.Error("query is not value")
		}

		// check the 3rd log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["req_id"] == "" {
			t.Error("req_id is empty")
		}
		if entry["payload"] != float64(42) {
			t.Error("payload is not 42")
		}

		// check the 4th log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["req_id"] == "" {
			t.Error("req_id is empty")
		}
		if entry["status"] != float64(200) {
			t.Error("status is not 200")
		}
		if entry["bytes"] != float64(5) {
			t.Error("bytes is not 5")
		}
		if _, err := time.ParseDuration(entry["time"].(string)); err != nil {
			t.Errorf("time is invalid: %v", err)
		}
		if entry["ctrl"] != "test" {
			t.Errorf("ctrl is not test, got %v", entry["ctrl"])
		}
		if entry["action"] != "goo" {
			t.Errorf("action is not goo, got %v", entry["action"])
		}
	})

	t.Run("logs error codes", func(t *testing.T) {
		var buf bytes.Buffer
		service := newService(&buf)

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return shogoa.MissingParamError("foo")
		}
		handler = ErrorHandler(service, false)(handler)
		handler = LogRequest(false)(handler)
		if err := handler(newRequest(service)); err != nil {
			t.Fatal(err)
		}

		var entry map[string]any
		decoder := json.NewDecoder(&buf)

		// check the 1st log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["req_id"] == "" {
			t.Error("req_id is empty")
		}
		if entry["POST"] != "/goo?param=value" {
			t.Error("POST is not /goo?param=value")
		}

		// check the 2nd log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["req_id"] == "" {
			t.Error("req_id is empty")
		}
		if entry["status"] != float64(http.StatusBadRequest) {
			t.Errorf("status is not %d, got %v", http.StatusBadRequest, entry["status"])
		}
		if entry["error"] == "" {
			t.Error("error is empty")
		}
		if entry["bytes"] != float64(124) {
			t.Errorf("bytes is not 124, got %v", entry["bytes"])
		}
		if _, err := time.ParseDuration(entry["time"].(string)); err != nil {
			t.Errorf("time is invalid: %v", err)
		}
		if entry["ctrl"] != "test" {
			t.Errorf("ctrl is not test, got %v", entry["ctrl"])
		}
		if entry["action"] != "goo" {
			t.Errorf("action is not goo, got %v", entry["action"])
		}
	})

	t.Run("hides secret headers", func(t *testing.T) {
		var buf bytes.Buffer
		service := newService(&buf)

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return service.Send(ctx, 200, "ok")
		}
		logRequest := LogRequest(true, "Secret")(handler)
		ctx, rw, req := newRequest(service)
		req.Header.Add("Secret", "super secret things")
		req.Header.Add("Not-So-Secret", "public")
		if err := logRequest(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		var entry map[string]any
		decoder := json.NewDecoder(&buf)

		// check the 1st log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}

		// check the 2nd log entry
		if err := decoder.Decode(&entry); err != nil {
			t.Fatal(err)
		}
		if entry["Not-So-Secret"] != "public" {
			t.Errorf("Not-So-Secret is not public, got %v", entry["Not-So-Secret"])
		}
		if entry["Secret"] != "<hidden>" {
			t.Errorf("Secret is not <hidden>, got %v", entry["Secret"])
		}
	})
}
