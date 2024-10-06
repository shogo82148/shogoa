package design

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestAttributeDefinition_Inherit(t *testing.T) {
	t.Run("does not change with a empty parent", func(t *testing.T) {
		parent := &AttributeDefinition{Type: Object{}}
		child := &AttributeDefinition{Type: Object{
			"c": &AttributeDefinition{Type: String},
		}}

		child.Inherit(parent)
		obj := child.Type.(Object)
		if _, ok := obj["c"]; !ok {
			t.Error("child does not inherit parent")
		}
	})

	t.Run("does not change with a parent that defines no inherited attributes", func(t *testing.T) {
		parent := &AttributeDefinition{Type: Object{
			"other": &AttributeDefinition{
				Type:         String,
				DefaultValue: "default",
			},
		}}
		child := &AttributeDefinition{Type: Object{
			"c": &AttributeDefinition{Type: String},
		}}

		child.Inherit(parent)
		obj := child.Type.(Object)
		if _, ok := obj["c"]; !ok {
			t.Error("child does not inherit parent")
		}
	})

	t.Run("inherit with a parent that defines an inherited attribute", func(t *testing.T) {
		parent := &AttributeDefinition{Type: Object{
			"c": &AttributeDefinition{
				Type:         String,
				DefaultValue: "default",
				Metadata:     map[string][]string{"swagger:read-only": nil},
			},
		}}
		child := &AttributeDefinition{Type: Object{
			"c": &AttributeDefinition{Type: String},
		}}

		child.Inherit(parent)
		obj := child.Type.(Object)
		if obj["c"].DefaultValue != "default" {
			t.Errorf("child does not inherit parent: DefaultValue: %v", obj["c"].DefaultValue)
		}
		if !obj["c"].IsReadOnly() {
			t.Error("child does not inherit parent: ReadOnly")
		}
	})
}

func TestAttributeDefinition_IsRequired(t *testing.T) {
	integer := &AttributeDefinition{Type: Integer}
	attribute := &AttributeDefinition{
		Type: Object{"required": integer},
		Validation: &dslengine.ValidationDefinition{
			Required: []string{"required"},
		},
	}
	if !attribute.IsRequired("required") {
		t.Error("required field is not required")
	}
	if attribute.IsRequired("non-required") {
		t.Error("non-required field is required")
	}
}

func TestActionDefinition_IterateHeaders(t *testing.T) {
	resource := &ResourceDefinition{}
	action := &ActionDefinition{
		Parent: resource,
		Headers: &AttributeDefinition{
			Type: Object{
				"a": &AttributeDefinition{Type: String},
			},
		},
	}

	names := []string{}
	// iterator tha collects header names
	it := func(name string, _ bool, _ *AttributeDefinition) error {
		names = append(names, name)
		return nil
	}
	if err := action.IterateHeaders(it); err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 {
		t.Errorf("got %d names, expected 1", len(names))
	}
	if names[0] != "a" {
		t.Errorf("got %s, expected a", names[0])
	}
}

func TestActionDefinition_Finalize(t *testing.T) {
	resource := &ResourceDefinition{
		Responses: map[string]*ResponseDefinition{
			"NotFound": {Name: "NotFound", Status: 404},
		},
	}
	action := &ActionDefinition{Parent: resource}

	action.Finalize()
	if _, ok := action.Responses["NotFound"]; !ok {
		t.Error("does not merge the resource responses")
	}
}

