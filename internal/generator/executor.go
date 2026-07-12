package generator

import (
	"fmt"

	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

// Invocation is one scheduled call to a generator with the loop scope it
// runs under. Resolver produces these; Executor runs them in order.
type Invocation struct {
	Name      string
	LoopStack []flow.LoopFrame
}

// Observer is notified around every generator invocation, with the root state as
// it stands at that moment.
//
// It exists so a caller can see what each generator *contributed*, which is not
// recoverable afterwards: generators merge into shared files like package.json,
// so the finished state records only that a generator touched a file, never what
// it put there. Bracketing each call is the only way to attribute the delta.
//
// Fingerprinting is the one user. See internal/fingerprint.
type Observer interface {
	BeforeGenerate(name string, root *state.VirtualProjectState)
	AfterGenerate(name string, root *state.VirtualProjectState)
}

// Executor runs ordered Invocations against a VirtualProjectState, building
// scoped Answers via flow.FlattenScope before each call.
type Executor struct {
	Registry *Registry
	Logger   dotapi.Logger

	// Observer, when set, is called around each generator. Optional.
	Observer Observer
}

func NewExecutor(reg *Registry, logger dotapi.Logger) *Executor {
	if logger == nil {
		logger = dotapi.DiscardLogger{}
	}
	return &Executor{Registry: reg, Logger: logger}
}

// Execute invokes generators in the given order. The provided ProjectSpec is
// passed to each generator; State accumulates writes across all invocations.
// On the first generator error, execution halts and the error is returned.
func (e *Executor) Execute(invocations []Invocation, s *spec.ProjectSpec, vstate *state.VirtualProjectState) error {
	if e.Registry == nil {
		return fmt.Errorf("generator: executor has no registry")
	}
	previous := []string{}
	for _, inv := range invocations {
		entry, ok := e.Registry.Get(inv.Name)
		if !ok {
			return &ErrUnknownGenerator{Name: inv.Name}
		}

		scoped := flow.FlattenScope(s.Answers, inv.LoopStack)

		// Scope file writes to apps/<name>/ when running inside an app loop.
		stateForInv := vstate
		if len(inv.LoopStack) > 0 {
			if appName, ok := scoped["app-name"].(string); ok && appName != "" {
				stateForInv = vstate.WithPrefix("apps/" + appName)
			}
		}

		ctx := &dotapi.Context{
			Spec:         s,
			Answers:      scoped,
			State:        stateForInv,
			PreviousGens: append([]string(nil), previous...),
			Logger:       e.Logger,
		}

		stateForInv.SetCurrentGenerator(inv.Name)
		if e.Observer != nil {
			e.Observer.BeforeGenerate(inv.Name, vstate)
		}
		if err := entry.Generator.Generate(ctx); err != nil {
			return &ErrGeneratorFailed{Name: inv.Name, Err: err}
		}
		if e.Observer != nil {
			e.Observer.AfterGenerate(inv.Name, vstate)
		}
		previous = append(previous, inv.Name)
	}
	return nil
}
