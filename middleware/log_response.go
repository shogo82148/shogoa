package middleware

import (
	"context"
	"net/http"

	"github.com/shogo82148/shogoa"
)

// loggingResponseWriter wraps an http.ResponseWriter and writes only raw
// response data (as text) to the context logger. assumes status and duration
// are logged elsewhere (i.e. by the LogRequest middleware).
type loggingResponseWriter struct {
	http.ResponseWriter
	ctx context.Context
}

// Write will write raw data to logger and response writer.
func (lrw *loggingResponseWriter) Write(buf []byte) (int, error) {
	shogoa.LogInfo(lrw.ctx, "response", "body", string(buf))
	return lrw.ResponseWriter.Write(buf)
}

// LogResponse creates a response logger middleware.
// Only Logs the raw response data without accumulating any statistics.
func LogResponse() shogoa.Middleware {
	return func(h shogoa.Handler) shogoa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			// chain a new logging writer to the current response writer.
			resp := shogoa.ContextResponse(ctx)
			resp.SwitchWriter(
				&loggingResponseWriter{
					ResponseWriter: resp.SwitchWriter(nil),
					ctx:            ctx,
				})

			// next
			return h(ctx, rw, req)
		}
	}
}
