package design

import (
	"slices"
	"testing"

	"github.com/shogo82148/shogoa/dslengine"
)

func TestCanonicalIdentifier(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{
			id:   "application/json",
			want: "application/json",
		},
		{
			id:   "application/json+xml; foo=bar",
			want: "application/json; foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := CanonicalIdentifier(tt.id)
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestExtractWildcards(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{
			path: "/foo",
			want: nil,
		},
		{
			path: "/a/:foo/:bar/b/:baz/c",
			want: []string{"foo", "bar", "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := ExtractWildcards(tt.path)
			if !slices.Equal(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestMediaTypeRoot_DSLName(t *testing.T) {
	root := MediaTypeRoot{}
	Design.MediaTypes = make(map[string]*MediaTypeDefinition)

	if root.DSLName() == "" {
		t.Error("DSLName is empty")
	}
}

func TestMediaTypeRoot_DependsOn(t *testing.T) {
	root := MediaTypeRoot{}
	Design.MediaTypes = make(map[string]*MediaTypeDefinition)
	if !slices.Equal(root.DependsOn(), []dslengine.Root{Design}) {
		t.Error("DependsOn is invalid")
	}
}

func TestMediaTypeRoot_IterateSets(t *testing.T) {
	tests := []struct {
		name string
		root MediaTypeRoot
		set  []string
	}{
		{
			name: "empty",
			root: MediaTypeRoot{},
			set:  []string{},
		},
		{
			name: "one",
			root: MediaTypeRoot{
				"foo": &MediaTypeDefinition{Identifier: "application/json"},
			},
			set: []string{"foo"},
		},
		{
			name: "two",
			root: MediaTypeRoot{
				"foo": &MediaTypeDefinition{Identifier: "application/json"},
				"bar": &MediaTypeDefinition{Identifier: "application/xml"},
			},
			set: []string{"foo", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.root
			Design.MediaTypes = make(map[string]*MediaTypeDefinition)
			var sets []dslengine.DefinitionSet
			it := func(s dslengine.DefinitionSet) error {
				sets = append(sets, s)
				return nil
			}
			root.IterateSets(it)
			if len(sets) != 1 {
				t.Errorf("expected 1, got %d", len(sets))
			}
			if len(sets[0]) != len(tt.set) {
				t.Errorf("expected %d, got %d", len(root), len(sets[0]))
			}
			for i, set := range sets[0] {
				if set != root[tt.set[i]] {
					t.Errorf("expected %v, got %v", root[tt.set[i]], set)
				}
			}
		})
	}
}
