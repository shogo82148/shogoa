package apidsl

import (
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
)

// API implements the top level API DSL. It defines the API name, default description and other
// default global property values.
func API(name string, dsl func()) *design.APIDefinition {
	if design.Design.Name != "" {
		dslengine.ReportError("multiple API definitions, only one is allowed")
		return nil
	}

	if name == "" {
		dslengine.ReportError("API name cannot be empty")
	}
	design.Design.Name = name
	design.Design.DSLFunc = dsl
	return design.Design
}
