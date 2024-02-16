package dslengine

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Error represents an error that occurred while running the API DSL.
// It contains the name of the file and line number of where the error
// occurred as well as the original Go error.
type Error struct {
	GoError error
	File    string
	Line    int
}

// MultiError collects all DSL errors. It implements error.
type MultiError []*Error

// DSL evaluation contexts stack
type contextStack []Definition

// Errors contains the DSL execution errors if any.
var Errors MultiError

// Global DSL evaluation stack
var ctxStack contextStack

// Registered DSL roots
var roots []Root

// DSL package paths used to compute error locations (skip the frames in these packages)
var dslPackages = map[string]bool{
	"github.com/shogo82148/shogoa/": true,
}

// Reset uses the registered RootFuncs to re-initialize the DSL roots.
// This is useful to tests.
func Reset() {
	for _, r := range roots {
		r.Reset()
	}
	Errors = nil
}

// Run runs the given root definitions. It iterates over the definition sets
// multiple times to first execute the DSL, the validate the resulting
// definitions and finally finalize them. The executed DSL may register new
// roots to have them be executed (last) in the same run.
func Run() error {
	if len(roots) == 0 {
		return nil
	}
	return nil
}

// Error returns the error message.
func (m MultiError) Error() string {
	msgs := make([]string, len(m))
	for i, de := range m {
		msgs[i] = de.Error()
	}
	return strings.Join(msgs, "\n")
}

// Error returns the underlying error message.
func (de *Error) Error() string {
	if err := de.GoError; err != nil {
		if de.File == "" {
			return err.Error()
		}
		return fmt.Sprintf("[%s:%d] %s", de.File, de.Line, err.Error())
	}
	return ""
}

// ReportError records a DSL error for reporting post DSL execution.
func ReportError(fm string, vals ...any) {
	var suffix string
	if cur := ctxStack.Current(); cur != nil {
		if ctx := cur.Context(); ctx != "" {
			suffix = fmt.Sprintf(" in %s", ctx)
		}
	} else {
		suffix = " (top level)"
	}

	vals = append(vals[:], suffix)
	err := fmt.Errorf(fm+"%s", vals...)
	file, line := computeErrorLocation()
	Errors = append(Errors, &Error{
		GoError: err,
		File:    file,
		Line:    line,
	})
}

// Current evaluation context, i.e. object being currently built by DSL
func (s contextStack) Current() Definition {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

// computeErrorLocation implements a heuristic to find the location in the user
// code where the error occurred. It walks back the callstack until the file
// doesn't match "/goa/design/*.go" or one of the DSL package paths.
// When successful it returns the file name and line number, empty string and
// 0 otherwise.
func computeErrorLocation() (file string, line int) {
	skipFunc := func(file string) bool {
		if strings.HasSuffix(file, "_test.go") { // Be nice with tests
			return false
		}
		file = filepath.ToSlash(file)
		for pkg := range dslPackages {
			if strings.Contains(file, pkg) {
				return true
			}
		}
		return false
	}
	depth := 2
	_, file, line, _ = runtime.Caller(depth)
	for skipFunc(file) {
		depth++
		_, file, line, _ = runtime.Caller(depth)
	}
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	wd, err = filepath.Abs(wd)
	if err != nil {
		return
	}
	f, err := filepath.Rel(wd, file)
	if err != nil {
		return
	}
	file = f
	return
}
