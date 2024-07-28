package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shogo82148/shogoa"
)

// ErrorHandler turns a Go error into an HTTP response. It should be placed in the middleware chain
// below the logger middleware so the logger properly logs the HTTP response. ErrorHandler
// understands instances of shogoa.ServiceError and returns the status and response body embodied in
// them, it turns other Go error types into a 500 internal error response.
// If verbose is false the details of internal errors is not included in HTTP responses.
// If you use github.com/pkg/errors then wrapping the error will allow a trace to be printed to the logs
func ErrorHandler(service *shogoa.Service, verbose bool) shogoa.Middleware {
	return func(h shogoa.Handler) shogoa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			e := h(ctx, rw, req)
			if e == nil {
				return nil
			}
			cause := cause(e)
			status := http.StatusInternalServerError
			var respBody any
			if err, ok := cause.(shogoa.ServiceError); ok {
				status = err.ResponseStatus()
				respBody = err
				shogoa.ContextResponse(ctx).ErrorCode = err.Token()
				rw.Header().Set("Content-Type", shogoa.ErrorMediaIdentifier)
			} else {
				respBody = e.Error()
				rw.Header().Set("Content-Type", "text/plain")
			}
			if status == http.StatusInternalServerError {
				reqID := ctx.Value(reqIDKey)
				if reqID == nil {
					reqID = shortID()
					ctx = context.WithValue(ctx, reqIDKey, reqID)
				}
				shogoa.LogError(ctx, "uncaught error", "err", fmt.Sprintf("%+v", e), "id", reqID, "msg", respBody)
				if !verbose {
					rw.Header().Set("Content-Type", shogoa.ErrorMediaIdentifier)
					msg := fmt.Sprintf("%s [%s]", http.StatusText(http.StatusInternalServerError), reqID)
					respBody = shogoa.ErrInternal(msg)
					// Preserve the ID of the original error as that's what gets logged, the client
					// received error ID must match the original
					if origErrID := shogoa.ContextResponse(ctx).ErrorCode; origErrID != "" {
						respBody.(*shogoa.ErrorResponse).ID = origErrID
					}
				}
			}
			return service.Send(ctx, status, respBody)
		}
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func cause(e error) error {
	type causer interface {
		Cause() error
	}
	for {
		cause, ok := e.(causer)
		if !ok {
			break
		}
		c := cause.Cause()
		if c == nil {
			break
		}
		e = c
	}
	return e
}
