package shogoa

import "testing"

func TestService(t *testing.T) {
	s := New("test")
	if s.Name != "test" {
		t.Errorf("expected service name to be 'test', got %s", s.Name)
	}
}
