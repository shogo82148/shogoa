package middleware

import (
	"context"
	"net/http"

	"github.com/shogo82148/shogoa"
	"github.com/shogo82148/shogoa/internal/randid"
)

const (
	// RequestIDHeader is the name of the header used to transmit the request ID.
	RequestIDHeader = "X-Request-Id"

	// DefaultRequestIDLengthLimit is the default maximum length for the request ID header value.
	DefaultRequestIDLengthLimit = 128
)

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
				id = randid.New(24)
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
