package design

import (
	. "github.com/shogo82148/shogoa/design"
	. "github.com/shogo82148/shogoa/design/apidsl"
)

var _ = API("adder", func() {
	Title("The adder API")
	Description("A teaser for shogoa")
	Host("localhost:8080")
	Scheme("http")
})

var _ = Resource("operands", func() {
	Action("add", func() {
		Routing(GET("add/:left/:right"))
		Description("add returns the sum of the left and right parameters in the response body")
		Params(func() {
			Param("left", Integer, "Left operand")
			Param("right", Integer, "Right operand")
		})
		Response(OK, "text/plain")
	})

})
