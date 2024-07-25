package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/shogo82148/shogoa"
)

const (
	// RequestIDHeader is the name of the header used to transmit the request ID.
	RequestIDHeader = "X-Request-Id"

	// DefaultRequestIDLengthLimit is the default maximum length for the request ID header value.
	DefaultRequestIDLengthLimit = 128
)

// Counter used to create new request ids.
var reqID int64

// Common prefix to all newly created request ids for this process.
var reqPrefix string

// Initialize common prefix on process startup.
func init() {
	// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L44-L50
	var buf [12]byte
	var b64 string
	replacer := strings.NewReplacer("+", "", "/", "")
	for len(b64) < 10 {
		if _, err := rand.Read(buf[:]); err != nil {
			panic(err)
		}
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = replacer.Replace(b64)
	}
	reqPrefix = string(b64[0:10])
}

// RequestIDWithHeader behaves like the middleware RequestID, but it takes the request id header
// as the (first) argument.
func RequestIDWithHeader(requestIDHeader string) shogoa.Middleware {
	return RequestIDWithHeaderAndLengthLimit(requestIDHeader, DefaultRequestIDLengthLimit)
}

// RequestIDWithHeaderAndLengthLimit behaves like the middleware RequestID, but it takes the
// request id header as the (first) argument and a length limit for truncation of the request
// header value if it exceeds a reasonable length. The limit can be negative for unlimited.
func RequestIDWithHeaderAndLengthLimit(requestIDHeader string, lengthLimit int) shogoa.Middleware {
	return func(h shogoa.Handler) shogoa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			id := req.Header.Get(requestIDHeader)
			if id == "" {
				id = fmt.Sprintf("%s-%d", reqPrefix, atomic.AddInt64(&reqID, 1))
			} else if lengthLimit >= 0 && len(id) > lengthLimit {
				id = id[:lengthLimit]
			}
			ctx = context.WithValue(ctx, reqIDKey, id)

			return h(ctx, rw, req)
		}
	}
}

// RequestID is a middleware that injects a request ID into the context of each request.
// Retrieve it using ctx.Value(ReqIDKey). If the incoming request has a RequestIDHeader header then
// that value is used else a random value is generated.
func RequestID() shogoa.Middleware {
	return RequestIDWithHeader(RequestIDHeader)
}

// ContextRequestID extracts the Request ID from the context.
func ContextRequestID(ctx context.Context) (reqID string) {
	id := ctx.Value(reqIDKey)
	if id != nil {
		reqID = id.(string)
	}
	return
}