func TestRouteDefinition_FullPath(t *testing.T) {
	Design.Reset()

	tests := []struct {
		name               string
		actionPath         string
		resourcePath       string
		parentResourcePath string
		want               string
	}{
		{
			name:               "with relative path",
			actionPath:         "/action",
			resourcePath:       "/resource",
			parentResourcePath: "/parent",
			want:               "/parent/resource/action",
		},
		{
			name:               "an action with absolute route",
			actionPath:         "//action",
			resourcePath:       "/resource",
			parentResourcePath: "/parent",
			want:               "/action",
		},
		{
			name:               "a resource with absolute route",
			actionPath:         "/action",
			resourcePath:       "//resource",
			parentResourcePath: "/parent",
			want:               "/resource/action",
		},
		{
			name:               "with trailing slashes",
			actionPath:         "/action/",
			resourcePath:       "/resource",
			parentResourcePath: "/parent",
			want:               "/parent/resource/action/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			showAct := &ActionDefinition{}
			showRoute := &RouteDefinition{
				Path:   tt.parentResourcePath,
				Parent: showAct,
			}
			showAct.Routes = []*RouteDefinition{showRoute}
			parentResource := &ResourceDefinition{
				Name:    "foo",
				Actions: map[string]*ActionDefinition{"show": showAct},
			}
			showAct.Parent = parentResource
			Design.Resources = map[string]*ResourceDefinition{"foo": parentResource}

			action := &ActionDefinition{}
			route := &RouteDefinition{
				Path:   tt.actionPath,
				Parent: action,
			}
			action.Routes = []*RouteDefinition{route}
			resource := &ResourceDefinition{
				Actions:    map[string]*ActionDefinition{"action": action},
				BasePath:   tt.resourcePath,
				ParentName: parentResource.Name,
			}
			action.Parent = resource

			got := route.FullPath()
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestActionDefinition_AllParams(t *testing.T) {
	Design.Reset()

	// Parent resource
	parent := func() *ResourceDefinition {
		baseParams := &AttributeDefinition{
			Type: Object{
				"pbasepath":  &AttributeDefinition{Type: String},
				"pbasequery": &AttributeDefinition{Type: String},
			},
		}
		parent := &ResourceDefinition{
			Name:                "parent",
			CanonicalActionName: "canonical",
			BasePath:            "/:pbasepath",
			Params:              baseParams,
		}
		canParams := &AttributeDefinition{
			Type: Object{
				"canpath":  &AttributeDefinition{Type: String},
				"canquery": &AttributeDefinition{Type: String},
			},
		}
		canonical := &ActionDefinition{
			Name:   "canonical",
			Parent: parent,
			Params: canParams,
		}
		croute := &RouteDefinition{
			Path:   "/:canpath",
			Parent: canonical,
		}
		canonical.Routes = []*RouteDefinition{croute}
		parent.Actions = map[string]*ActionDefinition{"canonical": canonical}
		return parent
	}()

	// Resource
	resource := func() *ResourceDefinition {
		baseParams := &AttributeDefinition{
			Type: Object{
				"basepath":  &AttributeDefinition{Type: String},
				"basequery": &AttributeDefinition{Type: String},
			},
		}
		resource := &ResourceDefinition{
			Name:       "child",
			ParentName: "parent",
			BasePath:   "/:basepath",
			Params:     baseParams,
		}
		return resource
	}()

	// Action
	action := func() *ActionDefinition {
		params := &AttributeDefinition{
			Type: Object{
				"path":     &AttributeDefinition{Type: String},
				"query":    &AttributeDefinition{Type: String},
				"basepath": &AttributeDefinition{Type: String},
			},
		}
		action := &ActionDefinition{
			Name:   "action",
			Parent: resource,
			Params: params,
		}
		route := &RouteDefinition{
			Path:   "/:path",
			Parent: action,
		}
		action.Routes = []*RouteDefinition{route}
		resource.Actions = map[string]*ActionDefinition{"action": action}
		return action
	}()

	Design.Resources = map[string]*ResourceDefinition{
		"resource": resource,
		"parent":   parent,
	}
	Design.BasePath = "/:apipath"
	Design.Params = &AttributeDefinition{
		Type: Object{
			"apipath":  &AttributeDefinition{Type: String},
			"apiquery": &AttributeDefinition{Type: String},
		},
	}

	// check all params
	allParams := action.AllParams().Type.ToObject()
	keys := []string{
		"apipath",
		"apiquery",
		"basepath",
		"basequery",
		"canpath",
		"path",
		"pbasepath",
		"query",
	}
	if len(allParams) != len(keys) {
		t.Errorf("got %d keys, expected %d", len(allParams), len(keys))
	}
	for _, key := range keys {
		if _, ok := allParams[key]; !ok {
			t.Errorf("missing key %q", key)
		}
	}

	// check the path params
	pathParams := action.PathParams().Type.ToObject()
	keys = []string{
		"path", "basepath", "canpath", "pbasepath", "apipath",
	}
	if len(pathParams) != len(keys) {
		t.Errorf("got %d keys, expected %d", len(pathParams), len(keys))
	}
	for _, key := range keys {
		if _, ok := pathParams[key]; !ok {
			t.Errorf("missing key %q", key)
		}
	}
}

func TestResourceDefinition_PathParams(t *testing.T) {
	t.Run("Given a resource with a nil base params", func(t *testing.T) {
		Design.Reset()

		resource := &ResourceDefinition{
			Name:     "resource",
			BasePath: "/:basepath",
		}
		Design.Resources = map[string]*ResourceDefinition{"resource": resource}

		pathParams := resource.PathParams().Type.ToObject()
		if len(pathParams) != 0 {
			t.Errorf("got %d keys, expected 0", len(pathParams))
		}
	})

	t.Run("Given a resource defining a subset of all base path params", func(t *testing.T) {
		Design.Reset()

		params := Object{
			"basepath": &AttributeDefinition{Type: String},
		}
		resource := &ResourceDefinition{
			Name:     "resource",
			BasePath: "/:basepath/:sub",
			Params:   &AttributeDefinition{Type: params},
		}
		Design.Resources = map[string]*ResourceDefinition{"resource": resource}

		pathParams := resource.PathParams().Type.ToObject()
		if len(pathParams) != 1 {
			t.Errorf("got %d keys, expected 1", len(pathParams))
		}
		if _, ok := pathParams["basepath"]; !ok {
			t.Error("missing key basepath")
		}
	})
}

func TestResourceDefinition_AllFileServers(t *testing.T) {
	resource := &ResourceDefinition{
		FileServers: []*FileServerDefinition{
			{FilePath: "/C"},
			{FilePath: "/B"},
			{FilePath: "/A"},
		},
	}

	got := []string{}
	for fs := range resource.AllFileServers() {
		got = append(got, fs.FilePath)
	}
	want := []string{"/A", "/B", "/C"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected file servers (-want, +got):\n%s", diff)
	}
}

func TestAPIDefinition_AllMediaTypes(t *testing.T) {
	api := &APIDefinition{
		MediaTypes: map[string]*MediaTypeDefinition{
			"application/example":  {Identifier: "application/example"},
			"application/example2": {Identifier: "application/example2"},
		},
	}

	got := []string{}
	for mt := range api.AllMediaTypes() {
		got = append(got, mt.Identifier)
	}
	want := []string{"application/example", "application/example2"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected media types (-want, +got):\n%s", diff)
	}
}

func TestAPIDefinition_AllUserTypes(t *testing.T) {
	api := &APIDefinition{
		Types: map[string]*UserTypeDefinition{
			"example":  {TypeName: "example"},
			"example2": {TypeName: "example2"},
		},
	}

	got := []string{}
	for ut := range api.AllUserTypes() {
		got = append(got, ut.TypeName)
	}
	want := []string{"example", "example2"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected user types (-want, +got):\n%s", diff)
	}
}

func TestAPIDefinition_AllResponses(t *testing.T) {
	api := &APIDefinition{
		Responses: map[string]*ResponseDefinition{
			"example":  {Name: "example"},
			"example2": {Name: "example2"},
		},
	}

	got := []string{}
	for res := range api.AllResponses() {
		got = append(got, res.Name)
	}
	want := []string{"example", "example2"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected responses (-want, +got):\n%s", diff)
	}
}

func TestAPIDefinition_AllResources(t *testing.T) {
	t.Run("sort by name", func(t *testing.T) {
		api := &APIDefinition{
			Resources: map[string]*ResourceDefinition{
				"A": {Name: "A"},
				"B": {Name: "B"},
				"C": {Name: "C"},
			},
		}

		got := []string{}
		for res := range api.AllResources() {
			got = append(got, res.Name)
		}
		want := []string{"A", "B", "C"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected resources (-want, +got):\n%s", diff)
		}
	})

	t.Run("parent is first", func(t *testing.T) {
		api := &APIDefinition{
			Resources: map[string]*ResourceDefinition{
				"A": {Name: "A", ParentName: "B"},
				"B": {Name: "B", ParentName: "C"},
				"C": {Name: "C"},
			},
		}

		got := []string{}
		for res := range api.AllResources() {
			got = append(got, res.Name)
		}
		want := []string{"C", "B", "A"}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("unexpected resources (-want, +got):\n%s", diff)
		}
	})
}

func TestAPIDefinition_AllSets(t *testing.T) {
	t.Run("should order nested resources", func(t *testing.T) {
		api := &APIDefinition{
			Name: "Test",
			Resources: map[string]*ResourceDefinition{
				"V": {Name: "V", ParentName: "W"},
				"W": {Name: "W", ParentName: "X"},
				"X": {Name: "X", ParentName: "Y"},
				"Y": {Name: "Y", ParentName: "Z"},
				"Z": {Name: "Z"},
			},
		}

		var resources []*ResourceDefinition
		for s := range api.AllSets() {
			if len(s) == 0 {
				continue
			}
			if _, ok := s[0].(*ResourceDefinition); !ok {
				continue
			}
			resources = make([]*ResourceDefinition, len(s))
			for i, res := range s {
				resources[i] = res.(*ResourceDefinition)
			}
		}

		if resources[0].Name != "Z" {
			t.Errorf("got %s, expected Z", resources[0].Name)
		}
		if resources[1].Name != "Y" {
			t.Errorf("got %s, expected Y", resources[1].Name)
		}
		if resources[2].Name != "X" {
			t.Errorf("got %s, expected X", resources[2].Name)
		}
		if resources[3].Name != "W" {
			t.Errorf("got %s, expected W", resources[3].Name)
		}
		if resources[4].Name != "V" {
			t.Errorf("got %s, expected V", resources[4].Name)
		}
	})

	t.Run("should order multiple nested resources", func(t *testing.T) {
		api := &APIDefinition{
			Name: "Test",
			Resources: map[string]*ResourceDefinition{
				"A": {Name: "A"},
				"B": {Name: "B", ParentName: "A"},
				"C": {Name: "C", ParentName: "A"},
				"I": {Name: "I"},
				"J": {Name: "J", ParentName: "K"},
				"K": {Name: "K", ParentName: "I"},
				"X": {Name: "X"},
				"Y": {Name: "Y"},
				"Z": {Name: "Z"},
			},
		}

		var resources []*ResourceDefinition
		for s := range api.AllSets() {
			if len(s) == 0 {
				continue
			}
			if _, ok := s[0].(*ResourceDefinition); !ok {
				continue
			}
			resources = make([]*ResourceDefinition, len(s))
			for i, res := range s {
				resources[i] = res.(*ResourceDefinition)
			}
		}

		if resources[0].Name != "A" {
			t.Errorf("got %s, expected A", resources[0].Name)
		}
		if resources[1].Name != "B" {
			t.Errorf("got %s, expected B", resources[1].Name)
		}
		if resources[2].Name != "C" {
			t.Errorf("got %s, expected C", resources[2].Name)
		}
		if resources[3].Name != "I" {
			t.Errorf("got %s, expected I", resources[3].Name)
		}
		if resources[4].Name != "K" {
			t.Errorf("got %s, expected K", resources[4].Name)
		}
		if resources[5].Name != "J" {
			t.Errorf("got %s, expected J", resources[5].Name)
		}
		if resources[6].Name != "X" {
			t.Errorf("got %s, expected X", resources[6].Name)
		}
		if resources[7].Name != "Y" {
			t.Errorf("got %s, expected Y", resources[7].Name)
		}
	})
}

func TestAPIDefinition_IterateSets(t *testing.T) {
	// a function that collects resource definitions for validation
	valFunc := func(validate func([]*ResourceDefinition)) func(dslengine.DefinitionSet) error {
		return func(s dslengine.DefinitionSet) error {
			if len(s) == 0 {
				return nil
			}

			if _, ok := s[0].(*ResourceDefinition); !ok {
				return nil
			}

			resources := make([]*ResourceDefinition, len(s))
			for i, res := range s {
				resources[i] = res.(*ResourceDefinition)
			}

			validate(resources)

			return nil
		}
	}
	t.Run("should order nested resources", func(t *testing.T) {
		var inspected bool
		validate := func(s []*ResourceDefinition) {
			if s[0].Name != "Z" {
				t.Errorf("got %s, expected Z", s[0].Name)
			}
			if s[1].Name != "Y" {
				t.Errorf("got %s, expected Y", s[1].Name)
			}
			if s[2].Name != "X" {
				t.Errorf("got %s, expected X", s[2].Name)
			}
			if s[3].Name != "W" {
				t.Errorf("got %s, expected W", s[3].Name)
			}
			if s[4].Name != "V" {
				t.Errorf("got %s, expected V", s[4].Name)
			}
			inspected = true
		}

		api := &APIDefinition{
			Name: "Test",
			Resources: map[string]*ResourceDefinition{
				"V": {Name: "V", ParentName: "W"},
				"W": {Name: "W", ParentName: "X"},
				"X": {Name: "X", ParentName: "Y"},
				"Y": {Name: "Y", ParentName: "Z"},
				"Z": {Name: "Z"},
			},
		}
		api.IterateSets(valFunc(validate))
		if !inspected {
			t.Error("did not iterate over the resources")
		}
	})

	t.Run("should order multiple nested resources", func(t *testing.T) {
		var inspected bool
		validate := func(s []*ResourceDefinition) {
			if s[0].Name != "A" {
				t.Errorf("got %s, expected A", s[0].Name)
			}
			if s[1].Name != "B" {
				t.Errorf("got %s, expected B", s[1].Name)
			}
			if s[2].Name != "C" {
				t.Errorf("got %s, expected C", s[2].Name)
			}
			if s[3].Name != "I" {
				t.Errorf("got %s, expected I", s[3].Name)
			}
			if s[4].Name != "K" {
				t.Errorf("got %s, expected K", s[4].Name)
			}
			if s[5].Name != "J" {
				t.Errorf("got %s, expected J", s[5].Name)
			}
			if s[6].Name != "X" {
				t.Errorf("got %s, expected X", s[6].Name)
			}
			if s[7].Name != "Y" {
				t.Errorf("got %s, expected Y", s[7].Name)
			}
			if s[8].Name != "Z" {
				t.Errorf("got %s, expected Z", s[8].Name)
			}
			inspected = true
		}

		api := &APIDefinition{
			Name: "Test",
			Resources: map[string]*ResourceDefinition{
				"A": {Name: "A"},
				"B": {Name: "B", ParentName: "A"},
				"C": {Name: "C", ParentName: "A"},
				"I": {Name: "I"},
				"J": {Name: "J", ParentName: "K"},
				"K": {Name: "K", ParentName: "I"},
				"X": {Name: "X"},
				"Y": {Name: "Y"},
				"Z": {Name: "Z"},
			},
		}
		api.IterateSets(valFunc(validate))
		if !inspected {
			t.Error("did not iterate over the resources")
		}
	})
}
