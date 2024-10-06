package randid

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	re := regexp.MustCompile(`^[a-zA-Z0-9]{24}$`)
	for range 100 {
		id := New(24)
		if !re.MatchString(id) {
			t.Errorf("invalid id: %s", id)
		}
	}
}

func BenchmarkNew(b *testing.B) {
	for range b.N {
		New(24)
	}
}
