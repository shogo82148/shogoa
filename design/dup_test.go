package design_test

import (
	"reflect"
	"testing"

	"github.com/shogo82148/shogoa/design"
)

func TestDup(t *testing.T) {
	t.Run("with a primitive type", func(t *testing.T) {
		dt := design.Integer
		dup := design.Dup(dt)
		if dup != dt {
			t.Errorf("Dup(%v) = %v; want %v", dt, dup, dt)
		}
	})

	t.Run("with an array type", func(t *testing.T) {
		elemType := design.Integer
		dt := &design.Array{
			ElemType: &design.AttributeDefinition{Type: elemType},
		}
		dup := design.Dup(dt)
		if dup == dt {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
		if dup.(*design.Array).ElemType == dt.ElemType {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
	})

	t.Run("with a hash type", func(t *testing.T) {
		keyType := design.String
		elemType := design.Integer
		dt := &design.Hash{
			KeyType:  &design.AttributeDefinition{Type: keyType},
			ElemType: &design.AttributeDefinition{Type: elemType},
		}
		dup := design.Dup(dt)
		if dup == dt {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
		if dup.(*design.Hash).KeyType == dt.KeyType {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
		if dup.(*design.Hash).ElemType == dt.ElemType {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
	})

	t.Run("with a user type", func(t *testing.T) {
		typeName := "foo"
		att := &design.AttributeDefinition{Type: design.Integer}
		dt := &design.UserTypeDefinition{
			TypeName:            typeName,
			AttributeDefinition: att,
		}
		dup := design.Dup(dt)
		if dup == dt {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
		if dup.(*design.UserTypeDefinition).AttributeDefinition == att {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
	})

	t.Run("with a media type", func(t *testing.T) {
		obj := design.Object{"att": &design.AttributeDefinition{Type: design.Integer}}
		ut := &design.UserTypeDefinition{
			TypeName:            "foo",
			AttributeDefinition: &design.AttributeDefinition{Type: obj},
		}
		identifier := "vnd.application/test"
		links := map[string]*design.LinkDefinition{
			"link": {Name: "att", View: "default"},
		}
		views := map[string]*design.ViewDefinition{
			"default": {
				Name:                "default",
				AttributeDefinition: &design.AttributeDefinition{Type: obj},
			},
		}
		dt := &design.MediaTypeDefinition{
			UserTypeDefinition: ut,
			Identifier:         identifier,
			Links:              links,
			Views:              views,
		}
		dup := design.Dup(dt)
		if dup == dt {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
		if dup.(*design.MediaTypeDefinition).UserTypeDefinition == ut {
			t.Errorf("Dup(%v) = %v; want not %v", dt, dup, dt)
		}
	})

	t.Run("with a media type referring to each other", func(t *testing.T) {
		mt := &design.MediaTypeDefinition{Identifier: "application/mt1"}
		mt2 := &design.MediaTypeDefinition{Identifier: "application/mt2"}
		obj1 := design.Object{"att": &design.AttributeDefinition{Type: mt2}}
		obj2 := design.Object{"att": &design.AttributeDefinition{Type: mt}}

		att1 := &design.AttributeDefinition{Type: obj1}
		ut := &design.UserTypeDefinition{AttributeDefinition: att1}
		link1 := &design.LinkDefinition{Name: "att", View: "default"}
		view1 := &design.ViewDefinition{AttributeDefinition: att1, Name: "default"}
		mt.UserTypeDefinition = ut
		mt.Links = map[string]*design.LinkDefinition{"att": link1}
		mt.Views = map[string]*design.ViewDefinition{"default": view1}

		att2 := &design.AttributeDefinition{Type: obj2}
		ut2 := &design.UserTypeDefinition{AttributeDefinition: att2}
		link2 := &design.LinkDefinition{Name: "att", View: "default"}
		view2 := &design.ViewDefinition{AttributeDefinition: att2, Name: "default"}
		mt2.UserTypeDefinition = ut2
		mt2.Links = map[string]*design.LinkDefinition{"att": link2}
		mt2.Views = map[string]*design.ViewDefinition{"default": view2}

		dup := design.Dup(mt)
		if dup == mt {
			t.Errorf("Dup(%v) = %v; want not %v", mt, dup, mt)
		}
		if !reflect.DeepEqual(dup, mt) {
			t.Errorf("Dup(%v) = %v; want %v", mt, dup, mt)
		}
		if dup.(*design.MediaTypeDefinition).UserTypeDefinition == ut {
			t.Errorf("Dup(%v) = %v; want not %v", mt, dup, mt)
		}
	})
}

func TestDupAtt(t *testing.T) {
	t.Run("with an attribute with a type which is a media type", func(t *testing.T) {
		att := &design.AttributeDefinition{Type: &design.MediaTypeDefinition{}}
		dup := design.DupAtt(att)
		if dup == att {
			t.Errorf("DupAtt(%v) = %v; want not %v", att, dup, att)
		}
		if dup.Type != att.Type {
			t.Errorf("DupAtt(%v) = %v; want %v", att, dup, att)
		}
	})
}
