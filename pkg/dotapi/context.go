package dotapi

import (
	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/state"
)

// Context is the per-invocation handle every generator receives.
//
//   - Spec is the full ProjectSpec (read-only conceptually).
//   - Answers is the FlattenScope view for this invocation: globals overlaid
//     with each enclosing loop frame, deepest wins.
//   - State is the in-memory project filesystem; all writes go here.
//   - PreviousGens is the ordered list of generator names already executed.
type Context struct {
	Spec         *spec.ProjectSpec
	Answers      map[string]interface{}
	State        *state.VirtualProjectState
	PreviousGens []string
	Logger       Logger
}

// Logger is the minimal logging surface generators use. The CLI provides a
// Lipgloss-backed implementation; tests can pass a discard logger.
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// DiscardLogger drops all log output. Useful in tests.
type DiscardLogger struct{}

func (DiscardLogger) Infof(string, ...interface{})  {}
func (DiscardLogger) Warnf(string, ...interface{})  {}
func (DiscardLogger) Errorf(string, ...interface{}) {}
