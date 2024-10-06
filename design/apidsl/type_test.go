package apidsl_test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestType(t *testing.T) {
	t.Run("with no dsl and no name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("", nil)
		if err := dslengine.Run(); err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with no dsl", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("foo", nil)
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run("with attributes", func(t *testing.T) {
		dslengine.Reset()
		ut := apidsl.Type("foo", func() {
			apidsl.Attribute("att")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if ut.AttributeDefinition == nil {
			t.Error("AttributeDefinition is nil")
		}
	})

	t.Run("with a name and uuid datatype", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("foo", func() {
			apidsl.Attribute("att", design.UUID)
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		ut := design.Design.Types["foo"]
		if err := ut.Validate("test", design.Design); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if ut.AttributeDefinition == nil {
			t.Error("AttributeDefinition is nil")
		}
		if ut.Type == nil {
			t.Error("Type is nil")
		}
		if _, ok := ut.Type.(design.Object); !ok {
			t.Errorf("Type is %T; want design.Object", ut.Type)
		}
		o := ut.Type.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if _, ok := o["att"]; !ok {
			t.Error("att is not found")
		}
		if o["att"].Type != design.UUID {
			t.Errorf("att.Type = %v; want %v", o["att"].Type, design.UUID)
		}
	})

	t.Run("with a name and date datatype", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("foo", func() {
			apidsl.Attribute("att", design.DateTime)
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		ut := design.Design.Types["foo"]
		if err := ut.Validate("test", design.Design); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if ut.AttributeDefinition == nil {
			t.Error("AttributeDefinition is nil")
		}
		if ut.Type == nil {
			t.Error("Type is nil")
		}
		if _, ok := ut.Type.(design.Object); !ok {
			t.Errorf("Type is %T; want design.Object", ut.Type)
		}
		o := ut.Type.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if _, ok := o["att"]; !ok {
			t.Error("att is not found")
		}
		if o["att"].Type != design.DateTime {
			t.Errorf("att.Type = %v; want %v", o["att"].Type, design.DateTime)
		}
	})
}

func TestArrayOf(t *testing.T) {
	t.Run("used on a global variable", func(t *testing.T) {
		dslengine.Reset()
		ut := apidsl.Type("example", func() {
			apidsl.Attribute("id")
		})
		ar := apidsl.ArrayOf(ut)
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if ar.Kind() != design.ArrayKind {
			t.Errorf("ar.Kind() = %v; want %v", ar.Kind(), design.ArrayKind)
		}
		if ar.ElemType.Type != ut {
			t.Errorf("ar.ElemType.Type = %v; want %v", ar.ElemType.Type, ut)
		}
	})

	t.Run("with a DSL", func(t *testing.T) {
		dslengine.Reset()
		ar := apidsl.ArrayOf(design.String, func() {
			apidsl.Pattern("foo")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if ar.Kind() != design.ArrayKind {
			t.Errorf("ar.Kind() = %v; want %v", ar.Kind(), design.ArrayKind)
		}
		if ar.ElemType.Type != design.String {
			t.Errorf("ar.ElemType.Type = %v; want %v", ar.ElemType.Type, design.String)
		}
		if ar.ElemType.Validation.Pattern != "foo" {
			t.Errorf("ar.ElemType.Validation.Pattern = %v; want foo", ar.ElemType.Validation.Pattern)
		}
	})

	t.Run("defined with the type name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("name", func() {
			apidsl.Attribute("id")
		})
		ar := apidsl.Type("names", func() {
			apidsl.Attribute("ut", apidsl.ArrayOf("name"))
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if ar.TypeName != "names" {
			t.Errorf("ar.TypeName = %v; want names", ar.TypeName)
		}
		o := ar.Type.ToObject()
		ut := o["ut"]
		if _, ok := ut.Type.(*design.Array); !ok {
			t.Errorf("ut.Type = %T; want *design.Array", ut.Type)
		}
		et := ut.Type.ToArray().ElemType
		if v := et.Type.(*design.UserTypeDefinition).TypeName; v != "name" {
			t.Errorf("et.Type.(*design.UserTypeDefinition).TypeName = %v; want name", v)
		}
	})

	t.Run("defined with a media type name", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/vnd.test", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("ut", apidsl.ArrayOf("application/vnd.test"))
			})
			apidsl.View("default", func() {
				apidsl.Attribute("ut")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if mt.TypeName != "Test" {
			t.Errorf("mt.TypeName = %v; want Test", mt.TypeName)
		}
		o := mt.Type.ToObject()
		ut := o["ut"]
		if _, ok := ut.Type.(*design.Array); !ok {
			t.Errorf("ut.Type = %T; want *design.Array", ut.Type)
		}
		et := ut.Type.ToArray().ElemType
		if v := et.Type.(*design.MediaTypeDefinition).TypeName; v != "Test" {
			t.Errorf("et.Type.(*design.MediaTypeDefinition).TypeName = %v; want Test", v)
		}
	})
}

func TestHashOf(t *testing.T) {
	t.Run("used on a global variable", func(t *testing.T) {
		dslengine.Reset()
		kt := apidsl.Type("key", func() {
			apidsl.Attribute("id")
		})
		vt := apidsl.Type("val", func() {
			apidsl.Attribute("id")
		})
		ha := apidsl.HashOf(kt, vt)
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if ha.Kind() != design.HashKind {
			t.Errorf("ha.Kind() = %v; want %v", ha.Kind(), design.HashKind)
		}
		if ha.KeyType.Type != kt {
			t.Errorf("ha.KeyType.Type = %v; want %v", ha.KeyType.Type, kt)
		}
		if ha.ElemType.Type != vt {
			t.Errorf("ha.ElemType.Type = %v; want %v", ha.ElemType.Type, vt)
		}
	})

	t.Run("with a DSL", func(t *testing.T) {
		dslengine.Reset()
		ha := apidsl.HashOf(design.String, design.String, func() {
			apidsl.Pattern("foo")
		}, func() {
			apidsl.Pattern("bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if ha.Kind() != design.HashKind {
			t.Errorf("ha.Kind() = %v; want %v", ha.Kind(), design.HashKind)
		}
		if ha.KeyType.Type != design.String {
			t.Errorf("ha.KeyType.Type = %v; want %v", ha.KeyType.Type, design.String)
		}
		if ha.KeyType.Validation.Pattern != "foo" {
			t.Errorf("ha.KeyType.Validation.Pattern = %v; want foo", ha.KeyType.Validation.Pattern)
		}
		if ha.ElemType.Type != design.String {
			t.Errorf("ha.ElemType.Type = %v; want %v", ha.ElemType.Type, design.String)
		}
		if ha.ElemType.Validation.Pattern != "bar" {
			t.Errorf("ha.ElemType.Validation.Pattern = %v; want bar", ha.ElemType.Validation.Pattern)
		}
	})
}
