package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/shogo82148/shogoa"
)

type reqIDKeyType struct{}

// ReqIDKey is the context key used by the RequestID middleware to store the request ID value.
var reqIDKey = reqIDKeyType{}

// LogRequest creates a request logger middleware.
// This middleware is aware of the RequestID middleware and if registered after it leverages the
// request ID for logging.
// If verbose is true then the middleware logs the request and response bodies.
func LogRequest(verbose bool, sensitiveHeaders ...string) shogoa.Middleware {
	var suppressed map[string]struct{}
	if len(sensitiveHeaders) > 0 {
		suppressed = make(map[string]struct{}, len(sensitiveHeaders))
		for _, sh := range sensitiveHeaders {
			suppressed[strings.ToLower(sh)] = struct{}{}
		}
	}

	return func(h shogoa.Handler) shogoa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			reqID := ctx.Value(reqIDKey)
			if reqID == nil {
				reqID = shortID()
			}
			ctx = shogoa.WithLogContext(ctx, "req_id", reqID)
			startedAt := time.Now()
			r := shogoa.ContextRequest(ctx)
			shogoa.LogInfo(ctx, "started", r.Method, r.URL.String(), "from", from(req),
				"ctrl", shogoa.ContextController(ctx), "action", shogoa.ContextAction(ctx))
			if verbose {
				if len(r.Header) > 0 {
					logCtx := make([]any, 2*len(r.Header))
					i := 0
					keys := make([]string, len(r.Header))
					for k := range r.Header {
						keys[i] = k
						i++
					}
					sort.Strings(keys)
					i = 0
					for _, k := range keys {
						v := r.Header[k]
						logCtx[i] = k
						if _, ok := suppressed[strings.ToLower(k)]; ok {
							logCtx[i+1] = "<hidden>"
						} else {
							logCtx[i+1] = any(strings.Join(v, ", "))
						}
						i = i + 2
					}
					shogoa.LogInfo(ctx, "headers", logCtx...)
				}
				if len(r.Params) > 0 {
					logCtx := make([]any, 2*len(r.Params))
					i := 0
					for k, v := range r.Params {
						logCtx[i] = k
						logCtx[i+1] = any(strings.Join(v, ", "))
						i = i + 2
					}
					shogoa.LogInfo(ctx, "params", logCtx...)
				}
				if r.ContentLength > 0 {
					if mp, ok := r.Payload.(map[string]any); ok {
						logCtx := make([]any, 2*len(mp))
						i := 0
						for k, v := range mp {
							logCtx[i] = k
							logCtx[i+1] = v
							i = i + 2
						}
						shogoa.LogInfo(ctx, "payload", logCtx...)
					} else {
						// Not the most efficient but this is used for debugging
						js, err := json.Marshal(r.Payload)
						if err != nil {
							js = []byte("<invalid JSON>")
						}
						shogoa.LogInfo(ctx, "payload", "raw", string(js))
					}
				}
			}
			err := h(ctx, rw, req)
			resp := shogoa.ContextResponse(ctx)
			if code := resp.ErrorCode; code != "" {
				shogoa.LogInfo(ctx, "completed", "status", resp.Status, "error", code,
					"bytes", resp.Length, "time", time.Since(startedAt).String(),
					"ctrl", shogoa.ContextController(ctx), "action", shogoa.ContextAction(ctx))
			} else {
				shogoa.LogInfo(ctx, "completed", "status", resp.Status,
					"bytes", resp.Length, "time", time.Since(startedAt).String(),
					"ctrl", shogoa.ContextController(ctx), "action", shogoa.ContextAction(ctx))
			}
			return err
		}
	}
}

// shortID produces a "unique" 6 bytes long string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
func shortID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// from makes a best effort to compute the request client IP.
func from(req *http.Request) string {
	if f := req.Header.Get("X-Forwarded-For"); f != "" {
		return f
	}
	f := req.RemoteAddr
	ip, _, err := net.SplitHostPort(f)
	if err != nil {
		return f
	}
	return ip
}
