package design

import (
	. "github.com/shogo82148/shogoa/design"
	. "github.com/shogo82148/shogoa/design/apidsl"
)

// This is the cellar application API design used by shogoa to generate
// the application code, client, tests, documentation etc.
var _ = API("cellar", func() {
	Title("The virtual wine cellar")
	Description("A basic example of a CRUD API implemented with shogoa")
	Contact(func() {
		Name("shogoa team")
		Email("admin@shogoa.design")
		URL("http://shogoa.design")
	})
	License(func() {
		Name("MIT")
		URL("https://github.com/shogo82148/shogoa/blob/master/LICENSE")
	})
	Docs(func() {
		Description("shogoa guide")
		URL("http://shogoa.design/getting-started.html")
	})
	Host("localhost:8081")
	Scheme("http")
	BasePath("/cellar")

	Origin("http://swagger.shogoa.design", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		MaxAge(600)
		Credentials()
	})

	ResponseTemplate(Created, func(pattern string) {
		Description("Resource created")
		Status(201)
		Headers(func() {
			Header("Location", String, "href to created resource", func() {
				Pattern(pattern)
			})
		})
	})
})
