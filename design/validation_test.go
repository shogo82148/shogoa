package design_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestValidation(t *testing.T) {
	t.Run("with a type attribute", func(t *testing.T) {

		t.Run("with a valid enum validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Enum("red", "blue")
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			want := []any{"red", "blue"}
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if diff := cmp.Diff(want, att.Validation.Values); diff != "" {
				t.Errorf("att.Validation.Values mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("with an incompatible enum validation type", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.Enum(1, "blue")
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a default value that doesn't exist in enum", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.Enum(1, 2, 3)
					apidsl.Default(4)
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a valid format validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Format("email")
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			want := "email"
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if diff := cmp.Diff(want, att.Validation.Format); diff != "" {
				t.Errorf("att.Validation.Format mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("with an invalid format validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Format("invalid")
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a valid pattern validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Pattern("^foo$")
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if att.Validation.Pattern != "^foo$" {
				t.Errorf("att.Validation.Pattern = %q; want %q", att.Validation.Pattern, "^foo$")
			}
		})

		t.Run("with an invalid pattern validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Pattern("[invalid")
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with an invalid format validation type", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.Format("email")
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a valid min value validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.Minimum(2)
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if *att.Validation.Minimum != 2 {
				t.Errorf("att.Validation.Minimum = %f; want 2", *att.Validation.Minimum)
			}
		})

		t.Run("with an invalid min value validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Minimum(2)
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a valid max value validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.Maximum(2)
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if *att.Validation.Maximum != 2 {
				t.Errorf("att.Validation.Maximum = %f; want 2", *att.Validation.Maximum)
			}
		})

		t.Run("with an invalid max value validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.Maximum(2)
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a valid min length validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", apidsl.ArrayOf(design.Integer), func() {
					apidsl.MinLength(2)
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if *att.Validation.MinLength != 2 {
				t.Errorf("att.Validation.MinLength = %d; want 2", *att.Validation.MinLength)
			}
		})

		t.Run("with an invalid min length validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.MinLength(2)
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a valid max length validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String, func() {
					apidsl.MaxLength(2)
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			o := design.Design.Types["bar"].Type.(design.Object)
			att := o["attName"]
			if *att.Validation.MaxLength != 2 {
				t.Errorf("att.Validation.MaxLength = %d; want 2", *att.Validation.MaxLength)
			}
		})

		t.Run("with an invalid max length validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.Integer, func() {
					apidsl.MaxLength(2)
				})
			})
			err := dslengine.Run()
			if err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("with a required field validation", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Type("bar", func() {
				apidsl.Attribute("attName", design.String)
				apidsl.Required("attName")
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
			validation := design.Design.Types["bar"].Validation
			if diff := cmp.Diff([]string{"attName"}, validation.Required); diff != "" {
				t.Errorf("validation.Required mismatch (-want +got):\n%s", diff)
			}
		})

	})

	t.Run("actions with different http methods", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("one", func() {
			apidsl.Action("first", func() {
				apidsl.Routing(apidsl.GET("/:first"))
			})
			apidsl.Action("second", func() {
				apidsl.Routing(apidsl.DELETE("/:second"))
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with an action", func(t *testing.T) {

		t.Run("which has a file type param", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Resource("foo", func() {
				apidsl.Action("bar", func() {
					apidsl.Routing(apidsl.GET("/buz"))
					apidsl.Params(func() {
						apidsl.Param("file", design.File) // action params cannot be a file
					})
				})
			})
			if err := dslengine.Run(); err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("which has a file array type param", func(t *testing.T) {
			dslengine.Reset()
			apidsl.Resource("foo", func() {
				apidsl.Action("bar", func() {
					apidsl.Routing(apidsl.GET("/buz"))
					apidsl.Params(func() {
						apidsl.Param("file_array", apidsl.ArrayOf(design.File)) // action params cannot be a file array
					})
				})
			})
			if err := dslengine.Run(); err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("which has a payload contains a file", func(t *testing.T) {
			dslengine.Reset()
			var payload = apidsl.Type("qux", func() {
				apidsl.Attribute("file", design.File)
				apidsl.Required("file")
			})
			apidsl.Resource("foo", func() {
				apidsl.Action("bar", func() {
					apidsl.Routing(apidsl.GET("/buz"))
					apidsl.Payload(payload) // action payloads cannot contain a file
				})
			})
			if err := dslengine.Run(); err == nil {
				t.Fatal("expected an error")
			}
		})

		t.Run("which has a payload contains a file and multipart form", func(t *testing.T) {
			dslengine.Reset()
			var payload = apidsl.Type("qux", func() {
				apidsl.Attribute("file", design.File)
				apidsl.Required("file")
			})
			apidsl.Resource("foo", func() {
				apidsl.Action("bar", func() {
					apidsl.Routing(apidsl.GET("/buz"))
					apidsl.Payload(payload)
					apidsl.MultipartForm()
				})
			})
			if err := dslengine.Run(); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("which has a response contains a file", func(t *testing.T) {
			dslengine.Reset()
			var response = apidsl.MediaType("application/vnd.shogoa.example", func() {
				apidsl.TypeName("quux")
				apidsl.Attributes(func() {
					apidsl.Attribute("file", design.File)
					apidsl.Required("file")
				})
				apidsl.View("default", func() {
					apidsl.Attribute("file")
				})
			})
			apidsl.Resource("foo", func() {
				apidsl.Action("bar", func() {
					apidsl.Routing(apidsl.GET("/buz"))
					apidsl.Response(design.OK, response) // action responses cannot contain a file
				})
			})
			if err := dslengine.Run(); err == nil {
				t.Fatal("expected an error")
			}
		})

	})
}
