package apidsl_test

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestResponse(t *testing.T) {
	t.Run("with no dsl and no name", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("")
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses[""]
		if err := res.Validate(); err == nil {
			t.Error("Validate() = nil; want an error")
		}
	})

	t.Run("with no dsl", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo")
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err == nil {
			t.Error("Validate() = nil; want an error")
		}
	})

	t.Run("with a status", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo", func() {
					apidsl.Status(201)
				})
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if res.Status != 201 {
			t.Errorf("Status = %d; want 201", res.Status)
		}
		if res.Parent == nil {
			t.Error("Parent = nil; want not nil")
		}
	})

	t.Run("with a type override", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				dt := apidsl.HashOf(design.String, design.Any)
				apidsl.Response("foo", dt, func() {
					apidsl.Status(201)
				})
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if res.Status != 201 {
			t.Errorf("Status = %d; want 201", res.Status)
		}
		if res.Type == nil {
			t.Error("Type = nil; want not nil")
		}
	})

	t.Run("with a status and description", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo", func() {
					apidsl.Status(201)
					apidsl.Description("desc")
				})
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if res.Status != 201 {
			t.Errorf("Status = %d; want 201", res.Status)
		}
		if res.Description != "desc" {
			t.Errorf("Description = %q; want %q", res.Description, "desc")
		}
	})

	t.Run("with a status and name override", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo", func() {
					apidsl.Status(201)
				})
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if res.Status != 201 {
			t.Errorf("Status = %d; want 201", res.Status)
		}
	})

	t.Run("with a status and media type", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo", func() {
					apidsl.Status(201)
					apidsl.Media("mt")
				})
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if res.Status != 201 {
			t.Errorf("Status = %d; want 201", res.Status)
		}
		if res.MediaType != "mt" {
			t.Errorf("MediaType = %q; want %q", res.MediaType, "mt")
		}
	})

	t.Run("with a status and headers", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo", func() {
					apidsl.Status(201)
					apidsl.Headers(func() {
						apidsl.Header("Location")
					})
				})
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if err := res.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if res.Status != 201 {
			t.Errorf("Status = %d; want 201", res.Status)
		}
		if res.Headers == nil {
			t.Error("Headers = nil; want not nil")
		}
	})

	t.Run("not from the shogoa default definitions", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("foo")
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["foo"]
		if res.Standard {
			t.Error("Standard = true; want false")
		}
	})

	t.Run("from the shogoa default definitions", func(t *testing.T) {
		dslengine.Reset()
		apidsl.Resource("res", func() {
			apidsl.Action("action", func() {
				apidsl.Response("Created")
			})
		})
		_ = dslengine.Run()

		res := design.Design.Resources["res"].Actions["action"].Responses["Created"]
		if !res.Standard {
			t.Error("Standard = false; want true")
		}
	})
}
