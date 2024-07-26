package shogoa

import (
	"context"
	"fmt"
	"net/http"
)

// Middleware represents the canonical shogoa middleware signature.
type Middleware func(Handler) Handler

// NewMiddleware creates a middleware from the given argument. The allowed types for the
// argument are:
//
//   - a shogoa middleware: shogoa.Middleware or func(shogoa.Handler) shogoa.Handler
//   - a shogoa handler: shogoa.Handler or func(context.Context, http.ResponseWriter, *http.Request) error
//   - an http middleware: func(http.Handler) http.Handler
//   - or an http handler: http.Handler or func(http.ResponseWriter, *http.Request)
//
// An error is returned if the given argument is not one of the types above.
func NewMiddleware(m any) (Middleware, error) {
	switch m := m.(type) {
	case Middleware:
		return m, nil
	case func(Handler) Handler:
		return m, nil
	case Handler:
		return handlerToMiddleware(m), nil
	case func(context.Context, http.ResponseWriter, *http.Request) error:
		return handlerToMiddleware(m), nil
	case func(http.Handler) http.Handler:
		return func(h Handler) Handler {
			return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) (err error) {
				m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					err = h(ctx, w, r)
				})).ServeHTTP(rw, req)
				return
			}
		}, nil
	case http.Handler:
		return httpHandlerToMiddleware(m.ServeHTTP), nil
	case func(http.ResponseWriter, *http.Request):
		return httpHandlerToMiddleware(m), nil
	default:
		return nil, fmt.Errorf("shogoa: invalid middleware %#v", m)
	}
}

// handlerToMiddleware creates a middleware from a raw handler.
// The middleware calls the handler and either breaks the middleware chain if the handler returns
// an error by also returning the error or calls the next handler in the chain otherwise.
func handlerToMiddleware(m Handler) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if err := m(ctx, rw, req); err != nil {
				return err
			}
			return h(ctx, rw, req)
		}
	}
}

// httpHandlerToMiddleware creates a middleware from a http.HandlerFunc.
// The middleware calls the ServerHTTP method exposed by the http handler and then calls the next
// middleware in the chain.
func httpHandlerToMiddleware(m http.HandlerFunc) Middleware {
	return func(h Handler) Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			m.ServeHTTP(rw, req)
			return h(ctx, rw, req)
		}
	}
}
