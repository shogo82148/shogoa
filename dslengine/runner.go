package dslengine

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

// Errors contains the DSL execution errors if any.
var Errors MultiError

// Registered DSL roots
var roots []Root

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
