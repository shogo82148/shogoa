package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtpkg "github.com/golang-jwt/jwt/v4"
	"github.com/shogo82148/shogoa"
)

func TestMiddleware_HMAC(t *testing.T) {
	t.Run("HMAC keys signed token", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New("keys", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("HMAC fails", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("should not be called")
		}

		middleware := New("otherkey", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("multiple keys", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New([]string{"firstkey", "keys"}, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should accept scopes specified using the 'scope' claim", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scope":"scope1","admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6InNjb3BlMSIsImFkbWluIjp0cnVlfQ.EwMZtpTUPUoKsiCHqH659JQeMLf3-KdboStmQKjv2IU")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		ctx = shogoa.WithRequiredScopes(ctx, []string{"scope1"})
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New("keys", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should accept scopes specified using the 'scopes' claim", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		ctx = shogoa.WithRequiredScopes(ctx, []string{"scope1"})
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New("keys", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should fall back to 'scopes' if 'scope' is null", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scope":null, "scopes":"scope1", "admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6bnVsbCwic2NvcGVzIjoic2NvcGUxIiwiYWRtaW4iOnRydWV9.h8L_MlWWyB0RnwaUBDVu8nGPn5wPSVPMEm42iH8Jxmg")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		ctx = shogoa.WithRequiredScopes(ctx, []string{"scope1"})
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New("keys", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should not fall back to 'scopes' if 'scope' is an empty string", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scope":"", "scopes":"scope1", "admin":true}, signed with "keys"
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6IiIsInNjb3BlcyI6InNjb3BlMSIsImFkbWluIjp0cnVlfQ.U5r-gAvk8SWRYBK3Hmj7zqHSQ0lSQO1wAAk0soyHkoU")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		ctx = shogoa.WithRequiredScopes(ctx, []string{"scope1"})
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("should not be called")
		}

		middleware := New("keys", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestMiddleware_RSA(t *testing.T) {
	t.Run("a single key", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// RS256 {"scopes":"scope1 scope2","admin":true}, signed with rsaKey1 below
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEgc2NvcGUyIiwiYWRtaW4iOnRydWV9.gT4gSGqXTCUZAJT_TWZ4eknazVo-ulMKwSpHoghWZU8Sm9QXt48ISwFAb_wW2xhR58MUNX95iuiex0bCWvze59r35dEQ2SOZixuDvE8srQi2SRk9qqsVV9-R361qf2D8KfLX9jQ7j-UB40bleg0fOyBAjPLPq0ggBigSjQ2yUz8YDKma-n6Ulc3LJ4gyozmb3MjO9RV2pdD3N-m6ttwkTkUE2jhsL6a3T8f0Y6xSGTMyZasKc6kHbUyz6NjAeplLhbkBDE8-Ak4GaLGlLnLzZ49oTVrh89yauciW5yLQCXzXt2PODqp6zXPC0FFcDr-2USCpA-nqaQQyhliMcgtqVw")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New(rsaPubKey1, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("keys that didn't the JWT", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// RS256 {"scopes":"scope1 scope2","admin":true}, signed with rsaKey1 below
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEgc2NvcGUyIiwiYWRtaW4iOnRydWV9.gT4gSGqXTCUZAJT_TWZ4eknazVo-ulMKwSpHoghWZU8Sm9QXt48ISwFAb_wW2xhR58MUNX95iuiex0bCWvze59r35dEQ2SOZixuDvE8srQi2SRk9qqsVV9-R361qf2D8KfLX9jQ7j-UB40bleg0fOyBAjPLPq0ggBigSjQ2yUz8YDKma-n6Ulc3LJ4gyozmb3MjO9RV2pdD3N-m6ttwkTkUE2jhsL6a3T8f0Y6xSGTMyZasKc6kHbUyz6NjAeplLhbkBDE8-Ak4GaLGlLnLzZ49oTVrh89yauciW5yLQCXzXt2PODqp6zXPC0FFcDr-2USCpA-nqaQQyhliMcgtqVw")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("should not be called")
		}

		middleware := New(rsaPubKey2, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("multiple keys", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// RS256 {"scopes":"scope1 scope2","admin":true}, signed with rsaKey1 below
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEgc2NvcGUyIiwiYWRtaW4iOnRydWV9.gT4gSGqXTCUZAJT_TWZ4eknazVo-ulMKwSpHoghWZU8Sm9QXt48ISwFAb_wW2xhR58MUNX95iuiex0bCWvze59r35dEQ2SOZixuDvE8srQi2SRk9qqsVV9-R361qf2D8KfLX9jQ7j-UB40bleg0fOyBAjPLPq0ggBigSjQ2yUz8YDKma-n6Ulc3LJ4gyozmb3MjO9RV2pdD3N-m6ttwkTkUE2jhsL6a3T8f0Y6xSGTMyZasKc6kHbUyz6NjAeplLhbkBDE8-Ak4GaLGlLnLzZ49oTVrh89yauciW5yLQCXzXt2PODqp6zXPC0FFcDr-2USCpA-nqaQQyhliMcgtqVw")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New([]*rsa.PublicKey{rsaPubKey1}, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMiddleware_ECDSA(t *testing.T) {
	t.Run("a single key", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// RS256 {"scopes":"scope1 scope2","admin":true}, signed with rsaKey1 below
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwic2NvcGVzIjoic2NvcGUxIHNjb3BlMiJ9.7gM4EblP4cvX5C6PBLSBFpKX2FQ9AsLNmOXEm86uvrd4czBfw1zDO24abQ7gtlbMcjuVvxrpIyRa7Nbbn31G7w")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New(ecPubKey1, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with keys that didn't the JWT", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// RS256 {"scopes":"scope1 scope2","admin":true}, signed with rsaKey1 below
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwic2NvcGVzIjoic2NvcGUxIHNjb3BlMiJ9.7gM4EblP4cvX5C6PBLSBFpKX2FQ9AsLNmOXEm86uvrd4czBfw1zDO24abQ7gtlbMcjuVvxrpIyRa7Nbbn31G7w")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("should not be called")
		}

		middleware := New(ecPubKey2, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("multiple keys", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocHeader,
			Name: "Authorization",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// RS256 {"scopes":"scope1 scope2","admin":true}, signed with rsaKey1 below
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwic2NvcGVzIjoic2NvcGUxIHNjb3BlMiJ9.7gM4EblP4cvX5C6PBLSBFpKX2FQ9AsLNmOXEm86uvrd4czBfw1zDO24abQ7gtlbMcjuVvxrpIyRa7Nbbn31G7w")
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New([]*ecdsa.PublicKey{ecPubKey1}, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMiddleware_HMAC_Query(t *testing.T) {
	t.Run("HMAC keys signed token", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocQuery,
			Name: "access_token",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		q := req.URL.Query()
		q.Set("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		req.URL.RawQuery = q.Encode()
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New("keys", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with keys that didn't the JWT", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocQuery,
			Name: "access_token",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		q := req.URL.Query()
		q.Set("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		req.URL.RawQuery = q.Encode()
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			panic("should not be called")
		}

		middleware := New("otherkey", nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("multiple keys", func(t *testing.T) {
		securityScheme := &shogoa.JWTSecurity{
			In:   shogoa.LocQuery,
			Name: "access_token",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// HS256 {"scopes":"scope1","admin":true}, signed with "keys"
		q := req.URL.Query()
		q.Set("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOiJzY29wZTEiLCJhZG1pbiI6dHJ1ZX0.UCvEfbD_yuS5dCZidxZgogVi2yF0ZVecMsQQbY1HJy0")
		req.URL.RawQuery = q.Encode()
		rw := httptest.NewRecorder()
		ctx := shogoa.NewContext(rw, req, nil)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			token := ContextJWT(ctx)
			if token == nil {
				t.Fatal("token is nil")
			}
			return nil
		}

		middleware := New([]string{"firstkey", "keys"}, nil, securityScheme)
		handler = middleware(handler)
		if err := handler(ctx, rw, req); err != nil {
			t.Fatal(err)
		}
	})
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// the private key of rsaPubKey1
/*
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEArZIJcPQd7aSGb80wgFpy5SVjzzsGpfIysZ30SdWlTcWMVbAT
XmsDNgw98TzIeoyikSbSHEeORbKWKS2clgNsdLjYKv3XLTBaXfLcU3x9mhnk/kUL
N/AQgyvsRGynPris2oVzGSib7uOZK/9+u+QAKIrp7prcmMmnwvdcjFXjwzx83RTF
1b+iuVGCdV0T4m1XQdm/YtIUh7JNbYrUolkdwZlOxMZuV0FDC+ms02+gyj580Pyl
TuAD4JmtSmmijyWfEx5dsZYtGALyUxcm5Hz15RP3FACrv4B++BHI6smO4sWdrSYV
l3sHJ60Bm6zbwuyB2twJPOdL5nVIGiIDdf+1IwIDAQABAoIBACF3MtLQfqS/QBx2
V4+n4NdFqkTegJ2mYRfV+K/zvXeNS02KMXHW+DuHiLnmmlDmpMGP1psVQN03XbR6
0uIprtOigCtp3f8cf4/1r315V05LB9fuwAb9BnIEGf3nZSe2u702VcbYCZi50WKm
VG0tvMoUXp5exYG//9SblQCJ3uxZf9D8y5RnrUZtP4Pnjkn/3YeJCF+Kked55Cvi
gv56/aiyWp9xEGsSWig5Zt8VNXihgT7D2KZzxcQDQlxw0CR5ECT7/4w7sZVvwc7B
I76JJDvpD0UGvzoUgx928efGKxJBrcjzvTNSKgHJYYCvaa6+qX2tjkmOqdG4xl27
/TaBISECgYEA4YJ32HKaS2ikn5J2C1qlHy4nRuVNhX8T9qvp6OBdbE2BQq3O5IUt
+wdTcjhD88aDdHCZmAd8i3FC4l+CKsQ5sTwRk0vTOZ7axC6+zDHg+na5/+NCq+xm
ffoaZ5jsZFyqfbsFn8NiLWLo2JSFV1AnUxwpPA2jbuylAuZVvVbLYHcCgYEAxQnO
L+U6NwTvN6EJgt3h1jHfXiQnapsj0O0XgK+g2K6vMovpXAf9noO+r3Qbx/aKxuRg
TvRQ08T5yyqysz+sYe0rp8oaMUhYQFMgJOUUBlmUVGxYdARMD6kTy/384B9Azoex
UCosMSEAD909MAsyQWB4X6OJKd+V68QpFYeIx7UCgYBHgaRY6PYOBU92He36abLE
MVFZBKrRMtt0s0yHgGV/SxA6wXxCMAzFdaw7IqZBbWgPiwjZET6nxLFNsLVItFIK
5h44k6mVss5xuNTdUM+i+/S8tCZW964EMkMfKHmE1XFmTuBYqY6/D4b/7hBeAFeH
3f0hQr3ZFYa5Zao4UIZKvwKBgGL8lhUBt8lENVlhEYIpLfeJfomw6AxqfAfN1GzV
zpyMxX9DQqz1ZrhnvzgtwHcoqHda6/c+TgzVfBhRDw12A4f+ulvE8HupuIw4NoHS
g8jc3+O5uoYuUnfbnRJyOsPtb4VSLgXz6deUmI9fugmU1l55tH93jMT4ijyzg2BJ
grGxAoGAWX24Yx9qoasqEQ2rgdTsgylwL28UczKQ5KNHt2PcEfPNw6/GpfK7YmlU
Heef2umEzb1K2ZK95wlMbF8zpNDWBf4PkxgfW+JEE+pO1kb5KXysBymymyXhGHAP
CwH9XHqbjVlsD358AbPeKqLgTCaGo9JgsEZDBpESmBDnIPUahMc=
-----END RSA PRIVATE KEY-----
*/

var rsaPubKey1 = must(jwtpkg.ParseRSAPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArZIJcPQd7aSGb80wgFpy
5SVjzzsGpfIysZ30SdWlTcWMVbATXmsDNgw98TzIeoyikSbSHEeORbKWKS2clgNs
dLjYKv3XLTBaXfLcU3x9mhnk/kULN/AQgyvsRGynPris2oVzGSib7uOZK/9+u+QA
KIrp7prcmMmnwvdcjFXjwzx83RTF1b+iuVGCdV0T4m1XQdm/YtIUh7JNbYrUolkd
wZlOxMZuV0FDC+ms02+gyj580PylTuAD4JmtSmmijyWfEx5dsZYtGALyUxcm5Hz1
5RP3FACrv4B++BHI6smO4sWdrSYVl3sHJ60Bm6zbwuyB2twJPOdL5nVIGiIDdf+1
IwIDAQAB
-----END PUBLIC KEY-----`)))

// the private key of rsaPubKey2
/*
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4jr/DGbPt0UDGvu6Xo2LV0F6Wf8OnyxF2IFPdG5B4X0YS3DC
9SF3clbbBivDVa2bEXppyj+eLEKlfohCWXTrJK0LxTEcneuDkF4re+BdP3q9cKRz
FtI/ZVhVnD7+PS1wps7OiTM0iOaIDo9+uFrC6zBTRAiPyrdwh1ApttLdoD6i5D9D
7zzvpTXLC/UWaRz/phAaaop6dPPR1YblZEckWgqTMC3KrRX/6QJFFfpgyQzFT09W
DYnmXl2gS7C2sk4UejygqmVg96JxaIaT3WiQSjxXddjR/krcA9EGNNEkpZB2W6Ux
6d63yWsNG9YJUacwI+M2q5ZW964J1s//FiNZZQIDAQABAoIBAQCoqYtU16Gs5Qq3
p0z/CVAFMY/iYMGp8fvwuhdemoULc5QVSnBPCTBgUljgdOggjFm74iPU4TEvllCD
0VqGDyDwKwNHdKH9KoTfsRWCOXfLx9pMjI4xSXZyPDU3U8+AFMyT0EMzDrXwCs8M
6/Zxw1jmtxSc+DUb0T9X4m/3GaaZvDGGShnU8/XnEh2uEHrNwWnGWYPJ/rZjNZPy
PZ9W2VpcHKBMVEowK/cOoouNuflAISoLCCLMNYygr9T4Ylm3HGP9o7JuWL+wGQsa
aXrE5qTOpsxmBqTQ8pglnxnhDEFXmx3O+bwRfIwDSYe+wvCINpdIstWuybh4Ed2i
ZgLTlx8BAoGBAP9LwmfZ/2XNHBzk+f09TnTnhXzVsKkHu5BlXvWoDigVv4Dzl44j
X1Ade5PjiOf0Jti2QCkAaI+CjExdP1zCqDZBQFpKI3QQgvlWoKXHVFV9ziC8gcX+
I6M8wmtIoK8ISnC6A5s1wKIvOPsZyP7aVZgu805BKfVqtFWCK42vnRVRAoGBAOLa
t2pOzVttd3vPgzGovD+Mf3RsPg6ygazj0GiDRspRCnoeopFEoBPFcKIQZlPp8rfT
NLOuwVkW5TqntrCW0UwixZMXicIaPDo0idXInIfP0+f7JxSYb5q7vmbyRt8uAYY9
GU4L/ZIn127JbgQ5n5nuODMvTe7m5Ky+FUYHw43VAoGAE6QOdtLstTZMfWSYXwVC
bfgJ6wq9pqNzqK5D2f5t6GOT8iXLeSH7iTxbb4tH0yCThISw9vaTFMdkZ9OctlQ7
gMEQZGHjzGAg03H4tghZ0qH1I8uc6FCfCUX5ZyuVQSIQKBAHiv9drJyZc6gOMJ03
jJfAHDsjMUBeU13KYAIswaECgYBTYiNSzv5KodTuTFsjsKrpDOJ4T6ULz+88NkyP
bdliWiFou8Pzc28HdWYuG6sRIwfVK6vOc+ibr3+4bJcJF5Z8zrcilt9K2kvS9SbI
zsFCZlC0jytRNaqoDGQzANCuDgH/bovTlTKyOzTDgwSORwP0F4zOu4+AxZu+Juw4
3nextQKBgEAGLuChkztZCVt0W2D8wJYFR7XjezcbsfpoXx9H8htk6u4STu9TwB76
DxoYj3qiTV2kRRBQQZRAli1TbDOnJuqFMnRL0aPsqebuW2sqY9Hx9G6TxokN8Nc6
RlTE+CbPcjBgAx+AANL/X2KYoXLAjOrYY5kQD8Qbt8Wkme7m6hiP
-----END RSA PRIVATE KEY-----
*/

var rsaPubKey2 = must(jwtpkg.ParseRSAPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4jr/DGbPt0UDGvu6Xo2L
V0F6Wf8OnyxF2IFPdG5B4X0YS3DC9SF3clbbBivDVa2bEXppyj+eLEKlfohCWXTr
JK0LxTEcneuDkF4re+BdP3q9cKRzFtI/ZVhVnD7+PS1wps7OiTM0iOaIDo9+uFrC
6zBTRAiPyrdwh1ApttLdoD6i5D9D7zzvpTXLC/UWaRz/phAaaop6dPPR1YblZEck
WgqTMC3KrRX/6QJFFfpgyQzFT09WDYnmXl2gS7C2sk4UejygqmVg96JxaIaT3WiQ
SjxXddjR/krcA9EGNNEkpZB2W6Ux6d63yWsNG9YJUacwI+M2q5ZW964J1s//FiNZ
ZQIDAQAB
-----END PUBLIC KEY-----`)))

// the private key of ecPubKey1
/*
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIM4zAVusfF+Xl4Z5a5LaspGk+OIwGQweubphSqC1R9+VoAoGCCqGSM49
AwEHoUQDQgAE3tWSknhfssUVytNbPz3TB7giFfxKtHsFW27Yls+Ohfuui9NW4eEk
fLOxYkTI9tyoKfh9Dan5kJFA7ZYEwZ0zMQ==
-----END EC PRIVATE KEY-----
*/

var ecPubKey1 = must(jwtpkg.ParseECPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3tWSknhfssUVytNbPz3TB7giFfxK
tHsFW27Yls+Ohfuui9NW4eEkfLOxYkTI9tyoKfh9Dan5kJFA7ZYEwZ0zMQ==
-----END PUBLIC KEY-----`)))

// the private key of ecPubKey2
/*
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKQ7EyFGaYMuFpMLnqK+mBnT9CrWOqzVxsF8wBlGrTq/oAoGCCqGSM49
AwEHoUQDQgAE8IX3mOtLvBpvrylaRjFpadqGrirXh9dkjJfM/t1dnLu5qPhybMIY
tEr3Xs8vYp2wyaSTVKsyj9y+t344T5Bhdw==
-----END EC PRIVATE KEY-----
*/

var ecPubKey2 = must(jwtpkg.ParseECPublicKeyFromPEM([]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8IX3mOtLvBpvrylaRjFpadqGrirX
h9dkjJfM/t1dnLu5qPhybMIYtEr3Xs8vYp2wyaSTVKsyj9y+t344T5Bhdw==
-----END PUBLIC KEY-----`)))
