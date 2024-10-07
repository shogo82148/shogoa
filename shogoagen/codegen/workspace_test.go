package codegen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shogo82148/shogoa/shogoagen/codegen"
)

func abs(elems ...string) string {
	r, err := filepath.Abs(filepath.Join(append([]string{""}, elems...)...))
	if err != nil {
		panic("abs: " + err.Error())
	}
	return r
}

func TestWorkspaceFor(t *testing.T) {
	t.Run("with GOMOD, with GO111MODULE=auto, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode workspace
		workspace, err := codegen.WorkspaceFor(abs(dir, "xx", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if workspace.Path != gopath {
			t.Errorf("unexpected workspace.Path: %s, want %s", workspace.Path, gopath)
		}
	})

	t.Run("with GOMOD, with GO111MODULE=auto, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a Module mode workspace
		workspace, err := codegen.WorkspaceFor(dir)
		if err != nil {
			t.Fatal(err)
		}
		if workspace.Path != dir {
			t.Errorf("unexpected workspace.Path: %s, want %s", workspace.Path, dir)
		}
	})

	t.Run("with GOMOD, with GO111MODULE=on, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return a Module mode workspace
		workspace, err := codegen.WorkspaceFor(abs(dir, "xx", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if workspace.Path != dir {
			t.Errorf("unexpected workspace.Path: %s, want %s", workspace.Path, dir)
		}
	})

	t.Run("with GOMOD, with GO111MODULE=on, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return a Module mode workspace
		workspace, err := codegen.WorkspaceFor(dir)
		if err != nil {
			t.Fatal(err)
		}
		if workspace.Path != dir {
			t.Errorf("unexpected workspace.Path: %s, want %s", workspace.Path, dir)
		}
	})

	t.Run("with GOMOD, with GO111MODULE=off, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return a Module mode workspace
		workspace, err := codegen.WorkspaceFor(abs(dir, "xx", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if workspace.Path != gopath {
			t.Errorf("unexpected workspace.Path: %s, want %s", workspace.Path, gopath)
		}
	})

	t.Run("with GOMOD, with GO111MODULE=off, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return a Module mode workspace
		_, err := codegen.WorkspaceFor(dir)
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=auto, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode workspace
		workspace, err := codegen.WorkspaceFor(abs(dir, "xx", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if workspace.Path != gopath {
			t.Errorf("unexpected workspace.Path: %s, want %s", workspace.Path, gopath)
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=auto, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return an error
		_, err := codegen.WorkspaceFor(dir)
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=on, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return an error
		_, err := codegen.WorkspaceFor(abs(dir, "xx", "bar", "xx", "42"))
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=on, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return an error
		_, err := codegen.WorkspaceFor(dir)
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})
}

func TestPackageFor(t *testing.T) {
	t.Run("with GOMOD, with GO111MODULE=auto, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "bar/xx")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=auto, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "bar/xx")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=on, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "xx/src/bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "xx/src/bar/xx")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=on, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "bar/xx")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=off, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "bar/xx")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=off, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return an error
		_, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=auto, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "bar/xx")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=auto, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return an error
		_, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=on, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return an error
		_, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=on, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return an error
		_, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=off, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return a GOPATH mode package
		pkg, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg.Path != "bar/xx" {
			t.Errorf("unexpected pkg.Path: %s, want %s", pkg.Path, "bar/xx")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=off, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return an error
		_, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Error("expected an error, but no error")
		}
	})
}

func TestPackage_Abs(t *testing.T) {
	// with GO111MODULE=auto
	dir := t.TempDir()
	gopath := abs(dir, "xx")
	t.Setenv("GOPATH", gopath)
	t.Setenv("GO111MODULE", "auto")

	// with go.mod
	if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("inside GOPATH", func(t *testing.T) {
		pkg, err := codegen.PackageFor(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}

		// should return the absolute path to the GOPATH directory
		got := pkg.Abs()
		want := abs(gopath, "src", "bar", "xx")
		if got != want {
			t.Errorf("unexpected pkg.Abs(): %s, want %s", got, want)
		}
	})

	t.Run("outside GOPATH", func(t *testing.T) {
		pkg, err := codegen.PackageFor(abs(dir, "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}

		// should return the absolute path to the Module directory
		got := pkg.Abs()
		want := abs(dir, "bar", "xx")
		if got != want {
			t.Errorf("unexpected pkg.Abs(): %s, want %s", got, want)
		}
	})
}

func TestPackagePath(t *testing.T) {
	t.Run("with GOMOD, with GO111MODULE=auto, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode package
		pkg, err := codegen.PackagePath(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "bar/xx/42")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=auto, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a Module mode package path
		pkg, err := codegen.PackagePath(abs(dir, "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "bar/xx/42")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=on, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return a Module mode package path
		pkg, err := codegen.PackagePath(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "xx/src/bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "xx/src/bar/xx/42")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=on, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return a Module mode package path
		pkg, err := codegen.PackagePath(abs(dir, "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "bar/xx/42")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=off, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return a GOPATH mode package path
		pkg, err := codegen.PackagePath(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "bar/xx/42")
		}
	})

	t.Run("with GOMOD, with GO111MODULE=off, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with go.mod
		if err := os.WriteFile(abs(dir, "go.mod"), []byte{}, 0o644); err != nil {
			t.Fatal(err)
		}

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return an error
		_, err := codegen.PackagePath(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=auto, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return a GOPATH mode package
		pkg, err := codegen.PackagePath(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "bar/xx/42")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=auto, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=auto
		t.Setenv("GO111MODULE", "auto")

		// should return an error
		_, err := codegen.PackagePath(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=on, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return an error
		_, err := codegen.PackagePath(abs(gopath, "src", "bar", "xx", "42"))
		if err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=on, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=on
		t.Setenv("GO111MODULE", "on")

		// should return an error
		_, err := codegen.PackagePath(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Errorf("expected an error, but no error")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=off, inside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return a GOPATH mode package
		pkg, err := codegen.PackagePath(abs(gopath, "src", "bar", "xx", "42"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg != "bar/xx/42" {
			t.Errorf("unexpected pkg: %s, want %s", pkg, "bar/xx/42")
		}
	})

	t.Run("with no GOMOD, with GO111MODULE=off, outside GOPATH", func(t *testing.T) {
		dir := t.TempDir()
		gopath := abs(dir, "xx")
		t.Setenv("GOPATH", gopath)

		// with GO111MODULE=off
		t.Setenv("GO111MODULE", "off")

		// should return an error
		_, err := codegen.PackagePath(abs(dir, "bar", "xx", "42"))
		if err == nil {
			t.Errorf("expected an error, but no error")
		}
	})
}
