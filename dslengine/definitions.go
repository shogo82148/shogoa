package dslengine

// Definition is the common interface implemented by all definitions.
type Definition interface {
	// Context is used to build error messages that refer to the definition.
	Context() string
}

// DefinitionSet contains DSL definitions that are executed as one unit.
// The slice elements may implement the Validate an, Source interfaces to enable the
// corresponding behaviors during DSL execution.
type DefinitionSet []Definition

type Root interface {
	// DSLName is displayed by the runner upon executing the DSL.
	// Registered DSL roots must have unique names.
	DSLName() string
	// DependsOn returns the list of other DSL roots this root depends on.
	// The DSL engine uses this function to order execution of the DSLs.
	DependsOn() []Root
	// IterateSets implements the visitor pattern: is is called by the engine so the
	// DSL can control the order of execution. IterateSets calls back the engine via
	// the given iterator as many times as needed providing the DSL definitions that
	// must be run for each callback.
	IterateSets(SetIterator)
	// Reset restores the root to pre DSL execution state.
	// This is mainly used by tests.
	Reset()
}

// SetIterator is the function signature used to iterate over definition sets with
// IterateSets.
type SetIterator func(s DefinitionSet) error
