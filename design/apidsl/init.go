package apidsl

import "github.com/shogo82148/shogoa/design"

// Setup API DSL roots.
func init() {
	design.Design = design.NewAPIDefinition()
	// TODO
	// design.GeneratedMediaTypes = make(design.MediaTypeRoot)
	// design.ProjectedMediaTypes = make(design.MediaTypeRoot)
	// dslengine.Register(design.Design)
	// dslengine.Register(design.GeneratedMediaTypes)
}
