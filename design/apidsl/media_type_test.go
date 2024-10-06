package apidsl_test

import (
	"regexp"
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestMediaType(t *testing.T) {
	t.Run("with no dsl and no identifier", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("", nil)
		if err := dslengine.Run(); err == nil {
			t.Error("Run() = nil; want an error")
		}

		if err := mt.Validate(); err == nil {
			t.Error("Validate() = nil; want an error")
		}
	})

	t.Run("with no dsl", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/foo", nil)
		if err := dslengine.Run(); err == nil {
			t.Error("Run() = nil; want an error")
		}

		if err := mt.Validate(); err == nil {
			t.Error("Validate() = nil; want an error")
		}
	})

	t.Run("with attributes", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/foo", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := mt.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if mt.AttributeDefinition == nil {
			t.Error("AttributeDefinition = nil; want not nil")
		}
		o := mt.Type.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if _, ok := o["name"]; !ok {
			t.Error("name is not found")
		}
	})

	t.Run("with a content type", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/foo", func() {
			apidsl.ContentType("application/json")
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := mt.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if mt.ContentType != "application/json" {
			t.Errorf("ContentType = %s; want application/json", mt.ContentType)
		}
	})

	t.Run("with a description", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/foo", func() {
			apidsl.Description("desc")
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := mt.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if mt.Description != "desc" {
			t.Errorf("Description = %q; want %q", mt.Description, "desc")
		}
	})

	t.Run("with links", func(t *testing.T) {
		dslengine.Reset()
		mt1 := design.NewMediaTypeDefinition("application/mt1", "application/mt1", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("foo")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("foo")
			})
			apidsl.View("link", func() {
				apidsl.Attribute("foo")
			})
		})
		mt2 := design.NewMediaTypeDefinition("application/mt2", "application/mt2", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("foo")
			})
			apidsl.View("l2v", func() {
				apidsl.Attribute("foo")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("foo")
			})
		})
		design.Design.MediaTypes = map[string]*design.MediaTypeDefinition{
			"application/mt1": mt1,
			"application/mt2": mt2,
		}
		mt := apidsl.MediaType("foo", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("l1", mt1)
				apidsl.Attribute("l2", mt2)
			})
			apidsl.Links(func() {
				apidsl.Link("l1")
				apidsl.Link("l2", "l2v")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("l1")
				apidsl.Attribute("l2")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := mt.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if len(mt.Links) != 2 {
			t.Errorf("len(mt.Links) = %d; want 2", len(mt.Links))
		}
		link1 := mt.Links["l1"]
		if link1.Name != "l1" {
			t.Errorf("link1.Name = %q; want %q", link1.Name, "l1")
		}
		if link1.View != "link" {
			t.Errorf("link1.View = %q; want %q", link1.View, "link")
		}
		if link1.Parent != mt {
			t.Error("link1.Parent is not set")
		}
		link2 := mt.Links["l2"]
		if link2.Name != "l2" {
			t.Errorf("link2.Name = %q; want %q", link2.Name, "l2")
		}
		if link2.View != "l2v" {
			t.Errorf("link2.View = %q; want %q", link2.View, "l2v")
		}
		if link2.Parent != mt {
			t.Error("link2.Parent is not set")
		}
	})

	t.Run("with views", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/foo", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
			apidsl.View("view", func() {
				apidsl.Attribute("name")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := mt.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if len(mt.Views) != 2 {
			t.Errorf("len(mt.Views) = %d; want 2", len(mt.Views))
		}
		v := mt.Views["view"]
		if v.Name != "view" {
			t.Errorf("v.Name = %q; want %q", v.Name, "view")
		}
		if v.Parent != mt {
			t.Error("v.Parent is not set")
		}
		o := v.AttributeDefinition.Type.(design.Object)
		if len(o) != 1 {
			t.Errorf("len(o) = %d; want 1", len(o))
		}
		if o["name"].Type != design.String {
			t.Errorf("o['name'].Type = %s; want string", o["name"].Type)
		}
	})
}

