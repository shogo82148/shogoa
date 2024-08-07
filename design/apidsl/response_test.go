package apidsl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

var _ = Describe("Response", func() {
	var name string
	var dt design.DataType
	var dsl func()

	var res *design.ResponseDefinition

	BeforeEach(func() {
		dslengine.Reset()
		name = ""
		dsl = nil
		dt = nil
	})

	JustBeforeEach(func() {
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				if dt != nil {
					apidsl.Response(name, dt, dsl)
				} else {
					apidsl.Response(name, dsl)
				}
			})
		})
		dslengine.Run()
		if r, ok := design.Design.Resources["res"]; ok {
			if a, ok := r.Actions["action"]; ok {
				res = a.Responses[name]
			}
		}
	})

	Context("with no dsl and no name", func() {
		It("produces an invalid action definition", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).Should(HaveOccurred())
		})
	})

	Context("with no dsl", func() {
		BeforeEach(func() {
			name = "foo"
		})

		It("produces an invalid action definition", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).Should(HaveOccurred())
		})
	})

	Context("with a status", func() {
		const status = 201

		BeforeEach(func() {
			name = "foo"
			dsl = func() {
				apidsl.Status(status)
			}
		})

		It("produces a valid action definition and sets the status and parent", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).ShouldNot(HaveOccurred())
			Ω(res.Status).Should(Equal(status))
			Ω(res.Parent).ShouldNot(BeNil())
		})
	})

	Context("with a type override", func() {
		const status = 201

		BeforeEach(func() {
			name = "foo"
			dsl = func() {
				apidsl.Status(status)
			}
			dt = apidsl.HashOf(design.String, design.Any)
		})

		It("produces a response definition with the given type", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Type).Should(Equal(dt))
			Ω(res.Validate()).ShouldNot(HaveOccurred())
		})
	})

	Context("with a status and description", func() {
		const status = 201
		const description = "desc"

		BeforeEach(func() {
			name = "foo"
			dsl = func() {
				apidsl.Status(status)
				apidsl.Description(description)
			}
		})

		It("sets the status and description", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).ShouldNot(HaveOccurred())
			Ω(res.Status).Should(Equal(status))
			Ω(res.Description).Should(Equal(description))
		})
	})

	Context("with a status and name override", func() {
		const status = 201

		BeforeEach(func() {
			name = "foo"
			dsl = func() {
				apidsl.Status(status)
			}
		})

		It("sets the status and name", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).ShouldNot(HaveOccurred())
			Ω(res.Status).Should(Equal(status))
		})
	})

	Context("with a status and media type", func() {
		const status = 201
		const mediaType = "mt"

		BeforeEach(func() {
			name = "foo"
			dsl = func() {
				apidsl.Status(status)
				apidsl.Media(mediaType)
			}
		})

		It("sets the status and media type", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).ShouldNot(HaveOccurred())
			Ω(res.Status).Should(Equal(status))
			Ω(res.MediaType).Should(Equal(mediaType))
		})
	})

	Context("with a status and headers", func() {
		const status = 201
		const headerName = "Location"

		BeforeEach(func() {
			name = "foo"
			dsl = func() {
				apidsl.Status(status)
				apidsl.Headers(func() {
					apidsl.Header(headerName)
				})
			}
		})

		It("sets the status and headers", func() {
			Ω(res).ShouldNot(BeNil())
			Ω(res.Validate()).ShouldNot(HaveOccurred())
			Ω(res.Status).Should(Equal(status))
			Ω(res.Headers).ShouldNot(BeNil())
			Ω(res.Headers.Type).Should(BeAssignableToTypeOf(design.Object{}))
			o := res.Headers.Type.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(headerName))
		})
	})

	Context("not from the shogoa default definitions", func() {
		BeforeEach(func() {
			name = "foo"
		})

		It("does not set the Standard flag", func() {
			Ω(res.Standard).Should(BeFalse())
		})
	})

	Context("from the shogoa default definitions", func() {
		BeforeEach(func() {
			name = "Created"
		})

		It("sets the Standard flag", func() {
			Ω(res.Standard).Should(BeTrue())
		})
	})

})
