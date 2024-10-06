package apidsl_test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestResource(t *testing.T) {
	t.Run("with no dsl and no name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("", nil)
		if err := dslengine.Run(); err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with no dsl", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("foo", nil)
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run("with a description", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.Description("desc")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.Description != "desc" {
			t.Errorf("unexpected description: %s", res.Description)
		}
	})

	t.Run("with a parent resource that does not exist", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.Parent("parent")
		})
		if err := dslengine.Run(); err == nil {
			t.Errorf("expected an error, but no error")
		}

		if res.ParentName != "parent" {
			t.Errorf("unexpected parent name: %s", res.ParentName)
		}
		if err := res.Validate(); err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with actions", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.Action("action", func() { apidsl.Routing(apidsl.PUT("/:id")) })
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(res.Actions) != 1 {
			t.Errorf("unexpected actions: %v", res.Actions)
		}
		if _, ok := res.Actions["action"]; !ok {
			t.Errorf("action not found")
		}
	})

	t.Run("with metadata and actions", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.Metadata("swagger:generate", "false")
			apidsl.Action("action", func() { apidsl.Routing(apidsl.PUT("/:id")) })
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(res.Actions) != 1 {
			t.Errorf("unexpected actions: %v", res.Actions)
		}
		if _, ok := res.Actions["action"]; !ok {
			t.Errorf("action not found")
		}
		if res.Metadata["swagger:generate"][0] != "false" {
			t.Errorf("unexpected metadata: %v", res.Metadata)
		}
	})

	t.Run("with metadata and files", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.Metadata("swagger:generate", "false")
			apidsl.Files("path", "filename")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(res.FileServers) != 1 {
			t.Errorf("unexpected file servers: %v", res.FileServers)
		}
	})

	t.Run("with a canonical action that does not exist", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.CanonicalActionName("can")
		})
		if err := dslengine.Run(); err == nil {
			t.Errorf("expected an error, but no error")
		}

		if res.CanonicalActionName != "can" {
			t.Errorf("unexpected canonical action name: %s", res.CanonicalActionName)
		}
		if err := res.Validate(); err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with a canonical action that does exist", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.Action("can", func() { apidsl.Routing(apidsl.PUT("/:id")) })
			apidsl.CanonicalActionName("can")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.CanonicalActionName != "can" {
			t.Errorf("unexpected canonical action name: %s", res.CanonicalActionName)
		}
		if err := res.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run("with a base path", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.BasePath("basePath")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.BasePath != "basePath" {
			t.Errorf("unexpected base path: %s", res.BasePath)
		}
	})

	t.Run("with base params", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.BasePath("basePath/:paramID")
			apidsl.Params(func() {
				apidsl.Param("paramID")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.BasePath != "basePath/:paramID" {
			t.Errorf("unexpected base path: %s", res.BasePath)
		}
		if _, ok := res.Params.Type.ToObject()["paramID"]; !ok {
			t.Errorf("paramID not found")
		}
	})

	t.Run("with a media type name", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.DefaultMedia("application/mt")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.MediaType != "application/mt" {
			t.Errorf("unexpected media type: %s", res.MediaType)
		}
	})

	t.Run("with a view name", func(t *testing.T) {
		dslengine.Reset()
		res := apidsl.Resource("foo", func() {
			apidsl.DefaultMedia("application/mt", "compact")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.MediaType != "application/mt" {
			t.Errorf("unexpected media type: %s", res.MediaType)
		}
		if res.DefaultViewName != "compact" {
			t.Errorf("unexpected default view name: %s", res.DefaultViewName)
		}
	})

	t.Run("with an invalid media type", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("foo", func() {
			apidsl.DefaultMedia(&design.MediaTypeDefinition{Identifier: "application/foo"})
		})
		if err := dslengine.Run(); err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with a valid media type", func(t *testing.T) {
		dslengine.Reset()
		mt := &design.MediaTypeDefinition{
			UserTypeDefinition: &design.UserTypeDefinition{
				TypeName: "typeName",
			},
			Identifier: "application/vnd.raphael.shogoa.test",
		}
		res := apidsl.Resource("foo", func() {
			apidsl.DefaultMedia(mt)
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := res.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if res.MediaType != "application/vnd.raphael.shogoa.test" {
			t.Errorf("unexpected media type: %s", res.MediaType)
		}
	})

	t.Run("with a valid media type using a modifier", func(t *testing.T) {
		dslengine.Reset()
		mt := &design.MediaTypeDefinition{
			UserTypeDefinition: &design.UserTypeDefinition{
				TypeName: "typeName",
			},
			Identifier: "application/vnd.raphael.shogoa.test+json",
		}
		res := apidsl.Resource("foo", func() {
			apidsl.DefaultMedia(mt)
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if err := res.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if res.MediaType != "application/vnd.raphael.shogoa.test+json" {
			t.Errorf("unexpected media type: %s", res.MediaType)
		}
	})

	t.Run("with a trait that does not exist", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("foo", func() {
			apidsl.UseTrait("Authenticated")
		})
		if err := dslengine.Run(); err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with a trait that exists", func(t *testing.T) {
		dslengine.Reset()
		apidsl.API("test", func() {
			apidsl.Trait("descTrait", func() {
				apidsl.Description("desc")
			})
		})
		res := apidsl.Resource("foo", func() {
			apidsl.UseTrait("descTrait")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if res.Description != "desc" {
			t.Errorf("unexpected description: %s", res.Description)
		}
	})
}
