package jwt

import (
	"context"

	jwt "github.com/golang-jwt/jwt/v4"
)

type contextKey struct{}

var jwtKey = contextKey{}

// WithJWT creates a child context containing the given JWT.
func WithJWT(ctx context.Context, t *jwt.Token) context.Context {
	return context.WithValue(ctx, jwtKey, t)
}

// ContextJWT retrieves the JWT token from a `context` that went through our security middleware.
func ContextJWT(ctx context.Context) *jwt.Token {
	token, ok := ctx.Value(jwtKey).(*jwt.Token)
	if !ok {
		return nil
	}
	return token
}
