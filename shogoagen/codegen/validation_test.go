package codegen_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/pointer"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

func TestValidator(t *testing.T) {
	t.Run("given an attribute definition and validations of enum", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Integer,
			Validation: &dslengine.ValidationDefinition{
				Values: []any{1, 2, 3},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if !(*val == 1 || *val == 2 || *val == 3) {
			err = shogoa.MergeErrors(err, shogoa.InvalidEnumValueError(` + "`context`" + `, *val, []interface{}{1, 2, 3}))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of pattern", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Integer,
			Validation: &dslengine.ValidationDefinition{
				Pattern: ".*",
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if ok := shogoa.ValidatePattern(` + "`.*`" + `, *val); !ok {
			err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`context`" + `, *val, ` + "`.*`" + `))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of min value 0", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Integer,
			Validation: &dslengine.ValidationDefinition{
				Minimum: pointer.Ptr(0.0),
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if *val < 0 {
			err = shogoa.MergeErrors(err, shogoa.InvalidRangeError(` + "`" + `context` + "`" + `, *val, 0, true))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of max value math.MaxInt64", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Integer,
			Validation: &dslengine.ValidationDefinition{
				Maximum: pointer.Ptr(float64(math.MaxInt64)),
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if *val > 9223372036854775807 {
			err = shogoa.MergeErrors(err, shogoa.InvalidRangeError(` + "`" + `context` + "`" + `, *val, 9223372036854775807, false))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of min value math.MinInt64", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Integer,
			Validation: &dslengine.ValidationDefinition{
				Minimum: pointer.Ptr(float64(math.MinInt64)),
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if *val < -9223372036854775808 {
			err = shogoa.MergeErrors(err, shogoa.InvalidRangeError(` + "`" + `context` + "`" + `, *val, -9223372036854775808, true))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of array min length 1", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: &design.Array{
				ElemType: &design.AttributeDefinition{
					Type: design.String,
				},
			},
			Validation: &dslengine.ValidationDefinition{
				MinLength: pointer.Ptr(1),
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if len(val) < 1 {
			err = shogoa.MergeErrors(err, shogoa.InvalidLengthError(` + "`" + `context` + "`" + `, val, len(val), 1, true))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of array elements", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: &design.Array{
				ElemType: &design.AttributeDefinition{
					Type: design.String,
					Validation: &dslengine.ValidationDefinition{
						Pattern: ".*",
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	for _, e := range val {
		if ok := shogoa.ValidatePattern(` + "`" + `.*` + "`" + `, e); !ok {
			err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`" + `context[*]` + "`" + `, e, ` + "`" + `.*` + "`" + `))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of hash elements (key, elem)", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: &design.Hash{
				KeyType: &design.AttributeDefinition{
					Type: design.String,
					Validation: &dslengine.ValidationDefinition{
						Pattern: ".*",
					},
				},
				ElemType: &design.AttributeDefinition{
					Type: design.String,
					Validation: &dslengine.ValidationDefinition{
						Pattern: ".*",
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	for k, e := range val {
		if ok := shogoa.ValidatePattern(` + "`" + `.*` + "`" + `, k); !ok {
			err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`" + `context[*]` + "`" + `, k, ` + "`" + `.*` + "`" + `))
		}
		if ok := shogoa.ValidatePattern(` + "`" + `.*` + "`" + `, e); !ok {
			err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`" + `context[*]` + "`" + `, e, ` + "`" + `.*` + "`" + `))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of hash elements (key, _)", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: &design.Hash{
				KeyType: &design.AttributeDefinition{
					Type: design.String,
					Validation: &dslengine.ValidationDefinition{
						Pattern: ".*",
					},
				},
				ElemType: &design.AttributeDefinition{
					Type: design.String,
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	for k, _ := range val {
		if ok := shogoa.ValidatePattern(` + "`" + `.*` + "`" + `, k); !ok {
			err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`" + `context[*]` + "`" + `, k, ` + "`" + `.*` + "`" + `))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of hash elements (key, _)", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: &design.Hash{
				KeyType: &design.AttributeDefinition{
					Type: design.String,
				},
				ElemType: &design.AttributeDefinition{
					Type: design.String,
					Validation: &dslengine.ValidationDefinition{
						Pattern: ".*",
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	for _, e := range val {
		if ok := shogoa.ValidatePattern(` + "`" + `.*` + "`" + `, e); !ok {
			err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`" + `context[*]` + "`" + `, e, ` + "`" + `.*` + "`" + `))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of string min length 2", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.String,
			Validation: &dslengine.ValidationDefinition{
				MinLength: pointer.Ptr(2),
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val != nil {
		if utf8.RuneCountInString(*val) < 2 {
			err = shogoa.MergeErrors(err, shogoa.InvalidLengthError(` + "`" + `context` + "`" + `, *val, utf8.RuneCountInString(*val), 2, true))
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of embedded object and the parent is optional", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Object{
				"foo": &design.AttributeDefinition{
					Type: design.Object{
						"bar": &design.AttributeDefinition{
							Type: design.Integer,
							Validation: &dslengine.ValidationDefinition{
								Values: []interface{}{1, 2, 3},
							},
						},
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val.Foo != nil {
		if val.Foo.Bar != nil {
			if !(*val.Foo.Bar == 1 || *val.Foo.Bar == 2 || *val.Foo.Bar == 3) {
				err = shogoa.MergeErrors(err, shogoa.InvalidEnumValueError(` + "`" + `context.foo.bar` + "`" + `, *val.Foo.Bar, []interface{}{1, 2, 3}))
			}
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of embedded object and the parent is required", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Object{
				"foo": &design.AttributeDefinition{
					Type: design.Object{
						"bar": &design.AttributeDefinition{
							Type: design.Integer,
							Validation: &dslengine.ValidationDefinition{
								Values: []interface{}{1, 2, 3},
							},
						},
					},
				},
			},
			Validation: &dslengine.ValidationDefinition{
				Required: []string{"foo"},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val.Foo == nil {
		err = shogoa.MergeErrors(err, shogoa.MissingAttributeError(` + "`context`" + `, "foo"))
	}
	if val.Foo != nil {
		if val.Foo.Bar != nil {
			if !(*val.Foo.Bar == 1 || *val.Foo.Bar == 2 || *val.Foo.Bar == 3) {
				err = shogoa.MergeErrors(err, shogoa.InvalidEnumValueError(` + "`" + `context.foo.bar` + "`" + `, *val.Foo.Bar, []interface{}{1, 2, 3}))
			}
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of embedded object with a child attribute with struct:tag:name metadata", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Object{
				"foo": &design.AttributeDefinition{
					Type: design.Object{
						"bar": &design.AttributeDefinition{
							Type: design.Integer,
							Validation: &dslengine.ValidationDefinition{
								Values: []interface{}{1, 2, 3},
							},
						},
					},
					Metadata: dslengine.MetadataDefinition{
						"struct:field:name": []string{"FOO"},
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val.FOO != nil {
		if val.FOO.Bar != nil {
			if !(*val.FOO.Bar == 1 || *val.FOO.Bar == 2 || *val.FOO.Bar == 3) {
				err = shogoa.MergeErrors(err, shogoa.InvalidEnumValueError(` + "`" + `context.foo.bar` + "`" + `, *val.FOO.Bar, []interface{}{1, 2, 3}))
			}
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of embedded object with a child attribute with a grand child attribute with struct:tag:name metadata", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Object{
				"foo": &design.AttributeDefinition{
					Type: design.Object{
						"bar": &design.AttributeDefinition{
							Type: design.Integer,
							Validation: &dslengine.ValidationDefinition{
								Values: []interface{}{1, 2, 3},
							},
							Metadata: dslengine.MetadataDefinition{
								"struct:field:name": []string{"FOO"},
							},
						},
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val.Foo != nil {
		if val.Foo.FOO != nil {
			if !(*val.Foo.FOO == 1 || *val.Foo.FOO == 2 || *val.Foo.FOO == 3) {
				err = shogoa.MergeErrors(err, shogoa.InvalidEnumValueError(` + "`" + `context.foo.bar` + "`" + `, *val.Foo.FOO, []interface{}{1, 2, 3}))
			}
		}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of required user type attribute with no validation with only direct required attributes", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Object{
				"foo": &design.AttributeDefinition{
					Type: &design.Array{
						ElemType: &design.AttributeDefinition{
							Type: &design.UserTypeDefinition{
								TypeName: "UT",
								AttributeDefinition: &design.AttributeDefinition{
									Type: design.Object{
										"bar": &design.AttributeDefinition{Type: design.String},
									},
								},
							},
						},
					},
				},
				"foo2": &design.AttributeDefinition{
					Type: &design.UserTypeDefinition{
						TypeName: "UT",
						AttributeDefinition: &design.AttributeDefinition{
							Type: design.Object{
								"bar": &design.AttributeDefinition{Type: design.String},
							},
						},
					},
				},
			},
			Validation: &dslengine.ValidationDefinition{
				Required: []string{"foo"},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val.Foo == nil {
		err = shogoa.MergeErrors(err, shogoa.MissingAttributeError(` + "`context`" + `, "foo"))
	}
	if val.Foo2 != nil {
	if err2 := val.Foo2.Validate(); err2 != nil {
		err = shogoa.MergeErrors(err, err2)
	}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

	t.Run("given an attribute definition and validations of required user type attribute with no validation with only direct required attributes", func(t *testing.T) {
		att := &design.AttributeDefinition{
			Type: design.Object{
				"foo": &design.AttributeDefinition{
					Type: &design.Array{
						ElemType: &design.AttributeDefinition{
							Type: &design.UserTypeDefinition{
								TypeName: "UT",
								AttributeDefinition: &design.AttributeDefinition{
									Type: design.Object{
										"bar": &design.AttributeDefinition{Type: design.String},
									},
									Validation: &dslengine.ValidationDefinition{
										Required: []string{"bar"},
									},
								},
							},
						},
					},
				},
				"foo2": &design.AttributeDefinition{
					Type: &design.UserTypeDefinition{
						TypeName: "UT",
						AttributeDefinition: &design.AttributeDefinition{
							Type: design.Object{
								"bar": &design.AttributeDefinition{Type: design.String},
							},
						},
					},
				},
			},
		}
		got := codegen.NewValidator().Code(att, false, false, false, "val", "context", 1, false)
		want := `	if val.Foo2 != nil {
	if err2 := val.Foo2.Validate(); err2 != nil {
		err = shogoa.MergeErrors(err, err2)
	}
	}`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected code (-want +got):\n%s", diff)
		}
	})

}
