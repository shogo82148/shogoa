package test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

// Global test definitions
const apiName = "API"
const apiDescription = "API description"
const resourceName = "R"
const resourceDescription = "R description"
const typeName = "T"
const typeDescription = "T description"
const mediaTypeIdentifier = "mt/json"
const mediaTypeDescription = "MT description"

var _ = apidsl.API(apiName, func() {
	apidsl.Description(apiDescription)
})

var _ = apidsl.Resource(resourceName, func() {
	apidsl.Description(resourceDescription)
})

var _ = apidsl.Type(typeName, func() {
	apidsl.Description(typeDescription)
	apidsl.Attribute("bar")
})

var _ = apidsl.MediaType(mediaTypeIdentifier, func() {
	apidsl.Description(mediaTypeDescription)
	apidsl.Attributes(func() { apidsl.Attribute("foo") })
	apidsl.View("default", func() { apidsl.Attribute("foo") })
})

func TestRunner(t *testing.T) {
	t.Run("with global DSL definitions", func(t *testing.T) {
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if design.Design == nil {
			t.Errorf("expected design, but got nil")
		}
		if design.Design.Name != apiName {
			t.Errorf("expected %s, but got %s", apiName, design.Design.Name)
		}
		if design.Design.Description != apiDescription {
			t.Errorf("expected %s, but got %s", apiDescription, design.Design.Description)
		}

		if _, ok := design.Design.Resources[resourceName]; !ok {
			t.Errorf("expected %s, but not found", resourceName)
		}
		res := design.Design.Resources[resourceName]
		if res.Name != resourceName {
			t.Errorf("expected %s, but got %s", resourceName, res.Name)
		}
		if res.Description != resourceDescription {
			t.Errorf("expected %s, but got %s", resourceDescription, res.Description)
		}

		if _, ok := design.Design.Types[typeName]; !ok {
			t.Errorf("expected %s, but not found", typeName)
		}
		typ := design.Design.Types[typeName]
		if typ.TypeName != typeName {
			t.Errorf("expected %s, but got %s", typeName, typ.TypeName)
		}
		if typ.Description != typeDescription {
			t.Errorf("expected %s, but got %s", typeDescription, typ.Description)
		}

		if _, ok := design.Design.MediaTypes[mediaTypeIdentifier]; !ok {
			t.Errorf("expected %s, but not found", mediaTypeIdentifier)
		}
		mt := design.Design.MediaTypes[mediaTypeIdentifier]
		if mt.Identifier != mediaTypeIdentifier {
			t.Errorf("expected %s, but got %s", mediaTypeIdentifier, mt.Identifier)
		}
		if mt.Description != mediaTypeDescription {
			t.Errorf("expected %s, but got %s", mediaTypeDescription, mt.Description)
		}
	})
}
