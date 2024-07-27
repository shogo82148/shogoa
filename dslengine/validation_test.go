package dslengine_test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestValidation(t *testing.T) {
	t.Run("a valid enum validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Enum("red", "blue")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		values := att.Validation.Values
		if len(values) != 2 {
			t.Errorf("unexpected values: %v", values)
		}
		if values[0] != "red" {
			t.Errorf("unexpected value: %v", values[0])
		}
		if values[1] != "blue" {
			t.Errorf("unexpected value: %v", values[1])
		}
	})

	t.Run("with an incompatible enum validation type", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.Integer, func() {
				apidsl.Enum(1, "blue")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid format validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Format("email")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		if att.Validation.Format != "email" {
			t.Errorf("unexpected format: %v", att.Validation.Format)
		}
	})

	t.Run("with an invalid format validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Format("emailz")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid pattern validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Pattern("^foo$")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		if att.Validation.Pattern != "^foo$" {
			t.Errorf("unexpected pattern: %v", att.Validation.Pattern)
		}
	})

	t.Run("with an invalid pattern validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Pattern("[invalid")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with an invalid format validation type", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.Integer, func() {
				apidsl.Format("email")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid min value validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.Integer, func() {
				apidsl.Minimum(2)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		if *att.Validation.Minimum != 2 {
			t.Errorf("unexpected minimum: %v", *att.Validation.Minimum)
		}
	})

	t.Run("with an invalid min value validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Minimum(2)
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid max value validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.Integer, func() {
				apidsl.Maximum(2)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		if *att.Validation.Maximum != 2 {
			t.Errorf("unexpected maximum: %v", *att.Validation.Maximum)
		}
	})

	t.Run("with an invalid max value validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String, func() {
				apidsl.Maximum(2)
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid min length validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", apidsl.ArrayOf(design.Integer), func() {
				apidsl.MinLength(2)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		if *att.Validation.MinLength != 2 {
			t.Errorf("unexpected minimum length: %v", *att.Validation.MinLength)
		}
	})

	t.Run("with an invalid min length validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.Integer, func() {
				apidsl.MinLength(2)
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid max length validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", apidsl.ArrayOf(design.Integer), func() {
				apidsl.MaxLength(2)
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		o := design.Design.Types["bar"].Type.(design.Object)
		att := o["attName"]
		if *att.Validation.MaxLength != 2 {
			t.Errorf("unexpected maximum length: %v", *att.Validation.MaxLength)
		}
	})

	t.Run("with an invalid max length validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.Integer, func() {
				apidsl.MaxLength(2)
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("with a valid required field validation", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			apidsl.Attribute("attName", design.String)
			apidsl.Required("attName")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		required := design.Design.Types["bar"].Validation.Required
		if len(required) != 1 {
			t.Errorf("unexpected required fields: %v", required)
		}
		if required[0] != "attName" {
			t.Errorf("unexpected required field: %v", required[0])
		}
	})
}
