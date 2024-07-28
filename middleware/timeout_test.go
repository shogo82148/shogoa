package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shogo82148/shogoa"
)

func TestTimeout(t *testing.T) {
	service := shogoa.New("test")
	service.Encoder.Register(shogoa.NewJSONEncoder, "*/*")
	service.Decoder.Register(shogoa.NewJSONDecoder, "*/*")

	handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		_, ok := ctx.Deadline()
		if !ok {
			t.Error("expected a deadline")
		}
		return service.Send(ctx, http.StatusOK, "ok")
	}
	handler = Timeout(time.Second)(handler)

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	ctx := shogoa.NewContext(rw, req, nil)
	if err := handler(ctx, rw, req); err != nil {
		t.Fatal(err)
	}
}
