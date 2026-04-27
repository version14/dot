package generator

import "fmt"

// ErrCircularDep is returned when DependsOn forms a cycle.
type ErrCircularDep struct {
	Cycle []string
}

func (e *ErrCircularDep) Error() string {
	return fmt.Sprintf("generator: circular dependency: %v", e.Cycle)
}

// ErrConflict is returned when two scheduled generators declare each other in ConflictsWith.
type ErrConflict struct {
	A, B string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("generator: conflict between %q and %q", e.A, e.B)
}

// ErrMissingDep is returned when a generator depends on one not in the resolved set.
type ErrMissingDep struct {
	Generator, Missing string
}

func (e *ErrMissingDep) Error() string {
	return fmt.Sprintf("generator: %q requires %q which is not registered", e.Generator, e.Missing)
}

// ErrUnknownGenerator is returned when a generator name has no registered Manifest.
type ErrUnknownGenerator struct{ Name string }

func (e *ErrUnknownGenerator) Error() string {
	return fmt.Sprintf("generator: unknown generator %q", e.Name)
}

// ErrGeneratorFailed wraps an error returned from a generator's Generate function.
type ErrGeneratorFailed struct {
	Name string
	Err  error
}

func (e *ErrGeneratorFailed) Error() string {
	return fmt.Sprintf("generator: %q failed: %v", e.Name, e.Err)
}

func (e *ErrGeneratorFailed) Unwrap() error { return e.Err }

// ErrValidationFailed is returned when a structural validator fails on a generated state.
type ErrValidationFailed struct {
	Generator, Validator, Reason string
}

func (e *ErrValidationFailed) Error() string {
	return fmt.Sprintf("generator: %q validator %q failed: %s", e.Generator, e.Validator, e.Reason)
}
