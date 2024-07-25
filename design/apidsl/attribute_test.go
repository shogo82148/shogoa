package apidsl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

// TestCD is a test container definition.
type TestCD struct {
	*design.AttributeDefinition
}

// Attribute returns a dummy attribute.
func (t *TestCD) Attribute() *design.AttributeDefinition {
	return t.AttributeDefinition
}

// DSL implements Source
func (t *TestCD) DSL() func() {
	return func() {
		apidsl.Attribute("foo")
	}
}

// Context implement Definition
func (t *TestCD) Context() string {
	return "test"
}

// DSLName returns the DSL name.
func (t *TestCD) DSLName() string {
	return "TestCD"
}

// DependsOn returns the DSL dependencies.
func (t *TestCD) DependsOn() []dslengine.Root {
	return nil
}

// IterateSets implement Root
func (t *TestCD) IterateSets(it dslengine.SetIterator) {
	it([]dslengine.Definition{t})
}

// Reset is a no-op
func (t *TestCD) Reset() {}

var _ = Describe("ContainerDefinition", func() {
	var att *design.AttributeDefinition
	var testCD *TestCD
	BeforeEach(func() {
		dslengine.Reset()
		att = &design.AttributeDefinition{Type: design.Object{}}
		testCD = &TestCD{AttributeDefinition: att}
		dslengine.Register(testCD)
	})

	JustBeforeEach(func() {
		err := dslengine.Run()
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("contains attributes", func() {
		Ω(testCD.Attribute()).Should(Equal(att))
	})
})

var _ = Describe("Attribute", func() {
	var name string
	var dataType interface{}
	var description string
	var dsl func()

	var parent *design.AttributeDefinition

	BeforeEach(func() {
		dslengine.Reset()
		name = ""
		dataType = nil
		description = ""
		dsl = nil
	})

	JustBeforeEach(func() {
		apidsl.Type("type", func() {
			if dsl == nil {
				if dataType == nil {
					apidsl.Attribute(name)
				} else if description == "" {
					apidsl.Attribute(name, dataType)
				} else {
					apidsl.Attribute(name, dataType, description)
				}
			} else if dataType == nil {
				apidsl.Attribute(name, dsl)
			} else if description == "" {
				apidsl.Attribute(name, dataType, dsl)
			} else {
				apidsl.Attribute(name, dataType, description, dsl)
			}
		})
		dslengine.Run()
		if t, ok := design.Design.Types["type"]; ok {
			parent = t.AttributeDefinition
		}
	})

	Context("with only a name", func() {
		BeforeEach(func() {
			name = "foo"
		})

		It("produces an attribute of type string", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.String))
		})
	})

	Context("with a name and datatype", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.Integer
		})

		It("produces an attribute of given type", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.Integer))
		})
	})

	Context("with a name and uuid datatype", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.UUID
		})

		It("produces an attribute of uuid type", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.UUID))
		})
	})

	Context("with a name and date datatype", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.DateTime
		})

		It("produces an attribute of date type", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.DateTime))
		})
	})

	Context("with a name, datatype and description", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.Integer
			description = "bar"
		})

		It("produces an attribute of given type and given description", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.Integer))
			Ω(o[name].Description).Should(Equal(description))
		})
	})

	Context("with a name and a DSL defining a 'readOnly' attribute", func() {
		BeforeEach(func() {
			name = "foo"
			dsl = func() { apidsl.ReadOnly() }
		})

		It("produces an attribute of type string set to readOnly", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.String))
			Ω(o[name].IsReadOnly()).Should(BeTrue())
		})
	})

	Context("with a name and a DSL defining an enum validation", func() {
		BeforeEach(func() {
			name = "foo"
			dsl = func() { apidsl.Enum("one", "two") }
		})

		It("produces an attribute of type string with a validation", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.String))
			Ω(o[name].Validation).ShouldNot(BeNil())
			Ω(o[name].Validation.Values).Should(Equal([]interface{}{"one", "two"}))
		})
	})

	Context("with a name, type datetime and a DSL defining a default value", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.DateTime
			dsl = func() { apidsl.Default("1978-06-30T10:00:00+09:00") }
		})

		It("produces an attribute of type string with a default value", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.DateTime))
			Ω(o[name].Validation).Should(BeNil())
			Ω(o[name].DefaultValue).Should(Equal(interface{}("1978-06-30T10:00:00+09:00")))
		})
	})

	Context("with a name, type integer and a DSL defining an enum validation", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.Integer
			dsl = func() { apidsl.Enum(1, 2) }
		})

		It("produces an attribute of type integer with a validation", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.Integer))
			Ω(o[name].Validation).ShouldNot(BeNil())
			Ω(o[name].Validation.Values).Should(Equal([]interface{}{1, 2}))
		})
	})

	Context("with a name, type integer, a description and a DSL defining an enum validation", func() {
		BeforeEach(func() {
			name = "foo"
			dataType = design.String
			description = "bar"
			dsl = func() { apidsl.Enum("one", "two") }
		})

		It("produces an attribute of type integer with a validation and the description", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.String))
			Ω(o[name].Validation).ShouldNot(BeNil())
			Ω(o[name].Validation.Values).Should(Equal([]interface{}{"one", "two"}))
		})
	})

	Context("with a name and type uuid", func() {
		BeforeEach(func() {
			name = "birthdate"
			dataType = design.UUID
		})

		It("produces an attribute of type date with a validation and the description", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.UUID))
		})
	})

	Context("with a name and type date", func() {
		BeforeEach(func() {
			name = "birthdate"
			dataType = design.DateTime
		})

		It("produces an attribute of type date with a validation and the description", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(design.DateTime))
		})
	})

	Context("with a name and a type defined by name", func() {
		var Foo *design.UserTypeDefinition

		BeforeEach(func() {
			name = "fooatt"
			dataType = "foo"
			Foo = apidsl.Type("foo", func() {
				apidsl.Attribute("bar")
			})
		})

		It("produces an attribute of the corresponding type", func() {
			t := parent.Type
			Ω(t).ShouldNot(BeNil())
			Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
			o := t.(design.Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(name))
			Ω(o[name].Type).Should(Equal(Foo))
		})
	})

	Context("with child attributes", func() {
		const childAtt = "childAtt"

		BeforeEach(func() {
			name = "foo"
			dsl = func() { apidsl.Attribute(childAtt) }
		})

		Context("on an attribute that is not an object", func() {
			BeforeEach(func() {
				dataType = design.Integer
			})

			It("fails", func() {
				Ω(dslengine.Errors).Should(HaveOccurred())
			})
		})

		Context("on an attribute that does not have a type", func() {
			It("sets the type to Object", func() {
				t := parent.Type
				Ω(t).ShouldNot(BeNil())
				Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
				o := t.(design.Object)
				Ω(o).Should(HaveLen(1))
				Ω(o).Should(HaveKey(name))
				Ω(o[name].Type).Should(BeAssignableToTypeOf(design.Object{}))
			})
		})

		Context("on an attribute of type Object", func() {
			BeforeEach(func() {
				dataType = design.Object{}
			})

			It("initializes the object attributes", func() {
				Ω(dslengine.Errors).ShouldNot(HaveOccurred())
				t := parent.Type
				Ω(t).ShouldNot(BeNil())
				Ω(t).Should(BeAssignableToTypeOf(design.Object{}))
				o := t.(design.Object)
				Ω(o).Should(HaveLen(1))
				Ω(o).Should(HaveKey(name))
				Ω(o[name].Type).Should(BeAssignableToTypeOf(design.Object{}))
				co := o[name].Type.(design.Object)
				Ω(co).Should(HaveLen(1))
				Ω(co).Should(HaveKey(childAtt))
			})
		})
	})
})
