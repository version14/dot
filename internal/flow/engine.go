package flow

import "fmt"

// FlowAdapter renders a Question to the user and returns their Answer.
// Implementations live outside the engine (e.g. internal/cli/prompt.go for Huh).
type FlowAdapter interface {
	Ask(q Question, ctx *FlowContext) (Answer, error)
}

// FlowEngine drives traversal of the flow graph, calling the adapter for each
// non-If question and routing through Next edges, fragments, and plugin
// injections (Replace / AddOption / InsertAfter).
type FlowEngine struct {
	Adapter   FlowAdapter
	Fragments *FragmentRegistry
	Hooks     *HookRegistry
}

func NewEngine(adapter FlowAdapter) *FlowEngine {
	return &FlowEngine{
		Adapter:   adapter,
		Fragments: NewFragmentRegistry(),
		Hooks:     NewHookRegistry(),
	}
}

// RegisterFragment adds a named fragment resolver to the engine.
func (e *FlowEngine) RegisterFragment(f *FlowFragment) error {
	return e.Fragments.Register(f)
}

// RegisterPlugin records a plugin so its injections can be added.
func (e *FlowEngine) RegisterPlugin(id PluginID) error {
	return e.Hooks.RegisterPlugin(id)
}

// Inject adds a plugin injection to the engine's hook registry.
func (e *FlowEngine) Inject(inj *Injection) error {
	return e.Hooks.Inject(inj)
}

// Run walks the graph from root and returns the populated FlowContext.
func (e *FlowEngine) Run(root Question) (*FlowContext, error) {
	ctx := &FlowContext{
		Answers:       map[string]AnswerNode{},
		LoopStack:     []LoopFrame{},
		VisitedNodes:  []string{},
		LoadedPlugins: pluginNames(e.Hooks.Plugins()),
	}

	state := &traversalState{}
	if err := e.traverse(root, ctx, state); err != nil {
		return nil, err
	}
	return ctx, nil
}

// traversalState tracks the per-Run continuation stack used by InsertAfter
// to return to the parent flow after an inserted chain ends.
type traversalState struct {
	continuations []*Next
}

func (s *traversalState) push(n *Next) {
	if n == nil {
		return
	}
	s.continuations = append(s.continuations, n)
}

func (s *traversalState) pop() *Next {
	n := len(s.continuations)
	if n == 0 {
		return nil
	}
	top := s.continuations[n-1]
	s.continuations = s.continuations[:n-1]
	return top
}

func (e *FlowEngine) traverse(q Question, ctx *FlowContext, state *traversalState) error {
	for q != nil {
		// Apply Replace injection first — the original target never runs if replaced.
		q = e.applyReplace(q)

		ctx.VisitedNodes = append(ctx.VisitedNodes, q.ID())

		next, err := e.step(q, ctx)
		if err != nil {
			return err
		}

		// Splice InsertAfter injections before computing the actual next.
		next = e.applyInsertAfter(q.ID(), next, state)

		q = e.resolveNext(next, ctx, state)
	}
	return nil
}

// step handles a single node: IfQuestion routes without input; OptionQuestion
// gets plugin-added options merged in; everything else is asked via the adapter.
func (e *FlowEngine) step(q Question, ctx *FlowContext) (*Next, error) {
	if ifq, ok := q.(*IfQuestion); ok {
		if ifq.Condition == nil {
			return nil, fmt.Errorf("flow: IfQuestion %q has nil Condition", ifq.ID())
		}
		if ifq.Condition(ctx) {
			return ifq.Then, nil
		}
		return ifq.Else, nil
	}

	q = e.applyAddOptions(q)

	if e.Adapter == nil {
		return nil, fmt.Errorf("flow: no adapter configured")
	}

	answer, err := e.Adapter.Ask(q, ctx)
	if err != nil {
		return nil, fmt.Errorf("flow: ask %q: %w", q.ID(), err)
	}

	ctx.Answers[q.ID()] = answer
	return q.Next(answer), nil
}

// applyReplace returns the first registered Replace question for q.ID(),
// or q unchanged if none.
func (e *FlowEngine) applyReplace(q Question) Question {
	if e.Hooks == nil {
		return q
	}
	replace, _, _ := e.Hooks.byKind(q.ID())
	if len(replace) == 0 {
		return q
	}
	return replace[0]
}

// applyAddOptions returns a clone of q with plugin-added options appended,
// or q unchanged if none apply or q is not an OptionQuestion.
func (e *FlowEngine) applyAddOptions(q Question) Question {
	opt, ok := q.(*OptionQuestion)
	if !ok || e.Hooks == nil {
		return q
	}
	_, extras, _ := e.Hooks.byKind(q.ID())
	if len(extras) == 0 {
		return q
	}
	merged := *opt
	merged.Options = append(append([]*Option(nil), opt.Options...), extras...)
	return &merged
}

// applyInsertAfter splices any InsertAfter injections between q and its
// natural Next. Multiple inserts run in registration order, each chaining
// to the next via the continuation stack; the last hands back to original.
func (e *FlowEngine) applyInsertAfter(targetID string, original *Next, state *traversalState) *Next {
	if e.Hooks == nil {
		return original
	}
	_, _, inserts := e.Hooks.byKind(targetID)
	if len(inserts) == 0 {
		return original
	}

	state.push(original)
	for i := len(inserts) - 1; i > 0; i-- {
		state.push(&Next{Question: inserts[i]})
	}
	return &Next{Question: inserts[0]}
}

// resolveNext converts an edge into the next concrete Question. End edges
// pop the continuation stack so InsertAfter chains return to the parent flow.
func (e *FlowEngine) resolveNext(next *Next, ctx *FlowContext, state *traversalState) Question {
	for next != nil {
		if next.End {
			cont := state.pop()
			if cont == nil {
				return nil
			}
			next = cont
			continue
		}
		if next.Question != nil {
			return next.Question
		}
		if next.Fragment != "" {
			next = e.Fragments.Resolve(next.Fragment, ctx)
			continue
		}
		return nil
	}
	cont := state.pop()
	if cont == nil {
		return nil
	}
	return e.resolveNext(cont, ctx, state)
}

func pluginNames(ids []PluginID) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, string(id))
	}
	return out
}
