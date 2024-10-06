package codegen_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

// UserDefinitionType is a user defined type used in the unit tests.
type UserDefinitionType string

func TestFinalizer(t *testing.T) {
	// given a recursive user type with an array attribute
	rt := &design.UserTypeDefinition{TypeName: "recursive"}
	ar := &design.Array{ElemType: &design.AttributeDefinition{Type: rt}}
	obj := &design.Object{
		"elems": &design.AttributeDefinition{Type: ar},
		"other": &design.AttributeDefinition{
			Type:         design.String,
			DefaultValue: "foo",
		},
	}
	rt.AttributeDefinition = &design.AttributeDefinition{Type: obj}

	testCases := []struct {
		name   string
		att    *design.AttributeDefinition
		target string
		want   string
	}{

		{
			name: "given an object with a primitive field",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type:         design.String,
						DefaultValue: "bar",
					},
				},
			},
			target: "ut",
			want: `var defaultFoo string = "bar"
if ut.Foo == nil {
	ut.Foo = &defaultFoo
}`,
		},

		{
			name: "given an object with a primitive Number field",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type:         design.Number,
						DefaultValue: 0.0,
					},
				},
			},
			target: "ut",
			want: `var defaultFoo float64 = 0.000000
if ut.Foo == nil {
	ut.Foo = &defaultFoo
}`,
		},

		{
			name: "given an object with a primitive Number field with a int default value",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type:         design.Number,
						DefaultValue: 50,
					},
				},
			},
			target: "ut",
			want: `var defaultFoo float64 = 50.000000
if ut.Foo == nil {
	ut.Foo = &defaultFoo
}`,
		},

		{
			name: "given an array field",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Array{
							ElemType: &design.AttributeDefinition{
								Type: design.String,
							},
						},
						DefaultValue: []any{"bar", "baz"},
					},
				},
			},
			target: "ut",
			want: `if ut.Foo == nil {
	ut.Foo = []string{"bar", "baz"}
}`,
		},

		{
			name: "given a hash field",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Hash{
							KeyType: &design.AttributeDefinition{
								Type: design.String,
							},
							ElemType: &design.AttributeDefinition{
								Type: design.String,
							},
						},
						DefaultValue: map[any]any{"bar": "baz"},
					},
				},
			},
			target: "ut",
			want: `if ut.Foo == nil {
	ut.Foo = map[string]string{"bar": "baz"}
}`,
		},

		{
			name: "given a datetime field",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type:         design.DateTime,
						DefaultValue: interface{}("1978-06-30T10:00:00+09:00"),
					},
				},
			},
			target: "ut",
			want: `var defaultFoo, _ = time.Parse(time.RFC3339, "1978-06-30T10:00:00+09:00")
if ut.Foo == nil {
	ut.Foo = &defaultFoo
}`,
		},

		{
			name:   "given a recursive user type with an array attribute",
			att:    &design.AttributeDefinition{Type: rt},
			target: "ut",
			want: `	for _, e := range ut.Elems {
		var defaultOther string = "foo"
		if e.Other == nil {
			e.Other = &defaultOther
}
	}
var defaultOther string = "foo"
if ut.Other == nil {
	ut.Other = &defaultOther
}`,
		},

		{
			name: "given an object with a user definition type",
			att: &design.AttributeDefinition{
				Type: &design.Object{
					"foo": &design.AttributeDefinition{
						Type: design.String,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:type": []string{"UserDefinitionType", "github.com/shogo82148/shogoa/shogoagen/codegen_test"},
						},
						DefaultValue: UserDefinitionType("bar"),
					},
				},
			},
			target: "ut",
			want: `var defaultFoo UserDefinitionType = "bar"
if ut.Foo == nil {
	ut.Foo = &defaultFoo
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			finalizer := codegen.NewFinalizer()
			code := finalizer.Code(tc.att, tc.target, 0)
			if diff := cmp.Diff(tc.want, code); diff != "" {
				t.Errorf("unexpected code (-want +got):\n%s", diff)
			}
		})
	}
}
