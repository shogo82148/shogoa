package apidsl_test

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestAction(t *testing.T) {
	t.Run("with only a name and a route", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("foo", nil)
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("with a name and DSL defining a route", func(t *testing.T) {
		dslengine.Reset()
		route := apidsl.GET("/:id")
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(route)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if action.Name != "foo" {
			t.Errorf("expected action name to be foo, got %s", action.Name)
		}
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if len(action.Routes) != 1 {
			t.Errorf("expected action to have a route, got %v", action.Routes)
		}
		if action.Routes[0] != route {
			t.Errorf("expected action to have a route, got %v", action.Routes)
		}
	})

	t.Run("with an empty params DSL", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Params(func() {})
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with a metadata", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id", func() {
					apidsl.Metadata("swagger:extension:x-get", `{"foo": "bar"}`)
				}))
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if action.Name != "foo" {
			t.Errorf("expected action name to be foo, got %s", action.Name)
		}
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if len(action.Routes) != 1 {
			t.Errorf("expected action to have a route, got %v", action.Routes)
		}
		want := dslengine.MetadataDefinition{
			"swagger:extension:x-get": []string{`{"foo": "bar"}`},
		}
		if diff := cmp.Diff(want, action.Routes[0].Metadata); diff != "" {
			t.Errorf("unexpected metadata: %s", diff)
		}
	})

	t.Run("with a string payload", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Payload(design.String)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if action.Payload == nil {
			t.Errorf("expected action to have a payload, got nil")
		}
		if action.Payload.Type != design.String {
			t.Errorf("expected action to have a string payload, got %s", action.Payload.Type)
		}
	})

	t.Run("with a name and DSL defining a description, route, headers, payload and responses", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("typeName", func() {
			apidsl.Attribute("name")
		})
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Description("description")
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Headers(func() { apidsl.Header("Foo") })
				apidsl.Payload("typeName")
				apidsl.Response(design.NoContent)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if action.Name != "foo" {
			t.Errorf("expected action name to be foo, got %s", action.Name)
		}
		if action.Description != "description" {
			t.Errorf("expected action description to be description, got %s", action.Description)
		}
		if len(action.Routes) != 1 {
			t.Errorf("expected action to have a route, got %v", action.Routes)
		}
		if len(action.Responses) != 1 {
			t.Errorf("expected action to have a response, got %v", action.Responses)
		}
		if _, ok := action.Responses["NoContent"]; !ok {
			t.Errorf("expected action to have a response, got %v", action.Responses)
		}
		headers := action.Headers.Type.(design.Object)
		if len(headers) != 1 {
			t.Errorf("expected action to have a header, got %v", headers)
		}
		if _, ok := headers["Foo"]; !ok {
			t.Errorf("expected action to have a header, got %v", headers)
		}
	})

	t.Run("with multiple headers sections", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("typeName", func() {
			apidsl.Attribute("name")
		})
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Headers(func() {
					apidsl.Header("Foo")
					apidsl.Required("Foo")
				})
				apidsl.Headers(func() {
					apidsl.Header("Foo2")
					apidsl.Required("Foo2")
				})
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if action.Name != "foo" {
			t.Errorf("expected action name to be foo, got %s", action.Name)
		}
		headers := action.Headers.Type.(design.Object)
		if len(headers) != 2 {
			t.Errorf("expected action to have two headers, got %v", headers)
		}
		if _, ok := headers["Foo"]; !ok {
			t.Errorf("expected action to have a header, got %v", headers)
		}
		if _, ok := headers["Foo2"]; !ok {
			t.Errorf("expected action to have a header, got %v", headers)
		}
		want := []string{"Foo", "Foo2"}
		if diff := cmp.Diff(want, action.Headers.Validation.Required); diff != "" {
			t.Errorf("unexpected required headers: %s", diff)
		}
	})

	t.Run("using a response with a media type modifier", func(t *testing.T) {
		dslengine.Reset()
		apidsl.MediaType("application/vnd.app.foo+json", func() {
			apidsl.Attributes(func() { apidsl.Attribute("foo") })
			apidsl.View("default", func() { apidsl.Attribute("foo") })
		})
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Response(design.OK, "application/vnd.app.foo+json")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if len(action.Responses) != 1 {
			t.Errorf("expected action to have a response, got %v", action.Responses)
		}
		resp := action.Responses["OK"]
		if resp.MediaType != "application/vnd.app.foo+json" {
			t.Errorf("expected response to have a media type, got %s", resp.MediaType)
		}
	})

	t.Run("using a response template", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("test", func() {
			apidsl.ResponseTemplate("tmpl", func(status, name string) {
				st, err := strconv.Atoi(status)
				if err != nil {
					dslengine.ReportError(err.Error())
					return
				}
				apidsl.Status(st)
			})
		})
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Response("tmpl", "200", "respName", func() {
					apidsl.Media("media")
				})
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["res"].Actions["foo"]
		if err := action.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if len(action.Responses) != 1 {
			t.Errorf("expected action to have a response, got %v", action.Responses)
		}
		resp := action.Responses["tmpl"]
		if resp.Name != "tmpl" {
			t.Errorf("expected response to have a name, got %s", resp.Name)
		}
		if resp.Status != 200 {
			t.Errorf("expected response to have a status, got %d", resp.Status)
		}
		if resp.MediaType != "media" {
			t.Errorf("expected response to have a media type, got %s", resp.MediaType)
		}
	})

	t.Run("using a response template, called incorrectly", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("test", func() {
			apidsl.ResponseTemplate("tmpl", func(status, name string) {
				st, err := strconv.Atoi(status)
				if err != nil {
					dslengine.ReportError(err.Error())
					return
				}
				apidsl.Status(st)
			})
		})
		apidsl.Resource("res", func() {
			apidsl.Action("foo", func() {
				apidsl.Routing(apidsl.GET("/:id"))
				apidsl.Response("tmpl", "not an integer", "respName", func() {
					apidsl.Media("media")
				})
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected error")
		}
	})
}

func TestPayload(t *testing.T) {
	t.Run("with a payload definition", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("foo", func() {
			apidsl.Action("bar", func() {
				apidsl.Routing(apidsl.GET(""))
				apidsl.Payload(func() {
					apidsl.Member("name")
					apidsl.Required("name")
				})
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["foo"].Actions["bar"]
		if action.Payload == nil {
			t.Error("expected payload to be defined")
		}
	})

	t.Run("with an array", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("foo", func() {
			apidsl.Action("bar", func() {
				apidsl.Routing(apidsl.GET(""))
				apidsl.Payload(apidsl.ArrayOf(design.Integer))
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["foo"].Actions["bar"]
		if !action.Payload.IsArray() {
			t.Error("expected payload to be an array")
		}
		array := action.Payload.ToArray()
		if array.ElemType.Type != design.Integer {
			t.Errorf("expected payload to have an integer element type, got %s", array.ElemType.Type)
		}
	})

	t.Run("with a hash", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("foo", func() {
			apidsl.Action("bar", func() {
				apidsl.Routing(apidsl.GET(""))
				apidsl.Payload(apidsl.HashOf(design.String, design.Integer))
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		action := design.Design.Resources["foo"].Actions["bar"]
		if !action.Payload.IsHash() {
			t.Error("expected payload to be a hash")
		}
		hash := action.Payload.ToHash()
		if hash.KeyType.Type != design.String {
			t.Errorf("expected payload to have a string key type, got %s", hash.KeyType.Type)
		}
		if hash.ElemType.Type != design.Integer {
			t.Errorf("expected payload to have an integer element type, got %s", hash.ElemType.Type)
		}
	})
}
