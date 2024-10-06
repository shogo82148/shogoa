package codegen_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

func TestAttributeImports(t *testing.T) {
	tests := []struct {
		name string
		att  *design.AttributeDefinition
		want []*codegen.ImportSpec
	}{
		{
			name: "given an attribute definition with fields of objects",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: design.String,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:type": []string{"json.RawMessage", "encoding/json"},
						},
					},
					"bar": &design.AttributeDefinition{Type: design.Integer},
				},
			},
			want: []*codegen.ImportSpec{
				{Path: "encoding/json"},
			},
		},

		{
			name: "given an attribute definition with fields of recursive objects",
			att: func() *design.AttributeDefinition {
				o := design.Object{
					"foo": &design.AttributeDefinition{Type: design.String},
				}
				o["foo"].Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				child := &design.AttributeDefinition{Type: o}

				po := design.Object{"child": child}
				po["child"].Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				parent := &design.AttributeDefinition{Type: po}

				o["parent"] = parent
				return &design.AttributeDefinition{
					Type: po,
				}
			}(),
			want: []*codegen.ImportSpec{
				{Path: "encoding/json"},
			},
		},

		{
			name: "given an attribute definition with fields of hash",
			att: &design.AttributeDefinition{
				Type: &design.Hash{
					KeyType: &design.AttributeDefinition{
						Type: design.Integer,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:type": []string{"json.RawMessage", "encoding/json"},
						},
					},
					ElemType: &design.AttributeDefinition{
						Type: design.Integer,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:type": []string{"json.RawMessage", "encoding/json"},
						},
					},
				},
			},
			want: []*codegen.ImportSpec{
				{Path: "encoding/json"},
			},
		},

		{
			name: "given an attribute definition with fields of array",
			att: &design.AttributeDefinition{
				Type: &design.Array{
					ElemType: &design.AttributeDefinition{
						Type: design.Integer,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:type": []string{"json.RawMessage", "encoding/json"},
						},
					},
				},
			},
			want: []*codegen.ImportSpec{
				{Path: "encoding/json"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var imports []*codegen.ImportSpec
			imports = codegen.AttributeImports(tt.att, imports, nil)
			if diff := cmp.Diff(tt.want, imports); diff != "" {
				t.Errorf("unexpected imports (-want +got):\n%s", diff)
			}
		})
	}
}
