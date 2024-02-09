package apidsl

import (
	"testing"

	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestAPI(t *testing.T) {
	dslengine.Reset()
	API("test", func() {})
	dslengine.Run()

	if design.Design.Name != "test" {
		t.Errorf("expected API name to be 'test', got %s", design.Design.Name)
	}
}
