package apidsl_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
	_ = it([]dslengine.Definition{t})
}

// Reset is a no-op
func (t *TestCD) Reset() {}

func TestContainerDefinition(t *testing.T) {
	att := &design.AttributeDefinition{Type: design.Object{}}
	testCD := &TestCD{AttributeDefinition: att}
	dslengine.Register(testCD)
	if err := dslengine.Run(); err != nil {
		t.Fatalf("Run() = %v; want nil", err)
	}
	if got, want := testCD.Attribute(), att; got != want {
		t.Errorf("Attribute() = %v; want %v", got, want)
	}
}

func TestAttribute(t *testing.T) {
	t.Run("with only a name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.String {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.String)
		}
	})

	t.Run("with a name and datatype", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.Integer)
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.Integer {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.Integer)
		}
	})

	t.Run("with a name and uuid datatype", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.UUID)
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.UUID {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.UUID)
		}
	})

	t.Run("with a name and date datatype", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.DateTime)
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.DateTime {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.DateTime)
		}
	})

	t.Run("with a name, datatype and description", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.Integer, "bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.Integer {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.Integer)
		}
		if o["foo"].Description != "bar" {
			t.Errorf("Description = %q; want %q", o["foo"].Description, "bar")
		}
	})

	t.Run("with a name and a DSL defining a 'readOnly' attribute", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", func() { apidsl.ReadOnly() })
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if !o["foo"].IsReadOnly() {
			t.Error("ReadOnly = false; want true")
		}
	})

	t.Run("with a name and a DSL defining an enum validation", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", func() { apidsl.Enum("one", "two") })
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.String {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.String)
		}

		want := []any{"one", "two"}
		got := o["foo"].Validation.Values
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validation.Values mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("with a name, type datetime and a DSL defining a default value", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.DateTime, func() { apidsl.Default("1978-06-30T10:00:00+09:00") })
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.DateTime {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.DateTime)
		}
		if o["foo"].Validation != nil {
			t.Error("Validation != nil; want nil")
		}
		if got, want := o["foo"].DefaultValue, any("1978-06-30T10:00:00+09:00"); got != want {
			t.Errorf("DefaultValue = %v; want %v", got, want)
		}
	})

	t.Run("with a name, type integer and a DSL defining an enum validation", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.Integer, func() { apidsl.Enum(1, 2) })
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.Integer {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.Integer)
		}
		if o["foo"].Validation == nil {
			t.Error("Validation = nil; want not nil")
		}

		got := o["foo"].Validation.Values
		want := []any{1, 2}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validation.Values mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("with a name, type integer, a description and a DSL defining an enum validation", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.String, "bar", func() { apidsl.Enum("one", "two") })
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.String {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.String)
		}
		if o["foo"].Validation == nil {
			t.Error("Validation = nil; want not nil")
		}

		got := o["foo"].Validation.Values
		want := []any{"one", "two"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validation.Values mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("with a name and type uuid", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.UUID)
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.UUID {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.UUID)
		}
	})

	t.Run("with a name and type date", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.DateTime)
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		if typ == nil {
			t.Error("Type = nil; want not nil")
		}
		o := typ.(design.Object)
		if o["foo"].Type != design.DateTime {
			t.Errorf("Type = %v; want %v", o["foo"].Type, design.DateTime)
		}
	})

	t.Run("with a name and a type defined by name", func(t *testing.T) {
		dslengine.Reset()
		foo := apidsl.Type("bar", func() {
			apidsl.Attribute("baz")
		})
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", "bar")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatalf("Run() = %v; want nil", err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		o := typ.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if o["foo"].Type != foo {
			t.Errorf("Type = %v; want %v", o["foo"].Type, foo)
		}
	})

	t.Run("with child attributes on an attribute that is not an object", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.Integer, func() {
				apidsl.Attribute("bar")
			})
		})
		if err := dslengine.Run(); err == nil {
			t.Error("Run() = nil; want an error")
		}
	})

	t.Run("with child attributes on an attribute that does not have a type", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", func() {
				apidsl.Attribute("bar")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		o := typ.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if _, ok := o["foo"].Type.(design.Object); !ok {
			t.Errorf("o[\"foo\"].Type is not an object")
		}
	})

	t.Run("with child attribute on an attribute of type Object", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Type("type", func() {
			apidsl.Attribute("foo", design.Object{}, func() {
				apidsl.Attribute("bar")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		parent := design.Design.Types["type"].AttributeDefinition
		typ := parent.Type
		o := typ.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if _, ok := o["foo"].Type.(design.Object); !ok {
			t.Errorf("o[\"foo\"].Type is not an object")
		}
		co := o["foo"].Type.(design.Object)
		if len(co) != 1 {
			t.Errorf("len(co) = %d; want 1", len(co))
		}
		if _, ok := co["bar"]; !ok {
			t.Errorf("co[\"bar\"] is not found")
		}
	})
}
