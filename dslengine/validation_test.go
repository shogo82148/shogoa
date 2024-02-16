package dslengine

import (
	"fmt"
	"testing"
)

var _ Definition = new(mockDefinition)

type mockDefinition struct{}

func (m *mockDefinition) Context() string {
	return "mockDefinition"
}

func TestValidationErrors(t *testing.T) {
	verr := &ValidationErrors{}
	if verr.Error() != "" {
		t.Errorf("expected empty error, got %s", verr.Error())
	}

	verr.Add(new(mockDefinition), "error %s", "1")
	if verr.Error() != "mockDefinition: error 1" {
		t.Errorf("unexpected error, got %s", verr.Error())
	}

	verr.Add(new(mockDefinition), "error %s", "2")
	if verr.Error() != "mockDefinition: error 1\nmockDefinition: error 2" {
		t.Errorf("unexpected error, got %s", verr.Error())
	}

	verr2 := &ValidationErrors{}
	verr2.Add(new(mockDefinition), "error %s", "3")
	verr2.Add(new(mockDefinition), "error %s", "4")
	verr.Merge(verr2)
	if verr.Error() != "mockDefinition: error 1\nmockDefinition: error 2\nmockDefinition: error 3\nmockDefinition: error 4" {
		t.Errorf("unexpected error, got %s", verr.Error())
	}

	verr2 = &ValidationErrors{}
	verr2.Add(new(mockDefinition), "error %s", "5")
	verr2.Add(new(mockDefinition), "error %s", "6")
	verr2.Merge(nil)
	verr.Merge(verr2)
	if verr.Error() != "mockDefinition: error 1\nmockDefinition: error 2\nmockDefinition: error 3\nmockDefinition: error 4\nmockDefinition: error 5\nmockDefinition: error 6" {
		t.Errorf("unexpected error, got %s", verr.Error())
	}

	verr = &ValidationErrors{}
	verr.AddError(new(mockDefinition), fmt.Errorf("error %s", "7"))
	if verr.Error() != "mockDefinition: error 7" {
		t.Errorf("unexpected error, got %s", verr.Error())
	}

	verr2 = &ValidationErrors{}
	verr2.Add(new(mockDefinition), "error %s", "8")
	verr2.Add(new(mockDefinition), "error %s", "9")
	verr.AddError(new(mockDefinition), verr2)
}

func TestValidationErrorsMergeWithSelf(t *testing.T) {
	verr := &ValidationErrors{}
	verr.Add(new(mockDefinition), "error %s", "self")
	verr.Merge(verr)
	if len(verr.Errors) != 1 {
		t.Errorf("expected 1 error after merging with self, got %d", len(verr.Errors))
	}
}

func TestValidationErrorsAddNilError(t *testing.T) {
	verr := &ValidationErrors{}
	verr.AddError(new(mockDefinition), nil)
	if len(verr.Errors) != 0 {
		t.Errorf("expected 0 errors after adding nil, got %d", len(verr.Errors))
	}
}
