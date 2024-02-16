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
	roots, err := SortRoots()
	if err != nil {
		return err
	}
	if len(roots) == 0 {
		return nil
	}

	// Execute the DSL
	Errors = nil
	executed := 0
	recursed := 0
	for executed < len(roots) {
		recursed++
		start := executed
		executed = len(roots)
		for _, root := range roots[start:] {
			root.IterateSets(runSet)
		}
		if recursed > 100 {
			// Let's cross that bridge once we get there
			return fmt.Errorf("too many generated roots, infinite loop?")
		}
	}
	if len(Errors) > 0 {
		return Errors
	}

	// Validate and finalize
	for _, root := range roots {
		root.IterateSets(validateSet)
	}
	if len(Errors) > 0 {
		return Errors
	}

	for _, root := range roots {
		root.IterateSets(finalizeSet)
	}

	return nil
}

// Execute runs the given DSL to initialize the given definition. It returns true on success.
// It returns false and appends to Errors on failure.
// Note that `Run` takes care of calling `Execute` on all definitions that implement Source.
// This function is intended for use by definitions that run the DSL at declaration time rather than
// store the DSL for execution by the dsl engine (usually simple independent definitions).
// The DSL should use ReportError to record DSL execution errors.
func Execute(dsl func(), def Definition) bool {
	if dsl == nil {
		return true
	}
	initCount := len(Errors)
	ctxStack = append(ctxStack, def)
	dsl()
	ctxStack = ctxStack[:len(ctxStack)-1]
	return len(Errors) <= initCount
}

// CurrentDefinition returns the definition whose initialization DSL is currently being executed.
func CurrentDefinition() Definition {
	current := ctxStack.Current()
	if current == nil {
		return &TopLevelDefinition{}
	}
	return current
}

// IsTopLevelDefinition returns true if the currently evaluated DSL is a root
// DSL (i.e. is not being run in the context of another definition).
func IsTopLevelDefinition() bool {
	_, ok := CurrentDefinition().(*TopLevelDefinition)
	return ok
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

var _ Definition = &TopLevelDefinition{}

// TopLevelDefinition represents the top-level file definitions, done
// with `var _ = `.  An instance of this object is returned by
// `CurrentDefinition()` when at the top-level.
type TopLevelDefinition struct{}

// Context tells the DSL engine which context we're in when showing
// errors.
func (t *TopLevelDefinition) Context() string { return "top-level" }

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

// runSet executes the DSL for all definitions in the given set. The definition DSLs may append to
// the set as they execute.
func runSet(set DefinitionSet) error {
	executed := 0
	recursed := 0
	for executed < len(set) {
		recursed++
		for _, def := range set[executed:] {
			executed++
			if source, ok := def.(Source); ok {
				if dsl := source.DSL(); dsl != nil {
					Execute(dsl, source)
				}
			}
		}
		if recursed > 100 {
			return fmt.Errorf("too many generated definitions, infinite loop?")
		}
	}
	return nil
}

// validateSet runs the validation on all the set definitions that define one.
func validateSet(set DefinitionSet) error {
	errors := &ValidationErrors{}
	for _, def := range set {
		if validate, ok := def.(Validate); ok {
			if err := validate.Validate(); err != nil {
				errors.AddError(def, err)
			}
		}
	}
	err := errors.AsError()
	if err != nil {
		Errors = append(Errors, &Error{GoError: err})
	}
	return err
}

// finalizeSet runs the validation on all the set definitions that define one.
func finalizeSet(set DefinitionSet) error {
	for _, def := range set {
		if finalize, ok := def.(Finalize); ok {
			finalize.Finalize()
		}
	}
	return nil
}

// SortRoots orders the DSL roots making sure dependencies are last. It returns an error if there
// is a dependency cycle.
func SortRoots() ([]Root, error) {
	if len(roots) == 0 {
		return nil, nil
	}

	// First flatten dependencies for each root
	rootDeps := make(map[string][]Root, len(roots))
	rootByName := make(map[string]Root, len(roots))
	for _, r := range roots {
		sorted := sortDependencies(r, func(r Root) []Root { return r.DependsOn() })
		length := len(sorted)
		for i := 0; i < length/2; i++ {
			sorted[i], sorted[length-i-1] = sorted[length-i-1], sorted[i]
		}
		rootDeps[r.DSLName()] = sorted
		rootByName[r.DSLName()] = r
	}
	// Now check for cycles
	for name, deps := range rootDeps {
		root := rootByName[name]
		for otherName, otherdeps := range rootDeps {
			other := rootByName[otherName]
			if root.DSLName() == other.DSLName() {
				continue
			}
			dependsOnOther := false
			for _, dep := range deps {
				if dep.DSLName() == other.DSLName() {
					dependsOnOther = true
					break
				}
			}
			if dependsOnOther {
				for _, dep := range otherdeps {
					if dep.DSLName() == root.DSLName() {
						return nil, fmt.Errorf("dependency cycle: %s and %s depend on each other (directly or not)",
							root.DSLName(), other.DSLName())
					}
				}
			}
		}
	}

	// Now sort top level DSLs
	var sorted []Root
	for _, r := range roots {
		s := sortDependencies(r, func(r Root) []Root { return rootDeps[r.DSLName()] })
		for _, s := range s {
			found := false
			for _, r := range sorted {
				if r.DSLName() == s.DSLName() {
					found = true
					break
				}
			}
			if !found {
				sorted = append(sorted, s)
			}
		}
	}
	return sorted, nil
}

// sortDependencies sorts the dependencies of the given root in the given slice.
func sortDependencies(root Root, depFunc func(Root) []Root) []Root {
	seen := make(map[string]bool, len(roots))
	var sorted []Root
	sortDependenciesR(root, seen, &sorted, depFunc)
	return sorted
}

// sortDependenciesR sorts the dependencies of the given root in the given slice.
func sortDependenciesR(root Root, seen map[string]bool, sorted *[]Root, depFunc func(Root) []Root) {
	for _, dep := range depFunc(root) {
		if !seen[dep.DSLName()] {
			seen[root.DSLName()] = true
			sortDependenciesR(dep, seen, sorted, depFunc)
		}
	}
	*sorted = append(*sorted, root)
}
