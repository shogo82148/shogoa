package design_test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
)

func TestSecuritySchemeDefinition(t *testing.T) {
	t.Run("with valid token and authorization URLs", func(t *testing.T) {
		def := &design.SecuritySchemeDefinition{
			TokenURL:         "http://example.com/token",
			AuthorizationURL: "http://example.com/auth",
		}
		if err := def.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
	})

	t.Run("with an invalid token URL", func(t *testing.T) {
		def := &design.SecuritySchemeDefinition{
			TokenURL:         ":",
			AuthorizationURL: "http://example.com/auth",
		}
		if err := def.Validate(); err == nil {
			t.Error("Validate() = nil; want an error")
		}
	})

	t.Run("with an absolute token URL", func(t *testing.T) {
		def := &design.SecuritySchemeDefinition{
			TokenURL: "http://example.com/token",
		}
		def.Finalize()
		if def.TokenURL != "http://example.com/token" {
			t.Errorf("Finalize() = %q; want %q", def.TokenURL, "http://example.com/token")
		}
	})

	t.Run("with a relative token URL", func(t *testing.T) {
		design.Design.Schemes = []string{"http"}
		design.Design.Host = "example.com"
		def := &design.SecuritySchemeDefinition{
			TokenURL: "/token",
		}
		def.Finalize()
		if def.TokenURL != "http://example.com/token" {
			t.Errorf("Finalize() = %q; want %q", def.TokenURL, "http://example.com/token")
		}
	})

	t.Run("with an invalid authorization URL", func(t *testing.T) {
		def := &design.SecuritySchemeDefinition{
			TokenURL:         "http://example.com/token",
			AuthorizationURL: ":",
		}
		if err := def.Validate(); err == nil {
			t.Error("Validate() = nil; want an error")
		}
	})

	t.Run("with an absolute authorization URL", func(t *testing.T) {
		def := &design.SecuritySchemeDefinition{
			AuthorizationURL: "http://example.com/auth",
		}
		def.Finalize()
		if def.AuthorizationURL != "http://example.com/auth" {
			t.Errorf("Finalize() = %q; want %q", def.AuthorizationURL, "http://example.com/auth")
		}
	})

	t.Run("with a relative authorization URL", func(t *testing.T) {
		design.Design.Schemes = []string{"http"}
		design.Design.Host = "example.com"
		def := &design.SecuritySchemeDefinition{
			AuthorizationURL: "/auth",
		}
		def.Finalize()
		if def.AuthorizationURL != "http://example.com/auth" {
			t.Errorf("Finalize() = %q; want %q", def.AuthorizationURL, "http://example.com/auth")
		}
	})
}
