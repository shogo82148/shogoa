package codegen_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/design/apidsl"
	"github.com/shogo82148/shogoa/dslengine"
	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

func TestGoify(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		firstUpper bool
		want       string
	}{
		{
			name:       "with first upper false",
			input:      "blue_id",
			firstUpper: false,
			want:       "blueID",
		},
		{
			name:       "with first upper false normal identifier",
			input:      "blue",
			firstUpper: false,
			want:       "blue",
		},
		{
			name:       "with first upper false and UUID",
			input:      "blue_uuid",
			firstUpper: false,
			want:       "blueUUID",
		},
		{
			name:       "with first upper true",
			input:      "blue_id",
			firstUpper: true,
			want:       "BlueID",
		},
		{
			name:       "with first upper true and UUID",
			input:      "blue_uuid",
			firstUpper: true,
			want:       "BlueUUID",
		},
		{
			name:       "with first upper true normal identifier",
			input:      "blue",
			firstUpper: true,
			want:       "Blue",
		},
		{
			name:       "with first upper false normal identifier",
			input:      "blue",
			firstUpper: false,
			want:       "blue",
		},
		{
			name:       "with first upper true normal identifier",
			input:      "Blue",
			firstUpper: true,
			want:       "Blue",
		},
		{
			name:       "with invalid identifier",
			input:      "Blue%50",
			firstUpper: true,
			want:       "Blue50",
		},
		{
			name:       "with invalid identifier first upper false",
			input:      "Blue%50",
			firstUpper: false,
			want:       "blue50",
		},
		{
			name:       "with only UUID and first upper false",
			input:      "UUID",
			firstUpper: false,
			want:       "uuid",
		},
		{
			name:       "with connectives invalid identifiers, first upper false",
			input:      "[[fields___type]]",
			firstUpper: false,
			want:       "fieldsType",
		},
		{
			name:       "with connectives invalid identifiers, first upper true",
			input:      "[[fields___type]]",
			firstUpper: true,
			want:       "FieldsType",
		},
		{
			name:       "with all invalid identifiers",
			input:      "[[",
			firstUpper: false,
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := codegen.Goify(tt.input, tt.firstUpper)
			if got != tt.want {
				t.Errorf("Goify(%v) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGoTypeDef(t *testing.T) {
	tests := []struct {
		name    string
		att     *design.AttributeDefinition
		private bool
		want    string
	}{
		{
			name: "given an attribute definition with fields of primitive types",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.Integer},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
			},
			want: "struct {\n" +
				"	Bar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	Baz *time.Time `form:\"baz,omitempty\" json:\"baz,omitempty\" yaml:\"baz,omitempty\" xml:\"baz,omitempty\"`\n" +
				"	Foo *int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"	Qux *uuid.UUID `form:\"qux,omitempty\" json:\"qux,omitempty\" yaml:\"qux,omitempty\" xml:\"qux,omitempty\"`\n" +
				"	Quz interface{} `form:\"quz,omitempty\" json:\"quz,omitempty\" yaml:\"quz,omitempty\" xml:\"quz,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of primitive types using struct tags metadata",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: design.Integer,
						Metadata: dslengine.MetadataDefinition{
							"struct:tag:foo":  []string{"bar", "baz"},
							"struct:tag:foo2": []string{"bar2"},
						},
					},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
			},
			want: "struct {\n" +
				"	Bar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	Baz *time.Time `form:\"baz,omitempty\" json:\"baz,omitempty\" yaml:\"baz,omitempty\" xml:\"baz,omitempty\"`\n" +
				"	Foo *int `foo:\"bar,baz\" foo2:\"bar2\"`\n" +
				"	Qux *uuid.UUID `form:\"qux,omitempty\" json:\"qux,omitempty\" yaml:\"qux,omitempty\" xml:\"qux,omitempty\"`\n" +
				"	Quz interface{} `form:\"quz,omitempty\" json:\"quz,omitempty\" yaml:\"quz,omitempty\" xml:\"quz,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of primitive types using struct field name metadata",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: design.Integer,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:name": []string{"serviceName", "unused"},
						},
					},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
			},
			want: "struct {\n" +
				"	Bar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	Baz *time.Time `form:\"baz,omitempty\" json:\"baz,omitempty\" yaml:\"baz,omitempty\" xml:\"baz,omitempty\"`\n" +
				"	ServiceName *int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"	Qux *uuid.UUID `form:\"qux,omitempty\" json:\"qux,omitempty\" yaml:\"qux,omitempty\" xml:\"qux,omitempty\"`\n" +
				"	Quz interface{} `form:\"quz,omitempty\" json:\"quz,omitempty\" yaml:\"quz,omitempty\" xml:\"quz,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of primitive types using struct field type metadata",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: design.Integer,
						Metadata: dslengine.MetadataDefinition{
							"struct:field:type": []string{"[]byte"},
						},
					},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
			},
			want: "struct {\n" +
				"	Bar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	Baz *time.Time `form:\"baz,omitempty\" json:\"baz,omitempty\" yaml:\"baz,omitempty\" xml:\"baz,omitempty\"`\n" +
				"	Foo *[]byte `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"	Qux *uuid.UUID `form:\"qux,omitempty\" json:\"qux,omitempty\" yaml:\"qux,omitempty\" xml:\"qux,omitempty\"`\n" +
				"	Quz interface{} `form:\"quz,omitempty\" json:\"quz,omitempty\" yaml:\"quz,omitempty\" xml:\"quz,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of primitive types that are required",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.Integer},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"foo", "bar", "baz", "qux", "quz"},
				},
			},
			want: "struct {\n" +
				"	Bar string `form:\"bar\" json:\"bar\" yaml:\"bar\" xml:\"bar\"`\n" +
				"	Baz time.Time `form:\"baz\" json:\"baz\" yaml:\"baz\" xml:\"baz\"`\n" +
				"	Foo int `form:\"foo\" json:\"foo\" yaml:\"foo\" xml:\"foo\"`\n" +
				"	Qux uuid.UUID `form:\"qux\" json:\"qux\" yaml:\"qux\" xml:\"qux\"`\n" +
				"	Quz interface{} `form:\"quz\" json:\"quz\" yaml:\"quz\" xml:\"quz\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of hash of primitive types",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Hash{
							KeyType:  &design.AttributeDefinition{Type: design.Integer},
							ElemType: &design.AttributeDefinition{Type: design.Integer},
						},
					},
				},
			},
			want: "struct {\n" +
				"\tFoo map[int]int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of array of primitive types",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Array{
							ElemType: &design.AttributeDefinition{Type: design.Integer},
						},
					},
				},
			},
			want: "struct {\n" +
				"\tFoo []int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of hash of objects",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Hash{
							KeyType: &design.AttributeDefinition{
								Type: design.Object{
									"keyAtt": &design.AttributeDefinition{Type: design.String},
								},
							},
							ElemType: &design.AttributeDefinition{
								Type: design.Object{
									"elemAtt": &design.AttributeDefinition{Type: design.Integer},
								},
							},
						},
					},
				},
			},
			want: "struct {\n" +
				"	Foo map[*struct {\n" +
				"		KeyAtt *string `form:\"keyAtt,omitempty\" json:\"keyAtt,omitempty\" yaml:\"keyAtt,omitempty\" xml:\"keyAtt,omitempty\"`\n" +
				"	}]*struct {\n" +
				"		ElemAtt *int `form:\"elemAtt,omitempty\" json:\"elemAtt,omitempty\" yaml:\"elemAtt,omitempty\" xml:\"elemAtt,omitempty\"`\n" +
				"	} `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of array of objects",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Array{
							ElemType: &design.AttributeDefinition{
								Type: design.Object{
									"bar": &design.AttributeDefinition{Type: design.Integer},
								},
							},
						},
					},
				},
			},
			want: "struct {\n" +
				"	Foo []*struct {\n" +
				"		Bar *int `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	} `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields of array of objects that are required",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Array{
							ElemType: &design.AttributeDefinition{
								Type: design.Object{
									"bar": &design.AttributeDefinition{Type: design.Integer},
								},
							},
						},
					},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"foo"},
				},
			},
			want: "struct {\n" +
				"	Foo []*struct {\n" +
				"		Bar *int `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	} `form:\"foo\" json:\"foo\" yaml:\"foo\" xml:\"foo\"`\n" +
				"}",
		},

		{
			name: "given an attribute definition with fields that are required",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.Integer},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"foo"},
				},
			},
			want: "struct {\n" +
				"	Foo int `form:\"foo\" json:\"foo\" yaml:\"foo\" xml:\"foo\"`\n" +
				"}",
		},

		{
			name: "given an array of primitive type",
			att: &design.AttributeDefinition{
				Type: &design.Array{
					ElemType: &design.AttributeDefinition{
						Type: design.Integer,
					},
				},
			},
			want: "[]int",
		},

		{
			name: "given an array of object type",
			att: &design.AttributeDefinition{
				Type: &design.Array{
					ElemType: &design.AttributeDefinition{
						Type: design.Object{
							"foo": &design.AttributeDefinition{Type: design.Integer},
							"bar": &design.AttributeDefinition{Type: design.String},
						},
					},
				},
			},
			want: "[]*struct {\n" +
				"\tBar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"\tFoo *int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},

		{
			name: "when generating an all-optional private struct, given an attribute definition with fields of primitive types",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.Integer},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
			},
			private: true,
			want: "struct {\n" +
				"	Bar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	Baz *time.Time `form:\"baz,omitempty\" json:\"baz,omitempty\" yaml:\"baz,omitempty\" xml:\"baz,omitempty\"`\n" +
				"	Foo *int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"	Qux *uuid.UUID `form:\"qux,omitempty\" json:\"qux,omitempty\" yaml:\"qux,omitempty\" xml:\"qux,omitempty\"`\n" +
				"	Quz interface{} `form:\"quz,omitempty\" json:\"quz,omitempty\" yaml:\"quz,omitempty\" xml:\"quz,omitempty\"`\n" +
				"}",
		},

		{
			name: "when generating an all-optional private struct, given an attribute definition with fields of primitive types that are required",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.Integer},
					"bar": &design.AttributeDefinition{Type: design.String},
					"baz": &design.AttributeDefinition{Type: design.DateTime},
					"qux": &design.AttributeDefinition{Type: design.UUID},
					"quz": &design.AttributeDefinition{Type: design.Any},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"foo", "bar", "baz", "qux", "quz"},
				},
			},
			private: true,
			want: "struct {\n" +
				"	Bar *string `form:\"bar,omitempty\" json:\"bar,omitempty\" yaml:\"bar,omitempty\" xml:\"bar,omitempty\"`\n" +
				"	Baz *time.Time `form:\"baz,omitempty\" json:\"baz,omitempty\" yaml:\"baz,omitempty\" xml:\"baz,omitempty\"`\n" +
				"	Foo *int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"	Qux *uuid.UUID `form:\"qux,omitempty\" json:\"qux,omitempty\" yaml:\"qux,omitempty\" xml:\"qux,omitempty\"`\n" +
				"	Quz interface{} `form:\"quz,omitempty\" json:\"quz,omitempty\" yaml:\"quz,omitempty\" xml:\"quz,omitempty\"`\n" +
				"}",
		},

		{
			name: "when generating an all-optional private struct, given an attribute definition with fields of hash of primitive types",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Hash{
							KeyType:  &design.AttributeDefinition{Type: design.Integer},
							ElemType: &design.AttributeDefinition{Type: design.Integer},
						},
					},
				},
			},
			private: true,
			want: "struct {\n" +
				"\tFoo map[int]int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`" +
				"\n}",
		},

		{
			name: "when generating an all-optional private struct, given an attribute definition with fields of array of primitive types",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{
						Type: &design.Array{
							ElemType: &design.AttributeDefinition{Type: design.Integer},
						},
					},
				},
			},
			private: true,
			want: "struct {\n" +
				"\tFoo []int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},

		{
			name: "when generating an all-optional private struct, given an attribute definition with fields that are required",
			att: &design.AttributeDefinition{
				Type: design.Object{
					"foo": &design.AttributeDefinition{Type: design.Integer},
				},
				Validation: &dslengine.ValidationDefinition{
					Required: []string{"foo"},
				},
			},
			private: true,
			want: "struct {\n" +
				"	Foo *int `form:\"foo,omitempty\" json:\"foo,omitempty\" yaml:\"foo,omitempty\" xml:\"foo,omitempty\"`\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := codegen.GoTypeDef(tt.att, 0, true, tt.private)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected code (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGoTypeTransform(t *testing.T) {
	t.Run("transforming simple objects", func(t *testing.T) {
		dslengine.Reset()
		source := apidsl.Type("Source", func() {
			apidsl.Attribute("att")
		})
		target := apidsl.Type("Target", func() {
			apidsl.Attribute("att")
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		got, err := codegen.GoTypeTransform(source, target, "", "Transform")
		if err != nil {
			t.Fatal(err)
		}
		want := `func Transform(source *Source) (target *Target) {
	target = new(Target)
	target.Att = source.Att
	return
}
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("transforming objects with attributes with map key metadata", func(t *testing.T) {
		dslengine.Reset()
		source := apidsl.Type("Source", func() {
			apidsl.Attribute("foo", func() {
				apidsl.Metadata(codegen.TransformMapKey, "key")
			})
		})
		target := apidsl.Type("Target", func() {
			apidsl.Attribute("bar", func() {
				apidsl.Metadata(codegen.TransformMapKey, "key")
			})
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		got, err := codegen.GoTypeTransform(source, target, "", "Transform")
		if err != nil {
			t.Fatal(err)
		}
		want := `func Transform(source *Source) (target *Target) {
	target = new(Target)
	target.Bar = source.Foo
	return
}
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("transforming objects with array attributes", func(t *testing.T) {
		dslengine.Reset()
		source := apidsl.Type("Source", func() {
			apidsl.Attribute("att", apidsl.ArrayOf(design.Integer))
		})
		target := apidsl.Type("Target", func() {
			apidsl.Attribute("att", apidsl.ArrayOf(design.Integer))
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		got, err := codegen.GoTypeTransform(source, target, "", "Transform")
		if err != nil {
			t.Fatal(err)
		}
		want := `func Transform(source *Source) (target *Target) {
	target = new(Target)
	target.Att = make([]int, len(source.Att))
	for i, v := range source.Att {
		target.Att[i] = source.Att[i]
	}
	return
}
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("transforming objects with hash attributes", func(t *testing.T) {
		dslengine.Reset()
		elem := apidsl.Type("elem", func() {
			apidsl.Attribute("foo", design.Integer)
			apidsl.Attribute("bar")
		})
		source := apidsl.Type("Source", func() {
			apidsl.Attribute("att", apidsl.HashOf(design.String, elem))
		})
		target := apidsl.Type("Target", func() {
			apidsl.Attribute("att", apidsl.HashOf(design.String, elem))
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		got, err := codegen.GoTypeTransform(source, target, "", "Transform")
		if err != nil {
			t.Fatal(err)
		}
		want := `func Transform(source *Source) (target *Target) {
	target = new(Target)
	target.Att = make(map[string]*Elem, len(source.Att))
	for k, v := range source.Att {
		var tk string
		tk = k
		var tv *Elem
		tv = new(Elem)
		tv.Bar = v.Bar
		tv.Foo = v.Foo
		target.Att[tk] = tv
	}
	return
}
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("transforming objects with recursive attributes", func(t *testing.T) {
		dslengine.Reset()
		inner := apidsl.Type("inner", func() {
			apidsl.Attribute("foo", design.Integer)
		})
		outer := apidsl.Type("outer", func() {
			apidsl.Attribute("in", inner)
		})
		array := apidsl.Type("array", func() {
			apidsl.Attribute("elem", apidsl.ArrayOf(outer))
		})
		hash := apidsl.Type("hash", func() {
			apidsl.Attribute("elem", apidsl.HashOf(design.Integer, outer))
		})
		source := apidsl.Type("Source", func() {
			apidsl.Attribute("outer", outer)
			apidsl.Attribute("array", array)
			apidsl.Attribute("hash", hash)
		})
		target := apidsl.Type("Target", func() {
			apidsl.Attribute("outer", outer)
			apidsl.Attribute("array", array)
			apidsl.Attribute("hash", hash)
		})
		if err := dslengine.Run(); err != nil {
			t.Fatal(err)
		}

		got, err := codegen.GoTypeTransform(source, target, "", "Transform")
		if err != nil {
			t.Fatal(err)
		}
		want := `func Transform(source *Source) (target *Target) {
	target = new(Target)
	target.Array = new(Array)
	target.Array.Elem = make([]*Outer, len(source.Array.Elem))
	for i, v := range source.Array.Elem {
		target.Array.Elem[i] = new(Outer)
		target.Array.Elem[i].In = new(Inner)
		target.Array.Elem[i].In.Foo = source.Array.Elem[i].In.Foo
	}
	target.Hash = new(Hash)
	target.Hash.Elem = make(map[int]*Outer, len(source.Hash.Elem))
	for k, v := range source.Hash.Elem {
		var tk int
		tk = k
		var tv *Outer
		tv = new(Outer)
		tv.In = new(Inner)
		tv.In.Foo = v.In.Foo
		target.Hash.Elem[tk] = tv
	}
	target.Outer = new(Outer)
	target.Outer.In = new(Inner)
	target.Outer.In.Foo = source.Outer.In.Foo
	return
}
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})
}

func TestGoTypeDesc(t *testing.T) {
	t.Run("with a type with a description", func(t *testing.T) {
		ut := &design.UserTypeDefinition{
			AttributeDefinition: &design.AttributeDefinition{
				Description: "foo",
			},
		}
		got := codegen.GoTypeDesc(ut, false)
		want := "foo"
		if got != want {
			t.Errorf("GoTypeDesc() = %v; want %v", got, want)
		}
	})

	t.Run("with a type with a description containing newlines", func(t *testing.T) {
		ut := &design.UserTypeDefinition{
			AttributeDefinition: &design.AttributeDefinition{
				Description: "foo\nbar",
			},
		}
		got := codegen.GoTypeDesc(ut, false)
		want := "foo\n// bar"
		if got != want {
			t.Errorf("GoTypeDesc() = %v; want %v", got, want)
		}
	})
}
