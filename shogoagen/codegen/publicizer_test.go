package codegen_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

func TestPublicizer(t *testing.T) {
	tests := []struct {
		name        string
		att         *design.AttributeDefinition
		sourceField string
		targetField string
		init        bool
		want        string
	}{
		{
			name:        "given a simple field",
			att:         &design.AttributeDefinition{Type: design.Integer},
			sourceField: "source",
			targetField: "target",
			want:        "target = source",
		},

		{
			name:        "given a simple field with init true",
			att:         &design.AttributeDefinition{Type: design.Integer},
			sourceField: "source",
			targetField: "target",
			init:        true,
			want:        "target := source",
		},

		{
			name: "given an object field",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.String},
					"bar": &design.AttributeDefinition{Type: design.Any},
					"baz": &design.AttributeDefinition{Type: design.Any},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"bar"},
				},
			},
			sourceField: "source",
			targetField: "target",
			want: `target = &struct {
	Bar interface{} ` + "`" + `form:"bar" json:"bar" yaml:"bar" xml:"bar"` + "`" + `
	Baz interface{} ` + "`" + `form:"baz,omitempty" json:"baz,omitempty" yaml:"baz,omitempty" xml:"baz,omitempty"` + "`" + `
	Foo *string ` + "`" + `form:"foo,omitempty" json:"foo,omitempty" yaml:"foo,omitempty" xml:"foo,omitempty"` + "`" + `
}{}
if source.Bar != nil {
	target.Bar = source.Bar
}
if source.Baz != nil {
	target.Baz = source.Baz
}
if source.Foo != nil {
	target.Foo = source.Foo
}`,
		},

		{
			name: "given an object field with metadata",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type:     design.String,
						Metadata: dslengine.MetadataDefinition{"struct:field:name": []string{"MetaFoo"}},
					},
					"bar": &design.AttributeDefinition{
						Type:     design.Any,
						Metadata: dslengine.MetadataDefinition{"struct:field:name": []string{"MetaBar"}},
					},
					"baz": &design.AttributeDefinition{
						Type:     design.Any,
						Metadata: dslengine.MetadataDefinition{"struct:field:name": []string{"MetaBaz"}},
					},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"bar"},
				},
			},
			sourceField: "source",
			targetField: "target",
			want: `target = &struct {
	MetaBar interface{} ` + "`" + `form:"bar" json:"bar" yaml:"bar" xml:"bar"` + "`" + `
	MetaBaz interface{} ` + "`" + `form:"baz,omitempty" json:"baz,omitempty" yaml:"baz,omitempty" xml:"baz,omitempty"` + "`" + `
	MetaFoo *string ` + "`" + `form:"foo,omitempty" json:"foo,omitempty" yaml:"foo,omitempty" xml:"foo,omitempty"` + "`" + `
}{}
if source.MetaBar != nil {
	target.MetaBar = source.MetaBar
}
if source.MetaBaz != nil {
	target.MetaBaz = source.MetaBaz
}
if source.MetaFoo != nil {
	target.MetaFoo = source.MetaFoo
}`,
		},

		{
			name: "given a user type",
			att: &design.AttributeDefinition{
				Type: &design.UserTypeDefinition{
					AttributeDefinition: &design.AttributeDefinition{
						Type: &design.Object{
							"foo": &design.AttributeDefinition{Type: design.String},
						},
					},
					TypeName: "TheUserType",
				},
			},
			sourceField: "source",
			targetField: "target",
			want:        `target = source.Publicize()`,
		},

		{
			name: "given an array field that contains primitive fields",
			att: &design.AttributeDefinition{
				Type: &design.Array{
					ElemType: &design.AttributeDefinition{
						Type: design.String,
					},
				},
			},
			sourceField: "source",
			targetField: "target",
			want:        `target = source`,
		},

		{
			name: "given an array field that contains user defined fields",
			att: &design.AttributeDefinition{
				Type: &design.Array{
					ElemType: &design.AttributeDefinition{
						Type: &design.UserTypeDefinition{
							AttributeDefinition: &design.AttributeDefinition{
								Type: design.Object{
									"foo": &design.AttributeDefinition{Type: design.String},
								},
							},
							TypeName: "TheUserType",
						},
					},
				},
			},
			sourceField: "source",
			targetField: "target",
			want: `target = make([]*TheUserType, len(source))
for i0, elem0 := range source {
	target[i0] = elem0.Publicize()
}`,
		},

		{
			name: "given a hash field that contains primitive fields",
			att: &design.AttributeDefinition{
				Type: &design.Hash{
					KeyType: &design.AttributeDefinition{
						Type: design.String,
					},
					ElemType: &design.AttributeDefinition{
						Type: design.String,
					},
				},
			},
			sourceField: "source",
			targetField: "target",
			want:        "target = source",
		},

		{
			name: "given a hash field that contains user defined fields",
			att: &design.AttributeDefinition{
				Type: &design.Hash{
					KeyType: &design.AttributeDefinition{
						Type: &design.UserTypeDefinition{
							AttributeDefinition: &design.AttributeDefinition{
								Type: &design.Object{
									"foo": &design.AttributeDefinition{Type: design.String},
								},
							},
							TypeName: "TheKeyType",
						},
					},
					ElemType: &design.AttributeDefinition{
						Type: &design.UserTypeDefinition{
							AttributeDefinition: &design.AttributeDefinition{
								Type: &design.Object{
									"bar": &design.AttributeDefinition{Type: design.String},
								},
							},
							TypeName: "TheElemType",
						},
					},
				},
			},
			sourceField: "source",
			targetField: "target",
			want: `target = make(map[*TheKeyType]*TheElemType, len(source))
for k0, v0 := range source {
	var pubk0 *TheKeyType
	if k0 != nil {
		pubk0 = k0.Publicize()
	}
	var pubv0 *TheElemType
	if v0 != nil {
		pubv0 = v0.Publicize()
	}
	target[pubk0] = pubv0
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := codegen.Publicizer(tt.att, tt.sourceField, tt.targetField, false, 0, tt.init)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected code (-want +got):\n%s", diff)
			}
		})
	}
}
