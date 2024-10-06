package apidsl_test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestSecurity(t *testing.T) {
	t.Run("should have no security DSL when none are defined", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("secure", nil)
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
		if design.Design.SecuritySchemes != nil {
			t.Errorf("SecuritySchemes = %v; want nil", design.Design.SecuritySchemes)
		}
	})

	t.Run("should be the fully valid and well defined, live on the happy path", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("secure", func() {
			apidsl.Host("example.com")
			apidsl.Scheme("http")

			apidsl.BasicAuthSecurity("basic_authz", func() {
				apidsl.Description("desc")
			})

			apidsl.OAuth2Security("googAuthz", func() {
				apidsl.Description("desc")
				apidsl.AccessCodeFlow("/auth", "/token")
				apidsl.Scope("user:read", "Read users")
			})

			apidsl.APIKeySecurity("a_key", func() {
				apidsl.Description("desc")
				apidsl.Query("access_token")
			})

			apidsl.JWTSecurity("jwt", func() {
				apidsl.Description("desc")
				apidsl.Header("Authorization")
				apidsl.TokenURL("/token")
				apidsl.Scope("user:read", "Read users")
				apidsl.Scope("user:write", "Write users")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
		if len(design.Design.SecuritySchemes) != 4 {
			t.Errorf("SecuritySchemes = %v; want 4", design.Design.SecuritySchemes)
		}

		if design.Design.SecuritySchemes[0].Kind != design.BasicAuthSecurityKind {
			t.Errorf("SecuritySchemes[0].Kind = %v; want %v", design.Design.SecuritySchemes[0].Kind, design.BasicAuthSecurityKind)
		}
		if design.Design.SecuritySchemes[0].Description != "desc" {
			t.Errorf("SecuritySchemes[0].Description = %v; want desc", design.Design.SecuritySchemes[0].Description)
		}

		if design.Design.SecuritySchemes[1].Kind != design.OAuth2SecurityKind {
			t.Errorf("SecuritySchemes[1].Kind = %v; want %v", design.Design.SecuritySchemes[1].Kind, design.OAuth2SecurityKind)
		}
		if design.Design.SecuritySchemes[1].AuthorizationURL != "http://example.com/auth" {
			t.Errorf("SecuritySchemes[1].AuthorizationURL = %v; want %v", design.Design.SecuritySchemes[1].AuthorizationURL, "http://example.com/auth")
		}
		if design.Design.SecuritySchemes[1].TokenURL != "http://example.com/token" {
			t.Errorf("SecuritySchemes[1].TokenURL = %v; want %v", design.Design.SecuritySchemes[1].TokenURL, "http://example.com/token")
		}
		if design.Design.SecuritySchemes[1].Flow != "accessCode" {
			t.Errorf("SecuritySchemes[1].Flow = %v; want %v", design.Design.SecuritySchemes[1].Flow, "accessCode")
		}

		if design.Design.SecuritySchemes[2].Kind != design.APIKeySecurityKind {
			t.Errorf("SecuritySchemes[2].Kind = %v; want %v", design.Design.SecuritySchemes[2].Kind, design.APIKeySecurityKind)
		}
		if design.Design.SecuritySchemes[2].In != "query" {
			t.Errorf("SecuritySchemes[2].In = %v; want %v", design.Design.SecuritySchemes[2].In, "query")
		}
		if design.Design.SecuritySchemes[2].Name != "access_token" {
			t.Errorf("SecuritySchemes[2].Name = %v; want %v", design.Design.SecuritySchemes[2].Name, "access_token")
		}

		if design.Design.SecuritySchemes[3].Kind != design.JWTSecurityKind {
			t.Errorf("SecuritySchemes[3].Kind = %v; want %v", design.Design.SecuritySchemes[3].Kind, design.JWTSecurityKind)
		}
		if design.Design.SecuritySchemes[3].TokenURL != "http://example.com/token" {
			t.Errorf("SecuritySchemes[3].TokenURL = %v; want %v", design.Design.SecuritySchemes[3].TokenURL, "http://example.com/token")
		}
		if len(design.Design.SecuritySchemes[3].Scopes) != 2 {
			t.Errorf("SecuritySchemes[3].Scopes = %v; want 2", design.Design.SecuritySchemes[3].Scopes)
		}
	})

	t.Run("should fallback properly to lower-level security", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.JWTSecurity("jwt", func() {
				apidsl.TokenURL("/token")
				apidsl.Scope("read", "Read")
				apidsl.Scope("write", "Write")
			})
			apidsl.BasicAuthSecurity("password")

			apidsl.Security("jwt")
		})
		apidsl.Resource("one", func() {
			apidsl.Action("first", func() {
				apidsl.Routing(apidsl.GET("/first"))
				apidsl.NoSecurity()
			})
			apidsl.Action("second", func() {
				apidsl.Routing(apidsl.GET("/second"))
			})
		})
		apidsl.Resource("two", func() {
			apidsl.Security("password")

			apidsl.Action("third", func() {
				apidsl.Routing(apidsl.GET("/third"))
			})
			apidsl.Action("fourth", func() {
				apidsl.Routing(apidsl.GET("/fourth"))
				apidsl.Security("jwt")
			})
		})
		apidsl.Resource("three", func() {
			apidsl.Action("fifth", func() {
				apidsl.Routing(apidsl.GET("/fifth"))
			})
		})
		apidsl.Resource("auth", func() {
			apidsl.NoSecurity()

			apidsl.Action("auth", func() {
				apidsl.Routing(apidsl.GET("/auth"))
			})
			apidsl.Action("refresh", func() {
				apidsl.Routing(apidsl.GET("/refresh"))
				apidsl.Security("jwt")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if len(design.Design.SecuritySchemes) != 2 {
			t.Errorf("SecuritySchemes = %v; want 2", design.Design.SecuritySchemes)
		}
		if v := design.Design.Resources["one"].Actions["first"].Security; v != nil {
			t.Errorf("Resources[one].Actions[first].Security = %v; want nil", v)
		}
		if v := design.Design.Resources["one"].Actions["second"].Security.Scheme.SchemeName; v != "jwt" {
			t.Errorf("Resources[one].Actions[second].Security.Scheme.SchemeName = %v; want jwt", v)
		}
		if v := design.Design.Resources["two"].Actions["third"].Security.Scheme.SchemeName; v != "password" {
			t.Errorf("Resources[two].Actions[third].Security.Scheme.SchemeName = %v; want password", v)
		}
		if v := design.Design.Resources["two"].Actions["fourth"].Security.Scheme.SchemeName; v != "jwt" {
			t.Errorf("Resources[two].Actions[fourth].Security.Scheme.SchemeName = %v; want jwt", v)
		}
		if v := design.Design.Resources["three"].Actions["fifth"].Security.Scheme.SchemeName; v != "jwt" {
			t.Errorf("Resources[three].Actions[fifth].Security.Scheme.SchemeName = %v; want jwt", v)
		}
		if v := design.Design.Resources["auth"].Actions["auth"].Security; v != nil {
			t.Errorf("Resources[auth].Actions[auth].Security = %v; want nil", v)
		}
		if v := design.Design.Resources["auth"].Actions["refresh"].Security.Scheme.SchemeName; v != "jwt" {
			t.Errorf("Resources[auth].Actions[refresh].Security.Scheme.SchemeName = %v; want jwt", v)
		}
	})
}

func TestSecurity_basic(t *testing.T) {
	t.Run("should fail because of duplicate in declaration", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.BasicAuthSecurity("broken_basic_authz", func() {
				apidsl.Description("desc")
				apidsl.Header("Authorization")
				apidsl.Query("access_token")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should fail because of invalid declaration of OAuth2Flow", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.BasicAuthSecurity("broken_basic_authz", func() {
				apidsl.Description("desc")
				apidsl.ImplicitFlow("invalid")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should fail because of invalid declaration of TokenURL", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.BasicAuthSecurity("broken_basic_authz", func() {
				apidsl.Description("desc")
				apidsl.TokenURL("/token")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("should fail because of invalid declaration of Header", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.BasicAuthSecurity("broken_basic_authz", func() {
				apidsl.Description("desc")
				apidsl.Header("invalid")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSecurity_OAuth2(t *testing.T) {
	t.Run("should pass with valid values when well defined", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.Host("example.com")
			apidsl.Scheme("http")
			apidsl.OAuth2Security("googAuthz", func() {
				apidsl.Description("Use Goog's Auth")
				apidsl.AccessCodeFlow("/auth", "/token")
				apidsl.Scope("scope:1", "Desc 1")
				apidsl.Scope("scope:2", "Desc 2")
			})
		})
		apidsl.Resource("one", func() {
			apidsl.Action("first", func() {
				apidsl.Routing(apidsl.GET("/first"))
				apidsl.Security("googAuthz", func() {
					apidsl.Scope("scope:1")
				})
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		scheme := design.Design.SecuritySchemes[0]
		if scheme.Description != "Use Goog's Auth" {
			t.Errorf("Description = %v; want Use Goog's Auth", scheme.Description)
		}
		if scheme.AuthorizationURL != "http://example.com/auth" {
			t.Errorf("AuthorizationURL = %v; want http://example.com/auth", scheme.AuthorizationURL)
		}
		if scheme.TokenURL != "http://example.com/token" {
			t.Errorf("TokenURL = %v; want http://example.com/token", scheme.TokenURL)
		}
		if scheme.Flow != "accessCode" {
			t.Errorf("Flow = %v; want accessCode", scheme.Flow)
		}
		if scheme.Scopes["scope:1"] != "Desc 1" {
			t.Errorf("Scopes[scope:1] = %v; want Desc 1", scheme.Scopes["scope:1"])
		}
		if scheme.Scopes["scope:2"] != "Desc 2" {
			t.Errorf("Scopes[scope:2] = %v; want Desc 2", scheme.Scopes["scope:2"])
		}
	})

	t.Run("should fail because of invalid declaration of Header", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("", func() {
			apidsl.OAuth2Security("googAuthz", func() {
				apidsl.Header("invalid")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
