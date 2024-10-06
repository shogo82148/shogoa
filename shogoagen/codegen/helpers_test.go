package codegen_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

func TestKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// should change uppercase letters to lowercase letters
		{
			input: "test-B",
			want:  "test-b",
		},
		{
			input: "teste",
			want:  "teste",
		},

		// should not add a dash before an abbreviation or acronym
		{
			input: "testABC",
			want:  "testabc",
		},

		// should add a dash before a title
		{
			input: "testAa",
			want:  "test-aa",
		},
		{
			input: "testAbc",
			want:  "test-abc",
		},

		// should replace underscores to dashes
		{
			input: "test_cA",
			want:  "test-ca",
		},
		{
			input: "test_D",
			want:  "test-d",
		},
	}

	for _, tt := range tests {
		if got := codegen.KebabCase(tt.input); got != tt.want {
			t.Errorf("KebabCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCommandLine(t *testing.T) {
	t.Run("with exported GOPATH", func(t *testing.T) {
		t.Setenv("GOPATH", "/xx")
		oldArgs := os.Args
		t.Cleanup(func() {
			os.Args = oldArgs
		})

		tests := []struct {
			name string
			args []string
			want string
		}{
			{
				name: "should not touch free arguments",
				args: []string{"foo", "/xx/bar/xx/42"},
				want: "$ foo /xx/bar/xx/42",
			},
			{
				name: "should replace GOPATH one match only in a long option",
				args: []string{"foo", "--opt=/xx/bar/xx/42"},
				want: "$ foo\n\t--opt=$(GOPATH)/bar/xx/42",
			},
			{
				name: "should not replace GOPATH if a match is not at the beginning of a long option",
				args: []string{"foo", "--opt=/bar/xx/42"},
				want: "$ foo\n\t--opt=/bar/xx/42",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				os.Args = tt.args
				if got := codegen.CommandLine(); got != tt.want {
					t.Errorf("CommandLine() = %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("with default GOPATH", func(t *testing.T) {
		t.Setenv("GOPATH", defaultGOPATH()) // Simulate a situation with no GOPATH exported.
		oldArgs := os.Args
		t.Cleanup(func() {
			os.Args = oldArgs
		})

		tests := []struct {
			name string
			args []string
			want string
		}{
			{
				name: "should not touch free arguments",
				args: []string{"foo", "/xx/bar/xx/42"},
				want: "$ foo /xx/bar/xx/42",
			},
			{
				name: "should replace GOPATH one match only in a long option",
				args: []string{"foo", "--opt=" + defaultGOPATH() + "/bar/xx/42"},
				want: "$ foo\n\t--opt=$(GOPATH)/bar/xx/42",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				os.Args = tt.args
				if got := codegen.CommandLine(); got != tt.want {
					t.Errorf("CommandLine() = %q, want %q", got, tt.want)
				}
			})
		}
	})
}

// Copied from go/build/build.go
func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}
