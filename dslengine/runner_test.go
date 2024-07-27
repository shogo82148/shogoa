package dslengine_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestDSL(t *testing.T) {
	current_line := func() int {
		_, _, line, _ := runtime.Caller(1)
		return line
	}

	t.Run("with cyclical type dependencies", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		apidsl.API("foo", func() {})
		apidsl.Type("type1Name", func() {
			apidsl.Attribute("att1Name", "type2Name")
		})
		apidsl.Type("type2Name", func() {
			apidsl.Attribute("att2Name", "type1Name")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		// verify the result
		if len(dslengine.Errors) != 0 {
			t.Errorf("unexpected errors: %v", dslengine.Errors)
		}
		if len(design.Design.Types) != 2 {
			t.Errorf("unexpected number of types: %v", design.Design.Types)
		}
		t1 := design.Design.Types["type1Name"]
		t2 := design.Design.Types["type2Name"]
		o1 := t1.Type.(design.Object)
		o2 := t2.Type.(design.Object)
		if o1["att1Name"].Type != t2 {
			t.Errorf("unexpected type want %v, got %v", t2, o1["att1Name"].Type)
		}
		if o2["att2Name"].Type != t1 {
			t.Errorf("unexpected type: want %v, got %v", t1, o2["att2Name"].Type)
		}
	})

	t.Run("with one error", func(t *testing.T) {
		// See NOTE below.
		var lineNumber = current_line() + 6

		// define the DSL
		dslengine.Reset()
		// NOTE: moving the line below requires updating
		// lineNumber above to match its number.
		dslengine.ReportError("err")

		// verify the result
		if len(dslengine.Errors) != 1 {
			t.Errorf("unexpected number of errors: %v", dslengine.Errors)
		}
		err := dslengine.Errors[0]
		if err.File != "runner_test.go" {
			t.Errorf("unexpected file: %v", err.File)
		}
		if err.Line != lineNumber {
			t.Errorf("unexpected line: %d, want %d", err.Line, lineNumber)
		}
	})

	t.Run("with multiple errors", func(t *testing.T) {
		// define the DSL
		dslengine.Reset()
		dslengine.ReportError("foo1")
		dslengine.ReportError("foo2")

		// verify the result
		if len(dslengine.Errors) != 2 {
			t.Errorf("unexpected number of errors: %v", dslengine.Errors)
		}
		msg := dslengine.Errors.Error()
		if !strings.Contains(msg, "foo1") {
			t.Errorf("unexpected message: %v", msg)
		}
		if !strings.Contains(msg, "foo2") {
			t.Errorf("unexpected message: %v", msg)
		}
	})

	t.Run("with invalid DSL", func(t *testing.T) {
		// See NOTE below.
		var lineNumber = current_line() + 7

		// define the DSL
		dslengine.Reset()
		apidsl.API("foo", func() {
			// NOTE: moving the line below requires updating
			// lineNumber above to match its number.
			apidsl.Attributes(func() {})
		})
		if err := dslengine.Run(); err == nil {
			t.Fatal("expected an error")
		}

		// verify the result
		if len(dslengine.Errors) != 1 {
			t.Errorf("unexpected number of errors: %v", dslengine.Errors)
		}
		err := dslengine.Errors[0]
		if !strings.Contains(err.Error(), "invalid use of Attributes") {
			t.Errorf("unexpected error: %v", err)
		}
		if err.File != "runner_test.go" {
			t.Errorf("unexpected file: %v", err.File)
		}
		if err.Line != lineNumber {
			t.Errorf("unexpected line: %v", err.Line)
		}
	})

	t.Run("with DSL calling a function with an invalid argument type", func(t *testing.T) {
		// See NOTE below.
		var lineNumber = current_line() + 7

		// define the DSL
		dslengine.Reset()
		apidsl.Type("bar", func() {
			// NOTE: moving the line below requires updating
			// lineNumber above to match its number.
			apidsl.Attribute("baz", 42)
		})
		if err := dslengine.Run(); err == nil {
			t.Fatal("expected an error")
		}

		// verify the result
		if len(dslengine.Errors) != 1 {
			t.Errorf("unexpected number of errors: %v", dslengine.Errors)
		}
		err := dslengine.Errors[0]
		if !strings.Contains(err.Error(), "cannot use 42 (type int) as type") {
			t.Errorf("unexpected error: %v", err)
		}
		if err.File != "runner_test.go" {
			t.Errorf("unexpected file: %v", err.File)
		}
		if err.Line != lineNumber {
			t.Errorf("unexpected line: %v", err.Line)
		}
	})

	t.Run("with DSL using an empty type", func(t *testing.T) {
		// define DSL
		apidsl.API("foo", func() {})
		apidsl.Resource("bar", func() {
			apidsl.Action("baz", func() {
				apidsl.Payload("use-empty")
			})
		})
		apidsl.Type("use-empty", func() {
			apidsl.Attribute("e", "empty")
		})
		apidsl.Type("empty", func() {})
	})
	if err := dslengine.Run(); err == nil {
		t.Fatal("expected an error")
	}
}
