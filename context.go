package shogoa

import (
	"context"
	"net/http"
	"net/url"
)

// Keys used to store data in context.
var (
	reqKey            = &contextKey{"request"}
	respKey           = &contextKey{"response"}
	ctrlKey           = &contextKey{"controller"}
	actionKey         = &contextKey{"action"}
	logKey            = &contextKey{"logger"}
	errKey            = &contextKey{"error"}
	securityScopesKey = &contextKey{"security-scope"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "shogoa context value " + k.name }

// RequestData provides access to the underlying HTTP request.
type RequestData struct {
	*http.Request

	// Payload returns the decoded request body.
	Payload any
	// Params contains the raw values for the parameters defined in the design including
	// path parameters, query string parameters and header parameters.
	Params url.Values
}

// ResponseData provides access to the underlying HTTP response.
type ResponseData struct {
	http.ResponseWriter

	// The service used to encode the response.
	Service *Service
	// ErrorCode is the code of the error returned by the action if any.
	ErrorCode string
	// Status is the response HTTP status code.
	Status int
	// Length is the response body length.
	Length int
}

// NewContext builds a new shogoa request context.
func NewContext(rw http.ResponseWriter, req *http.Request, params url.Values) context.Context {
	ctx := req.Context()
	request := &RequestData{Request: req, Params: params}
	response := &ResponseData{ResponseWriter: rw}
	ctx = context.WithValue(ctx, respKey, response)
	ctx = context.WithValue(ctx, reqKey, request)
	return ctx
}

// WithAction creates a context with the given action name.
func WithAction(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, actionKey, action)
}

// WithLogger sets the request context logger and returns the resulting new context.
func WithLogger(ctx context.Context, logger LogAdapter) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

// WithLogContext instantiates a new logger by appending the given key/value pairs to the context
// logger and setting the resulting logger in the context.
func WithLogContext(ctx context.Context, keyvals ...any) context.Context {
	logger := ContextLogger(ctx)
	if logger == nil {
		return ctx
	}
	nl := logger.New(keyvals...)
	return WithLogger(ctx, nl)
}

// WithError creates a context with the given error.
func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errKey, err)
}

// ContextController extracts the controller name from the given context.
func ContextController(ctx context.Context) string {
	if c := ctx.Value(ctrlKey); c != nil {
		return c.(string)
	}
	return "<unknown>"
}

// ContextAction extracts the action name from the given context.
func ContextAction(ctx context.Context) string {
	if a := ctx.Value(actionKey); a != nil {
		return a.(string)
	}
	return "<unknown>"
}

// ContextRequest extracts the request data from the given context.
func ContextRequest(ctx context.Context) *RequestData {
	if r := ctx.Value(reqKey); r != nil {
		return r.(*RequestData)
	}
	return nil
}

// ContextResponse extracts the response data from the given context.
func ContextResponse(ctx context.Context) *ResponseData {
	if r := ctx.Value(respKey); r != nil {
		return r.(*ResponseData)
	}
	return nil
}

// ContextLogger extracts the logger from the given context.
func ContextLogger(ctx context.Context) LogAdapter {
	if v := ctx.Value(logKey); v != nil {
		return v.(LogAdapter)
	}
	return nil
}

// ContextError extracts the error from the given context.
func ContextError(ctx context.Context) error {
	if err := ctx.Value(errKey); err != nil {
		return err.(error)
	}
	return nil
}

// SwitchWriter overrides the underlying response writer. It returns the response
// writer that was previously set.
func (r *ResponseData) SwitchWriter(rw http.ResponseWriter) http.ResponseWriter {
	rwo := r.ResponseWriter
	r.ResponseWriter = rw
	return rwo
}

// Written returns true if the response was written, false otherwise.
func (r *ResponseData) Written() bool {
	return r.Status != 0
}

// WriteHeader records the response status code and calls the underlying writer.
func (r *ResponseData) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write records the amount of data written and calls the underlying writer.
func (r *ResponseData) Write(b []byte) (int, error) {
	if !r.Written() {
		r.WriteHeader(http.StatusOK)
	}
	r.Length += len(b)
	return r.ResponseWriter.Write(b)
}
