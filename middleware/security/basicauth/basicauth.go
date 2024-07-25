package basicauth

import (
	"context"
	"net/http"

	"github.com/shogo82148/shogoa"
)

// ErrBasicAuthFailed means it wasn't able to authenticate you with your login/password.
var ErrBasicAuthFailed = shogoa.NewErrorClass("basic_auth_failed", 401)

// New creates a static username/password auth middleware.
//
// Example:
//
//	app.UseBasicAuth(basicauth.New("admin", "password"))
//
// It doesn't get simpler than that.
//
// If you want to handle the username and password checks dynamically,
// copy the source of `New`, it's 8 lines and you can tweak at will.
func New(username, password string) shogoa.Middleware {
	middleware, _ := shogoa.NewMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			return ErrBasicAuthFailed("Authentication failed")
		}
		return nil
	})
	return middleware
}