func TestTestMediaType_duplicated(t *testing.T) {
	dslengine.Reset()
	apidsl.MediaType("application/foo", func() {
		apidsl.Attributes(func() {
			apidsl.Attribute("name")
		})
		apidsl.View("default", func() {
			apidsl.Attribute("name")
		})
	})
	duplicate := apidsl.MediaType("application/foo", func() {
		apidsl.Attributes(func() {
			apidsl.Attribute("name")
		})
		apidsl.View("default", func() {
			apidsl.Attribute("name")
		})
	})
	apidsl.Resource("foo", func() {
		apidsl.Action("show", func() {
			apidsl.Routing(apidsl.GET(""))
			apidsl.Response(design.OK, func() {
				apidsl.Media(duplicate)
			})
		})
	})

	_ = dslengine.Run()
}

func TestCollectionOf(t *testing.T) {
	t.Run("used on a global variable", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/vnd.example", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
		collection := apidsl.CollectionOf(mt)
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := collection.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if collection.Identifier != "application/vnd.example; type=collection" {
			t.Errorf("Identifier = %q; want %q", collection.Identifier, "application/vnd.example; type=collection")
		}
		if _, ok := design.Design.MediaTypes[collection.Identifier]; !ok {
			t.Errorf("design.Design.MediaTypes[%q] is not found", collection.Identifier)
		}
	})

	t.Run("defined with a collection identifier", func(t *testing.T) {
		dslengine.Reset()
		mt := apidsl.MediaType("application/vnd.example", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
		collection := apidsl.CollectionOf(mt, "application/vnd.examples")
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if err := collection.Validate(); err != nil {
			t.Errorf("Validate() = %v; want nil", err)
		}
		if collection.Identifier != "application/vnd.examples" {
			t.Errorf("Identifier = %q; want %q", collection.Identifier, "application/vnd.examples")
		}
		if _, ok := design.Design.MediaTypes[collection.Identifier]; !ok {
			t.Errorf("design.Design.MediaTypes[%q] is not found", collection.Identifier)
		}
	})

	t.Run("defined with the media type identifier", func(t *testing.T) {
		dslengine.Reset()
		apidsl.MediaType("application/vnd.example+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("name")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
		collection := apidsl.MediaType("application/vnd.parent+json", func() {
			apidsl.Attribute("mt", apidsl.CollectionOf("application/vnd.example"))
			apidsl.View("default", func() {
				apidsl.Attribute("mt")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		if collection.Identifier != "application/vnd.parent+json" {
			t.Errorf("Identifier = %q; want %q", collection.Identifier, "application/vnd.parent+json")
		}
		if collection.TypeName != "Parent" {
			t.Errorf("TypeName = %q; want %q", collection.TypeName, "Parent")
		}
		mt := collection.Type.(design.Object)["mt"]
		if mt.Type.Name() != "array" {
			t.Errorf("mt.Type.Name() = %q; want %q", mt.Type.Name(), "array")
		}
		et := mt.Type.ToArray().ElemType
		if et.Type.(*design.MediaTypeDefinition).Identifier != "application/vnd.example+json" {
			t.Errorf("et.Type.(*design.MediaTypeDefinition).Identifier = %q; want %q", et.Type.(*design.MediaTypeDefinition).Identifier, "application/vnd.example+json")
		}
	})
}

func TestExample(t *testing.T) {
	t.Run("produces a media type with examples", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := apidsl.MediaType("application/vnd.example+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 description", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 description", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", design.Integer, "test3 description", func() {
					apidsl.Minimum(1)
				})
				apidsl.Attribute("test4", design.String, func() {
					apidsl.Format("email")
					apidsl.Pattern("@")
				})
				apidsl.Attribute("test5", design.Any)

				apidsl.Attribute("test-failure1", design.Integer, func() {
					apidsl.Minimum(0)
					apidsl.Maximum(0)
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if o["test1"].Example != "test1" {
			t.Errorf("o['test1'].Example = %q; want %q", o["test1"].Example, "test1")
		}
		if o["test2"].Example != "-" {
			t.Errorf("o['test2'].Example = %q; want %q", o["test2"].Example, "-")
		}
		if v := o["test3"].Example.(int); v < 1 {
			t.Errorf("o['test3'].Example = %v; want >= 1", v)
		}
		if ok, err := regexp.MatchString(`\w+@`, o["test4"].Example.(string)); !ok || err != nil {
			t.Errorf("o['test4'].Example = %q; want to match %q", o["test4"].Example, `\w+@`)
		}
		if o["test5"].Example == nil {
			t.Errorf("o['test5'].Example = nil; want not nil")
		}
		if o["test-failure1"].Example != 0 {
			t.Errorf("o['test-failure1'].Example = %v; want 0", o["test-failure1"].Example)
		}
	})

	t.Run("produces a media type with HashOf examples", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		ut := apidsl.Type("example", func() {
			apidsl.Attribute("test1", design.Integer)
			apidsl.Attribute("test2", design.Any)
		})
		mt := apidsl.MediaType("application/vnd.example+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", apidsl.HashOf(design.String, design.Integer))
				apidsl.Attribute("test2", apidsl.HashOf(design.Any, design.String))
				apidsl.Attribute("test3", apidsl.HashOf(design.String, design.Any))
				apidsl.Attribute("test4", apidsl.HashOf(design.Any, design.Any))

				apidsl.Attribute("test-with-user-type-1", apidsl.HashOf(design.String, ut))
				apidsl.Attribute("test-with-user-type-2", apidsl.HashOf(design.Any, ut))

				apidsl.Attribute("test-with-array-1", apidsl.HashOf(design.String, apidsl.ArrayOf(design.Integer)))
				apidsl.Attribute("test-with-array-2", apidsl.HashOf(design.String, apidsl.ArrayOf(design.Any)))
				apidsl.Attribute("test-with-array-3", apidsl.HashOf(design.String, apidsl.ArrayOf(ut)))
				apidsl.Attribute("test-with-array-4", apidsl.HashOf(design.Any, apidsl.ArrayOf(design.String)))
				apidsl.Attribute("test-with-array-5", apidsl.HashOf(design.Any, apidsl.ArrayOf(design.Any)))
				apidsl.Attribute("test-with-array-6", apidsl.HashOf(design.Any, apidsl.ArrayOf(ut)))

				apidsl.Attribute("test-with-example-1", apidsl.HashOf(design.String, design.Boolean), func() {
					apidsl.Example(map[string]bool{})
				})
				apidsl.Attribute("test-with-example-2", apidsl.HashOf(design.Any, design.Boolean), func() {
					apidsl.Example(map[string]int{})
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if _, ok := o["test1"].Example.(map[string]int); !ok {
			t.Errorf("o['test1'].Example = %T; want map[string]int", o["test1"].Type)
		}
		if _, ok := o["test2"].Example.(map[any]string); !ok {
			t.Errorf("o['test2'].Example = %T; want map[any]string", o["test2"].Type)
		}
		if _, ok := o["test3"].Example.(map[string]any); !ok {
			t.Errorf("o['test3'].Example = %T; want map[string]any", o["test3"].Type)
		}
		if _, ok := o["test4"].Example.(map[any]any); !ok {
			t.Errorf("o['test4'].Example = %T; want map[any]any", o["test4"].Type)
		}

		for _, attr := range o["test-with-user-type-1"].Example.(map[string]map[string]any) {
			if _, ok := attr["test1"].(int); !ok {
				t.Errorf("attr['test1'] = %T; want int", attr["test1"])
			}
			if _, ok := attr["test2"]; !ok {
				t.Errorf("attr['test2'] is not found")
			}
		}
		for _, attr := range o["test-with-user-type-2"].Example.(map[any]map[string]any) {
			if _, ok := attr["test1"].(int); !ok {
				t.Errorf("attr['test1'] = %T; want int", attr["test1"])
			}
			if _, ok := attr["test2"]; !ok {
				t.Errorf("attr['test2'] is not found")
			}
		}

		if _, ok := o["test-with-array-1"].Example.(map[string][]int); !ok {
			t.Errorf("o['test-with-array-1'].Example = %T; want map[string][]int", o["test-with-array-1"].Example)
		}
		if _, ok := o["test-with-array-2"].Example.(map[string][]any); !ok {
			t.Errorf("o['test-with-array-2'].Example = %T; want map[string][]string", o["test-with-array-2"].Example)
		}
		if _, ok := o["test-with-array-3"].Example.(map[string][]map[string]any); !ok {
			t.Errorf("o['test-with-array-3'].Example = %T; want map[string][]map[string]any", o["test-with-array-3"].Example)
		}
		if _, ok := o["test-with-array-4"].Example.(map[any][]string); !ok {
			t.Errorf("o['test-with-array-4'].Example = %T; want map[any][]string", o["test-with-array-4"].Example)
		}
		if _, ok := o["test-with-array-5"].Example.(map[any][]any); !ok {
			t.Errorf("o['test-with-array-5'].Example = %T; want map[any][]any", o["test-with-array-5"].Example)
		}
		if _, ok := o["test-with-array-6"].Example.(map[any][]map[string]any); !ok {
			t.Errorf("o['test-with-array-6'].Example = %T; want map[any][]map[string]any", o["test-with-array-6"].Example)
		}

		if _, ok := o["test-with-example-1"].Example.(map[string]bool); !ok {
			t.Errorf("o['test-with-example-1'].Example = %T; want map[string]bool", o["test-with-example-1"].Example)
		}
		if _, ok := o["test-with-example-2"].Example.(map[string]int); !ok {
			t.Errorf("o['test-with-example-2'].Example = %T; want map[any]bool", o["test-with-example-2"].Example)
		}
	})

	t.Run("produces a media type with examples in cyclical dependencies", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := apidsl.MediaType("vnd.application/foo", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("foo", "vnd.application/bar")
				apidsl.Attribute("others", design.Integer, func() {
					apidsl.Minimum(3)
					apidsl.Maximum(3)
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("foo")
				apidsl.Attribute("others")
			})
		})
		mt2 := apidsl.MediaType("vnd.application/bar", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("bar", mt)
				apidsl.Attribute("others", design.Integer, func() {
					apidsl.Minimum(1)
					apidsl.Maximum(2)
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("bar")
				apidsl.Attribute("others")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		child := o["foo"].Example.(map[string]any)
		if _, ok := child["bar"]; !ok {
			t.Errorf("child['bar'] is not found")
		}
		if v := child["others"].(int); v < 1 {
			t.Errorf("child['others'] = %v; want >= 1", v)
		}
		if v := child["others"].(int); v > 2 {
			t.Errorf("child['others'] = %v; want <= 2", v)
		}
		attr := mt.Type.ToObject()["others"]
		if v := attr.Example.(int); v != 3 {
			t.Errorf("attr.Example = %v; want 3", v)
		}

		o2 := mt2.Type.ToObject()
		child2 := o2["bar"].Example.(map[string]any)
		if _, ok := child2["foo"]; !ok {
			t.Errorf("child2['foo'] is not found")
		}
		if v := child2["others"].(int); v != 3 {
			t.Errorf("child2['others'] = %v; want 3", v)
		}
		attr2 := mt2.Type.ToObject()["others"]
		if v := attr2.Example.(int); v < 1 {
			t.Errorf("attr2.Example = %v; want >= 1", v)
		}
		if v := attr2.Example.(int); v > 2 {
			t.Errorf("attr2.Example = %v; want <= 2", v)
		}
	})

	t.Run("produces media type examples from the linked media type", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := apidsl.MediaType("application/vnd.example+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 desc", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", design.Integer, "test3 desc", func() {
					apidsl.Minimum(1)
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
				apidsl.Attribute("test2")
				apidsl.Attribute("test3")
			})
		})
		pmt := apidsl.MediaType("application/vnd.example.parent+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 desc", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", design.Integer, "test3 desc", func() {
					apidsl.Minimum(1)
				})
				apidsl.Attribute("test4", mt, "test4 desc")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
				apidsl.Attribute("test2")
				apidsl.Attribute("test3")
				apidsl.Attribute("test4")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if o["test1"].Example != "test1" {
			t.Errorf("o['test1'].Example = %q; want %q", o["test1"].Example, "test1")
		}
		if o["test2"].Example != "-" {
			t.Errorf("o['test2'].Example = %q; want %q", o["test2"].Example, "-")
		}
		if v := o["test3"].Example.(int); v < 1 {
			t.Errorf("o['test3'].Example = %v; want >= 1", v)
		}

		o2 := pmt.Type.ToObject()
		if o2["test1"].Example != "test1" {
			t.Errorf("o2['test1'].Example = %q; want %q", o2["test1"].Example, "test1")
		}
		if o2["test2"].Example != "-" {
			t.Errorf("o2['test2'].Example = %q; want %q", o2["test2"].Example, "-")
		}
		if v := o2["test3"].Example.(int); v < 1 {
			t.Errorf("o2['test3'].Example = %v; want >= 1", v)
		}
		child := o2["test4"].Example.(map[string]any)
		if child["test1"] != "test1" {
			t.Errorf("child['test1'] = %q; want %q", child["test1"], "test1")
		}
		if child["test2"] != "-" {
			t.Errorf("child['test2'] = %q; want %q", child["test2"], "-")
		}
		if v := child["test3"].(int); v < 1 {
			t.Errorf("child['test3'] = %v; want >= 1", v)
		}
	})

	t.Run("produces media type examples from the linked media type collection with custom examples", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := apidsl.MediaType("application/vnd.example+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 desc", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", design.Integer, "test3 desc", func() {
					apidsl.Minimum(1)
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
				apidsl.Attribute("test2")
				apidsl.Attribute("test3")
			})
		})

		pmt := apidsl.MediaType("application/vnd.example.parent+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 desc", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", design.String, "test3 desc", func() {
					apidsl.Pattern("^1$")
				})
				apidsl.Attribute("test4", apidsl.CollectionOf(mt), "test4 desc")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
				apidsl.Attribute("test2")
				apidsl.Attribute("test3")
				apidsl.Attribute("test4")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if o["test1"].Example != "test1" {
			t.Errorf("o['test1'].Example = %q; want %q", o["test1"].Example, "test1")
		}
		if o["test2"].Example != "-" {
			t.Errorf("o['test2'].Example = %q; want %q", o["test2"].Example, "-")
		}
		if v := o["test3"].Example.(int); v < 1 {
			t.Errorf("o['test3'].Example = %v; want >= 1", v)
		}

		o2 := pmt.Type.ToObject()
		if o2["test1"].Example != "test1" {
			t.Errorf("o2['test1'].Example = %q; want %q", o2["test1"].Example, "test1")
		}
		if o2["test2"].Example != "-" {
			t.Errorf("o2['test2'].Example = %q; want %q", o2["test2"].Example, "-")
		}
		if o2["test3"].Example != "1" {
			t.Errorf("o2['test3'].Example = %q; want %q", o2["test3"].Example, "1")
		}
		child := o2["test4"].Example.([]map[string]any)
		if len(child) != 1 {
			t.Errorf("len(child) = %d; want 1", len(child))
		}
		if child[0]["test1"] != "test1" {
			t.Errorf("child[0]['test1'] = %q; want %q", child[0]["test1"], "test1")
		}
		if child[0]["test2"] != "-" {
			t.Errorf("child[0]['test2'] = %q; want %q", child[0]["test2"], "-")
		}
		if v := child[0]["test3"].(int); v < 1 {
			t.Errorf("child[0]['test3'] = %v; want >= 1", v)
		}
	})

	t.Run("produces media type examples from the linked media type without custom examples", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := apidsl.MediaType("application/vnd.example.child+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
			})
		})

		pmt := apidsl.MediaType("application/vnd.example.parent+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 desc", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", mt, "test3 desc")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
				apidsl.Attribute("test2")
				apidsl.Attribute("test3")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if o["test1"].Example == "" {
			t.Errorf("o['test1'].Example = %q; want not empty", o["test1"].Example)
		}

		o2 := pmt.Type.ToObject()
		if o2["test1"].Example != "test1" {
			t.Errorf("o2['test1'].Example = %q; want %q", o2["test1"].Example, "test1")
		}
		if o2["test2"].Example != "-" {
			t.Errorf("o2['test2'].Example = %q; want %q", o2["test2"].Example, "-")
		}
		child := o2["test3"].Example.(map[string]any)
		if child["test1"] == "" {
			t.Errorf("child['test1'].Example = %q; want not empty", child["test1"])
		}
	})

	t.Run("produces media type examples from the linked media type collection without custom examples", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := apidsl.MediaType("application/vnd.example.child+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
			})
		})
		pmt := apidsl.MediaType("application/vnd.example.parent+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", design.String, "test1 desc", func() {
					apidsl.Example("test1")
				})
				apidsl.Attribute("test2", design.String, "test2 desc", func() {
					apidsl.NoExample()
				})
				apidsl.Attribute("test3", apidsl.CollectionOf(mt), "test3 desc")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if o["test1"].Example == "" {
			t.Errorf("o['test1'].Example = %q; want not empty", o["test1"].Example)
		}

		o2 := pmt.Type.ToObject()
		if o2["test1"].Example != "test1" {
			t.Errorf("o2['test1'].Example = %q; want %q", o2["test1"].Example, "test1")
		}
		if o2["test2"].Example != "-" {
			t.Errorf("o2['test2'].Example = %q; want %q", o2["test2"].Example, "-")
		}
		child := o2["test3"].Example.([]map[string]any)
		if len(child) < 1 {
			t.Errorf("len(child) = %d; want >= 1", len(child))
		}
	})

	t.Run("produces a media type with appropriate MinLength and MaxLength examples", func(t *testing.T) {
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		ut := apidsl.Type("example", func() {
			apidsl.Attribute("test1", design.Integer, func() {
				apidsl.Minimum(-200)
				apidsl.Maximum(-100)
			})
		})
		mt := apidsl.MediaType("application/vnd.example+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("test1", apidsl.ArrayOf(design.Any), func() {
					apidsl.MinLength(0)
					apidsl.MaxLength(10)
				})
				apidsl.Attribute("test2", apidsl.ArrayOf(design.Any), func() {
					apidsl.MinLength(1000)
					apidsl.MaxLength(2000)
				})
				apidsl.Attribute("test3", apidsl.ArrayOf(design.Any), func() {
					apidsl.MinLength(1000)
					apidsl.MaxLength(1000)
				})

				apidsl.Attribute("test-failure1", apidsl.ArrayOf(ut), func() {
					apidsl.MinLength(0)
					apidsl.MaxLength(0)
				})
			})
			apidsl.View("default", func() {
				apidsl.Attribute("test1")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		o := mt.Type.ToObject()
		if len(o["test1"].Example.([]any)) > 10 {
			t.Errorf("len(o['test1'].Example) = %d; want <= 10", len(o["test1"].Example.([]any)))
		}
		if len(o["test2"].Example.([]any)) != 10 {
			t.Errorf("len(o['test2'].Example) = %d; want 10", len(o["test2"].Example.([]any)))
		}
		if len(o["test3"].Example.([]any)) != 10 {
			t.Errorf("len(o['test3'].Example) = %d; want 10", len(o["test3"].Example.([]any)))
		}
		if o["test-failure1"].Example != nil {
			t.Errorf("o['test-failure1'].Example = %v; want nil", o["test-failure1"].Example)
		}
	})
}
