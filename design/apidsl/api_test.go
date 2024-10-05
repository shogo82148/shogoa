package apidsl_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestAPI(t *testing.T) {
	t.Run("with no dsl", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", nil)
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Name != "foo" {
			t.Errorf("Name = %q; want %q", design.Design.Name, "foo")
		}
	})

	t.Run("with an already defined API with the same name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", nil)
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
		apidsl.API("foo", nil)

		if dslengine.Errors == nil {
			t.Error("Errors = nil; want an error")
		}
	})

	t.Run("with an already defined API with a different name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", nil)
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
		apidsl.API("news", nil)

		if dslengine.Errors == nil {
			t.Error("Errors = nil; want an error")
		}
	})

	t.Run("with a description", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Description("bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Description != "bar" {
			t.Errorf("Description = %q; want %q", design.Design.Description, "bar")
		}
	})

	t.Run("with a title", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Title("bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Title != "bar" {
			t.Errorf("Title = %q; want %q", design.Design.Title, "bar")
		}
	})

	t.Run("with a version", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Version("bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Version != "bar" {
			t.Errorf("Version = %q; want %q", design.Design.Version, "bar")
		}
	})

	t.Run("with a terms of service", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.TermsOfService("bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.TermsOfService != "bar" {
			t.Errorf("TermsOfService = %q; want %q", design.Design.TermsOfService, "bar")
		}
	})

	t.Run("with contact information", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Contact(func() {
				apidsl.Name("contactName")
				apidsl.Email("contactEmail")
				apidsl.URL("http://example.com")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Contact == nil {
			t.Error("Contact = nil; want not nil")
		}
		if design.Design.Contact.Name != "contactName" {
			t.Errorf("Contact.Name = %q; want %q", design.Design.Contact.Name, "contactName")
		}
		if design.Design.Contact.Email != "contactEmail" {
			t.Errorf("Contact.Email = %q; want %q", design.Design.Contact.Email, "contactEmail")
		}
		if design.Design.Contact.URL != "http://example.com" {
			t.Errorf("Contact.URL = %q; want %q", design.Design.Contact.URL, "http://example.com")
		}
	})

	t.Run("with license information", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.License(func() {
				apidsl.Name("licenseName")
				apidsl.URL("http://example.com")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.License == nil {
			t.Error("License = nil; want not nil")
		}
		if design.Design.License.Name != "licenseName" {
			t.Errorf("License.Name = %q; want %q", design.Design.License.Name, "licenseName")
		}
		if design.Design.License.URL != "http://example.com" {
			t.Errorf("License.URL = %q; want %q", design.Design.License.URL, "http://example.com")
		}
	})

	t.Run("with Consumes", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Consumes("application/json")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if len(design.Design.Consumes) != 1 {
			t.Errorf("Consumes = %d; want 1", len(design.Design.Consumes))
		}
		if design.Design.Consumes[0].MIMETypes[0] != "application/json" {
			t.Errorf("Consumes[0].MIMETypes[0] = %q; want %q", design.Design.Consumes[0].MIMETypes[0], "application/json")
		}
		if design.Design.Consumes[0].PackagePath != "" {
			t.Errorf("Consumes[0].PackagePath = %q; want %q", design.Design.Consumes[0].PackagePath, "")
		}
	})

	t.Run("using a custom encoding package", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Consumes("application/json", func() {
				apidsl.Package("github.com/shogo82148/shogoa/encoding/json")
				apidsl.Function("NewFoo")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if len(design.Design.Consumes) != 1 {
			t.Errorf("Consumes = %d; want 1", len(design.Design.Consumes))
		}
		if design.Design.Consumes[0].MIMETypes[0] != "application/json" {
			t.Errorf("Consumes[0].MIMETypes[0] = %q; want %q", design.Design.Consumes[0].MIMETypes[0], "application/json")
		}
		if design.Design.Consumes[0].PackagePath != "github.com/shogo82148/shogoa/encoding/json" {
			t.Errorf("Consumes[0].PackagePath = %q; want %q", design.Design.Consumes[0].PackagePath, "github.com/shogo82148/shogoa/encoding/json")
		}
		if design.Design.Consumes[0].Function != "NewFoo" {
			t.Errorf("Consumes[0].Function = %q; want %q", design.Design.Consumes[0].Function, "NewFoo")
		}
	})

	t.Run("with a BasePath", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.BasePath("/foo")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.BasePath != "/foo" {
			t.Errorf("BasePath = %q; want %q", design.Design.BasePath, "/foo")
		}
	})

	t.Run("with Params", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Params(func() {
				apidsl.Param("accountID", design.Integer, "the account ID")
				apidsl.Param("id", design.String, "the widget ID")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Params.Type == nil {
			t.Error("Params.Type = nil; want not nil")
		}
		params := design.Design.Params.Type.ToObject()
		if len(params) != 2 {
			t.Errorf("Params = %d; want 2", len(params))
		}
		accountID := params["accountID"]
		if accountID.Type != design.Integer {
			t.Errorf("Params[accountID].Type = %v; want %v", accountID.Type, design.Integer)
		}
		if accountID.Description != "the account ID" {
			t.Errorf("Params[accountID].Description = %q; want %q", accountID.Description, "the account ID")
		}
		id := params["id"]
		if id.Type != design.String {
			t.Errorf("Params[id].Type = %v; want %v", id.Type, design.String)
		}
		if id.Description != "the widget ID" {
			t.Errorf("Params[id].Description = %q; want %q", id.Description, "the widget ID")
		}
	})

	t.Run("with Params and a BasePath using them", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Params(func() {
				apidsl.Param("accountID", design.Integer, "the account ID")
				apidsl.Param("id", design.String, "the widget ID")
			})
			apidsl.BasePath("/:accountID/:id")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if design.Design.Params.Type == nil {
			t.Error("Params.Type = nil; want not nil")
		}
		params := design.Design.Params.Type.ToObject()
		if len(params) != 2 {
			t.Errorf("Params = %d; want 2", len(params))
		}
		accountID := params["accountID"]
		if accountID.Type != design.Integer {
			t.Errorf("Params[accountID].Type = %v; want %v", accountID.Type, design.Integer)
		}
		if accountID.Description != "the account ID" {
			t.Errorf("Params[accountID].Description = %q; want %q", accountID.Description, "the account ID")
		}
		id := params["id"]
		if id.Type != design.String {
			t.Errorf("Params[id].Type = %v; want %v", id.Type, design.String)
		}
		if design.Design.BasePath != "/:accountID/:id" {
			t.Errorf("BasePath = %q; want %q", design.Design.BasePath, "/:accountID/:id")
		}
	})

	t.Run("with conflicting resource and API base params", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Params(func() {
				apidsl.Param("accountID", design.Integer, "the account ID")
			})
			apidsl.BasePath("/:accountID")
		})
		apidsl.Resource("bar", func() {
			apidsl.BasePath("/:accountID")
		})
		if err := dslengine.Run(); err == nil {
			t.Error("Run() = nil; want an error")
		}
	})

	t.Run("with an absolute resource base path", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Params(func() {
				apidsl.Param("accountID", design.Integer, "the account ID")
			})
			apidsl.BasePath("/:accountID")
		})
		apidsl.Resource("bar", func() {
			apidsl.Params(func() {
				apidsl.Param("accountID", design.Integer, "the account ID")
			})
			apidsl.BasePath("//:accountID")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with ResponseTemplates", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.ResponseTemplate("NotFound2", func() {
				apidsl.Description("Resource Not Found")
				apidsl.Status(404)
				apidsl.Media("text/plain")
			})
			apidsl.ResponseTemplate("OK", func(mt string) {
				apidsl.Description("All good")
				apidsl.Status(200)
				apidsl.Media(mt)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		want := design.ResponseDefinition{
			Name:        "NotFound2",
			Description: "Resource Not Found",
			Status:      404,
			MediaType:   "text/plain",
		}
		got := *design.Design.Responses["NotFound2"]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Responses[NotFound2] (-want +got):\n%s", diff)
		}

		if len(design.Design.ResponseTemplates) != 1 {
			t.Errorf("ResponseTemplates = %d; want 1", len(design.Design.ResponseTemplates))
		}
		if design.Design.ResponseTemplates["OK"] == nil {
			t.Error("ResponseTemplates[OK] = nil; want not nil")
		}
	})

	t.Run("with Traits", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Trait("Authenticated", func() {
				apidsl.Headers(func() {
					apidsl.Header("Auth-Token")
					apidsl.Required("Auth-Token")
				})
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if len(design.Design.Traits) != 1 {
			t.Errorf("Traits = %d; want 1", len(design.Design.Traits))
		}
		if design.Design.Traits["Authenticated"] == nil {
			t.Error("Traits[Authenticated] = nil; want not nil")
		}
	})

	t.Run("using Traits", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Trait("Authenticated", func() {
				apidsl.Attributes(func() {
					apidsl.Attribute("foo")
				})
			})
		})
		apidsl.MediaType("application/vnd.foo", func() {
			apidsl.UseTrait("Authenticated")
			apidsl.Attributes(func() {
				apidsl.Attribute("bar")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("bar")
				apidsl.Attribute("foo")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if len(design.Design.Traits) != 1 {
			t.Errorf("Traits = %d; want 1", len(design.Design.Traits))
		}
		if design.Design.Traits["Authenticated"] == nil {
			t.Error("Traits[Authenticated] = nil; want not nil")
		}
		if len(design.Design.MediaTypes) != 1 {
			t.Errorf("MediaTypes = %d; want 1", len(design.Design.MediaTypes))
		}
		foo := design.Design.MediaTypes["application/vnd.foo"]
		if foo.Type.ToObject() == nil {
			t.Error("Type.ToObject() = nil; want not nil")
		}
		o := foo.Type.ToObject()
		if _, ok := o["foo"]; !ok {
			t.Error("Type.ToObject() does not have foo")
		}
		if _, ok := o["bar"]; !ok {
			t.Error("Type.ToObject() does not have bar")
		}
	})

	t.Run("using variadic Traits", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("foo", func() {
			apidsl.Trait("Authenticated", func() {
				apidsl.Attributes(func() {
					apidsl.Attribute("foo")
				})
			})
			apidsl.Trait("AuthenticatedTwo", func() {
				apidsl.Attributes(func() {
					apidsl.Attribute("baz")
				})
			})
		})
		apidsl.MediaType("application/vnd.foo", func() {
			apidsl.UseTrait("Authenticated", "AuthenticatedTwo")
			apidsl.Attributes(func() {
				apidsl.Attribute("bar")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("bar")
				apidsl.Attribute("foo")
				apidsl.Attribute("baz")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := design.Design.Validate(); err != nil {
			t.Error(err)
		}
		if len(design.Design.Traits) != 2 {
			t.Errorf("Traits = %d; want 2", len(design.Design.Traits))
		}
		if _, ok := design.Design.Traits["Authenticated"]; !ok {
			t.Error("Traits does not have Authenticated")
		}
		if _, ok := design.Design.Traits["AuthenticatedTwo"]; !ok {
			t.Error("Traits does not have AuthenticatedTwo")
		}
		foo := design.Design.MediaTypes["application/vnd.foo"]
		o := foo.Type.ToObject()
		if _, ok := o["foo"]; !ok {
			t.Error("Type.ToObject() does not have foo")
		}
		if _, ok := o["bar"]; !ok {
			t.Error("Type.ToObject() does not have bar")
		}
		if _, ok := o["baz"]; !ok {
			t.Error("Type.ToObject() does not have baz")
		}
	})
}
