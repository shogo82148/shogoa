package middleware

import (
	"context"
	"net/http"
	"regexp"

	"github.com/shogo82148/shogoa"
)

// RequireHeader requires a request header to match a value pattern. If the
// header is missing or does not match then the failureStatus is the response
// (e.g. http.StatusUnauthorized). If pathPattern is nil then any path is
// included. If requiredHeaderValue is nil then any value is accepted so long as
// the header is non-empty.
func RequireHeader(
	service *shogoa.Service,
	pathPattern *regexp.Regexp,
	requiredHeaderName string,
	requiredHeaderValue *regexp.Regexp,
	failureStatus int,
) shogoa.Middleware {

	return func(h shogoa.Handler) shogoa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) (err error) {
			if pathPattern != nil && !pathPattern.MatchString(req.URL.Path) {
				return h(ctx, rw, req)
			}

			matched := false
			headerValue := req.Header.Get(requiredHeaderName)
			if len(headerValue) > 0 {
				if requiredHeaderValue == nil {
					matched = true
				} else {
					matched = requiredHeaderValue.MatchString(headerValue)
				}
			}
			if matched {
				err = h(ctx, rw, req)
			} else {
				err = service.Send(ctx, failureStatus, http.StatusText(failureStatus))
			}
			return
		}
	}
}
