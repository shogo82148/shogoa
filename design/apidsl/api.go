package apidsl

import "github.com/shogo82148/shogoa/design"

// API implements the top level API DSL. It defines the API name, default description and other
// default global property values.
func API(name string, dsl func()) *design.APIDefinition {
	if design.Design.Name != "" {
		// report error
		panic("multiple API definitions, only one is allowed")
	}

	if name == "" {
		// report error
		panic("API name cannot be empty")
	}
	design.Design.Name = name
	design.Design.DSLFunc = dsl
	return design.Design
}
