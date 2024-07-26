package shogoa

import (
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	t.Run("shogoa Middleware", func(t *testing.T) {
		called := false
		myMiddleware := Middleware(func(h Handler) Handler {
			called = true
			return h
		})
		got, err := NewMiddleware(myMiddleware)
		if err != nil {
			t.Fatal(err)
		}
		got(nil)
		if !called {
			t.Fatal("middleware not called")
		}
	})

	t.Run("shogoa middleware func", func(t *testing.T) {
		called := false
		myMiddleware := func(h Handler) Handler {
			called = true
			return h
		}
		got, err := NewMiddleware(myMiddleware)
		if err != nil {
			t.Fatal(err)
		}
		got(nil)
		if !called {
			t.Fatal("middleware not called")
		}
	})

	// TODO add more tests
}
