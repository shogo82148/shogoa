package design_test

import (
	"mime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestDataType_IsObject(t *testing.T) {
	tests := []struct {
		name     string
		dataType design.DataType
		want     bool
	}{
		{
			name:     "primitive",
			dataType: design.String,
			want:     false,
		},
		{
			name:     "array",
			dataType: &design.Array{ElemType: &design.AttributeDefinition{Type: design.String}},
			want:     false,
		},
		{
			name: "hash",
			dataType: &design.Hash{
				KeyType:  &design.AttributeDefinition{Type: design.String},
				ElemType: &design.AttributeDefinition{Type: design.String},
			},
			want: false,
		},
		{
			name: "nil user type",
			dataType: &design.UserTypeDefinition{
				AttributeDefinition: &design.AttributeDefinition{Type: nil},
			},
		},
		{
			name:     "object",
			dataType: &design.Object{},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dataType.IsObject(); got != tt.want {
				t.Errorf("DataType.IsObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaTypeDefinition_Project(t *testing.T) {
	t.Run("with a media type with multiple views", func(t *testing.T) {
		design.Design.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)

		mt := &design.MediaTypeDefinition{
			Identifier: "vnd.application/foo",
			UserTypeDefinition: &design.UserTypeDefinition{
				AttributeDefinition: &design.AttributeDefinition{
					Type: design.Object{
						"att1": &design.AttributeDefinition{Type: design.Integer},
						"att2": &design.AttributeDefinition{Type: design.String},
					},
				},
				TypeName: "Foo",
			},
			Views: map[string]*design.ViewDefinition{
				"default": {
					Name: "default",
					AttributeDefinition: &design.AttributeDefinition{
						Type: design.Object{
							"att1": &design.AttributeDefinition{Type: design.String},
							"att2": &design.AttributeDefinition{Type: design.String},
						},
					},
				},
				"tiny": {
					Name: "tiny",
					AttributeDefinition: &design.AttributeDefinition{
						Type: design.Object{
							"att2": &design.AttributeDefinition{Type: design.String},
						},
					},
				},
			},
		}

		t.Run("using the empty view", func(t *testing.T) {
			if _, _, err := mt.Project(""); err == nil {
				t.Error("expected error")
			}
		})

		t.Run("using the default view", func(t *testing.T) {
			projected, _, err := mt.Project("default")
			if err != nil {
				t.Fatal(err)
			}

			// returns a media type with an identifier view param
			_, params, err := mime.ParseMediaType(projected.Identifier)
			if err != nil {
				t.Fatal(err)
			}
			if params["view"] != "default" {
				t.Errorf("invalid view: %s", params["view"])
			}

			// returns a media type with only a default view
			if len(projected.Views) != 1 {
				t.Errorf("invalid views: %v", projected.Views)
			}
			if _, ok := projected.Views["default"]; !ok {
				t.Error("missing default view")
			}

			// returns a media type with the default view attributes
			att := projected.Type.ToObject()["att1"]
			if att.Type.Kind() != design.IntegerKind {
				t.Errorf("invalid att1 type: %v", att.Type)
			}
		})

		t.Run("using the tiny view", func(t *testing.T) {
			projected, _, err := mt.Project("tiny")
			if err != nil {
				t.Fatal(err)
			}

			// returns a media type with an identifier view param
			_, params, err := mime.ParseMediaType(projected.Identifier)
			if err != nil {
				t.Fatal(err)
			}
			if params["view"] != "tiny" {
				t.Errorf("invalid view: %s", params["view"])
			}

			// returns a media type with only a tiny view
			if len(projected.Views) != 1 {
				t.Errorf("invalid views: %v", projected.Views)
			}
			if _, ok := projected.Views["default"]; !ok {
				t.Error("missing default view")
			}

			// returns a media type with the tiny view attributes
			att := projected.Type.ToObject()["att2"]
			if att.Type.Kind() != design.StringKind {
				t.Errorf("invalid att2 type: %v", att.Type)
			}
		})

		t.Run("on a collection", func(t *testing.T) {
			dslengine.Reset()
			mt := apidsl.CollectionOf(design.Dup(mt))
			dslengine.Execute(mt.DSL(), mt)
			mt.GenerateExample(design.NewRandomGenerator(""), nil)

			projected, _, err := mt.Project("default")
			if err != nil {
				t.Fatal(err)
			}
			// resets the example
			if projected.Example != nil {
				t.Error("unexpected example")
			}
		})
	})

	t.Run("with a media type with a links attribute", func(t *testing.T) {
		design.Design.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
		mt := &design.MediaTypeDefinition{
			UserTypeDefinition: &design.UserTypeDefinition{
				AttributeDefinition: &design.AttributeDefinition{
					Type: design.Object{
						"att1":  &design.AttributeDefinition{Type: design.Integer},
						"links": &design.AttributeDefinition{Type: design.String},
					},
				},
				TypeName: "Foo",
			},
			Identifier: "vnd.application/foo",
			Views: map[string]*design.ViewDefinition{
				"default": {
					Name: "default",
					AttributeDefinition: &design.AttributeDefinition{
						Type: design.Object{
							"att1":  &design.AttributeDefinition{Type: design.String},
							"links": &design.AttributeDefinition{Type: design.String},
						},
					},
				},
			},
		}

		projected, _, err := mt.Project("default")
		if err != nil {
			t.Fatal(err)
		}
		att := projected.Type.ToObject()["links"]
		if att.Type.Kind() != design.StringKind {
			t.Errorf("invalid links type: %v", att.Type)
		}
	})

	t.Run("with media types with view attributes with a cyclical dependency", func(t *testing.T) {
		design.Design.Reset()
		dslengine.Reset()
		design.ProjectedMediaTypes = make(design.MediaTypeRoot)

		apidsl.API("test", func() {})
		mt := apidsl.MediaType("vnd.application/MT1", func() {
			apidsl.TypeName("Mt1")
			apidsl.Attributes(func() {
				apidsl.Attribute("att", "vnd.application/MT2", func() {
					apidsl.Metadata("foo", "bar")
				})
			})
			apidsl.Links(func() {
				apidsl.Link("att", "default")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("att")
				apidsl.Attribute("links")
			})
			apidsl.View("tiny", func() {
				apidsl.Attribute("att", func() {
					apidsl.View("tiny")
				})
			})
		})
		apidsl.MediaType("vnd.application/MT2", func() {
			apidsl.TypeName("Mt2")
			apidsl.Attributes(func() {
				apidsl.Attribute("att2", mt)
			})
			apidsl.Links(func() {
				apidsl.Link("att2", "default")
			})
			apidsl.View("default", func() {
				apidsl.Attribute("att2")
				apidsl.Attribute("links")
			})
			apidsl.View("tiny", func() {
				apidsl.Attribute("links")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}
		if dslengine.Errors != nil {
			t.Fatal(dslengine.Errors)
		}

		// using the default view
		projected, links, err := mt.Project("default")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := projected.Type.ToObject()["att"]; !ok {
			t.Error("missing att attribute")
		}
		l := projected.Type.ToObject()["links"]
		if l.Type.(*design.UserTypeDefinition).AttributeDefinition != links.AttributeDefinition {
			t.Error("invalid links attribute")
		}
		metadata := links.Type.ToObject()["att"].Metadata
		want := dslengine.MetadataDefinition{"foo": []string{"bar"}}
		if diff := cmp.Diff(want, metadata); diff != "" {
			t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
		}
	})
}
