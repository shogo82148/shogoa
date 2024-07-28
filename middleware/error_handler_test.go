package middleware

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/shogo82148/shogoa"
)

// errorResponse contains the details of a error response. It implements ServiceError.
type errorResponse struct {
	// ID is the unique error instance identifier.
	ID string `json:"id" yaml:"id" xml:"id" form:"id"`
	// Code identifies the class of errors.
	Code string `json:"code" yaml:"code" xml:"code" form:"code"`
	// Status is the HTTP status code used by responses that cary the error.
	Status int `json:"status" yaml:"status" xml:"status" form:"status"`
	// Detail describes the specific error occurrence.
	Detail string `json:"detail" yaml:"detail" xml:"detail" form:"detail"`
	// Meta contains additional key/value pairs useful to clients.
	Meta map[string]any `json:"meta,omitempty" yaml:"meta,omitempty" xml:"meta,omitempty" form:"meta,omitempty"`
}

// Error returns the error occurrence details.
func (e *errorResponse) Error() string {
	msg := fmt.Sprintf("[%s] %d %s: %s", e.ID, e.Status, e.Code, e.Detail)
	for k, v := range e.Meta {
		msg += ", " + fmt.Sprintf("%s: %v", k, v)
	}
	return msg
}

func TestErrorHandler(t *testing.T) {
	t.Run("verbose", func(t *testing.T) {
		// build a service
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return errors.New("boom")
		}
		errorHandler := ErrorHandler(service, true)(handler)

		// run the handler
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := errorHandler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// check the result
		result := rw.Result()
		if result.StatusCode != http.StatusInternalServerError {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("unexpected content type: %s", result.Header.Get("Content-Type"))
		}
		body, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != `"boom"`+"\n" {
			t.Errorf("unexpected body: %s", body)
		}
	})

	t.Run("not verbose", func(t *testing.T) {
		// build a service
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return errors.New("boom")
		}
		errorHandler := ErrorHandler(service, false)(handler)

		// run the handler
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := errorHandler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// check the result
		result := rw.Result()
		if result.StatusCode != http.StatusInternalServerError {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Content-Type") != shogoa.ErrorMediaIdentifier {
			t.Errorf("unexpected content type: %s", result.Header.Get("Content-Type"))
		}
		var decoded errorResponse
		if err := service.Decoder.Decode(&decoded, result.Body, "application/json"); err != nil {
			t.Fatal(err)
		}
		if decoded.Code != "internal" {
			t.Errorf("unexpected code: %s", decoded.Code)
		}
		if ok, err := regexp.MatchString(`Internal Server Error \[.*\]`, decoded.Detail); err != nil || !ok {
			t.Errorf("unexpected detail: %s", decoded.Detail)
		}
		if decoded.Status != http.StatusInternalServerError {
			t.Errorf("unexpected status: %d", decoded.Status)
		}
	})

	t.Run("shogoa 500 error", func(t *testing.T) {
		// build a service
		shogoaErr := shogoa.ErrInternal("shogoa-500-boom")
		origID := shogoaErr.(shogoa.ServiceError).Token()
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return shogoaErr
		}
		errorHandler := ErrorHandler(service, false)(handler)

		// run the handler
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := errorHandler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// check the result
		result := rw.Result()
		if result.StatusCode != http.StatusInternalServerError {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Content-Type") != shogoa.ErrorMediaIdentifier {
			t.Errorf("unexpected content type: %s", result.Header.Get("Content-Type"))
		}
		var decoded errorResponse
		if err := service.Decoder.Decode(&decoded, result.Body, "application/json"); err != nil {
			t.Fatal(err)
		}
		if decoded.ID != origID {
			t.Errorf("unexpected ID: %s", decoded.ID)
		}
		if decoded.Code != "internal" {
			t.Errorf("unexpected code: %s", decoded.Code)
		}
		if ok, err := regexp.MatchString(`Internal Server Error \[.*\]`, decoded.Detail); err != nil || !ok {
			t.Errorf("unexpected detail: %s", decoded.Detail)
		}
		if decoded.Status != http.StatusInternalServerError {
			t.Errorf("unexpected status: %d", decoded.Status)
		}
	})

	t.Run("shogoa 504 error", func(t *testing.T) {
		// build a service
		meaningful := shogoa.NewErrorClass("shogoa-504-boom", http.StatusGatewayTimeout)
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return meaningful("gatekeeper says no")
		}
		errorHandler := ErrorHandler(service, false)(handler)

		// run the handler
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := errorHandler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// check the result
		result := rw.Result()
		if result.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Content-Type") != shogoa.ErrorMediaIdentifier {
			t.Errorf("unexpected content type: %s", result.Header.Get("Content-Type"))
		}
		var decoded errorResponse
		if err := service.Decoder.Decode(&decoded, result.Body, "application/json"); err != nil {
			t.Fatal(err)
		}
		if decoded.Code != "shogoa-504-boom" {
			t.Errorf("unexpected code: %s", decoded.Code)
		}
		if decoded.Detail != "gatekeeper says no" {
			t.Errorf("unexpected detail: %s", decoded.Detail)
		}
		if decoded.Status != http.StatusGatewayTimeout {
			t.Errorf("unexpected status: %d", decoded.Status)
		}
	})

	t.Run("custom shogoa error", func(t *testing.T) {
		// build a service
		gerr := shogoa.NewErrorClass("code", http.StatusTeapot)("teapot", "foobar", 42)
		service := shogoa.New("test")
		service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
		service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return gerr
		}
		errorHandler := ErrorHandler(service, true)(handler)

		// run the handler
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		if err := errorHandler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}

		// check the result
		result := rw.Result()
		if result.StatusCode != http.StatusTeapot {
			t.Errorf("unexpected status code: %d", result.StatusCode)
		}
		if result.Header.Get("Content-Type") != shogoa.ErrorMediaIdentifier {
			t.Errorf("unexpected content type: %s", result.Header.Get("Content-Type"))
		}
		var decoded errorResponse
		if err := service.Decoder.Decode(&decoded, result.Body, "application/json"); err != nil {
			t.Fatal(err)
		}
		if decoded.Code != "code" {
			t.Errorf("unexpected code: %s", decoded.Code)
		}
		if decoded.Detail != "teapot" {
			t.Errorf("unexpected detail: %s", decoded.Detail)
		}
		if decoded.Status != http.StatusTeapot {
			t.Errorf("unexpected status: %d", decoded.Status)
		}
	})
}
