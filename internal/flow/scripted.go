package flow

import "fmt"

// ScriptedRunner implements FlowRunner by answering every question from a
// recorded map instead of a terminal. It is how a Flow is replayed headlessly —
// by test-flow from its fixtures, and by fingerprinting to see what each
// Generator contributes under a given set of Answers.
//
// Plugin injections fire via the supplied HookRegistry (and FragmentRegistry),
// so inserted, replaced and added-option questions appear in the engine's
// traversal exactly as they would interactively.
type ScriptedRunner struct {
	adapter   *scriptedAsker
	hooks     *HookRegistry
	fragments *FragmentRegistry
}

// NewScriptedRunner returns a FlowRunner that answers from answers. hooks and
// fragments may be nil.
func NewScriptedRunner(
	answers map[string]Answer,
	hooks *HookRegistry,
	fragments *FragmentRegistry,
) *ScriptedRunner {
	return &ScriptedRunner{
		adapter:   &scriptedAsker{answers: answers},
		hooks:     hooks,
		fragments: fragments,
	}
}

func (r *ScriptedRunner) Run(root Question) (*FlowContext, error) {
	eng := NewEngine(r.adapter)
	if r.hooks != nil {
		eng.Hooks = r.hooks
	}
	if r.fragments != nil {
		eng.Fragments = r.fragments
	}
	return eng.Run(root)
}

// scriptedAsker answers each question from a recorded map.
//
// LoopQuestion handling: when the scripted answer for a loop is a JSON array
// (i.e. []interface{} after unmarshal), each element is the answer-map for one
// iteration of the loop body. The body is walked in order, looking up each child
// question's answer from the per-iteration map. This mirrors
// HuhFormRunner.runLoopSubForms, but synchronously and without any UI.
type scriptedAsker struct {
	answers map[string]Answer
}

func (a *scriptedAsker) Ask(q Question, ctx *FlowContext) (Answer, error) {
	if loop, ok := q.(*LoopQuestion); ok {
		return a.askLoop(loop, ctx)
	}

	id := q.ID()
	ans, ok := a.answers[id]
	if !ok {
		return nil, fmt.Errorf("flow: no scripted answer for question %q", id)
	}
	return ans, nil
}

// askLoop runs the full body sub-graph for each scripted iteration using a fresh
// engine, so conditional questions in the body are followed or skipped based on
// that iteration's answers.
func (a *scriptedAsker) askLoop(loop *LoopQuestion, _ *FlowContext) (Answer, error) {
	raw, ok := a.answers[loop.ID()]
	if !ok {
		return nil, fmt.Errorf("flow: no scripted iterations for loop %q", loop.ID())
	}
	iters, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("flow: loop %q expects an array of objects, got %T", loop.ID(), raw)
	}

	out := make([]map[string]Answer, len(iters))
	for i, iter := range iters {
		iterMap, ok := iter.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("flow: loop %q iteration %d must be an object", loop.ID(), i)
		}

		// Layer iteration answers over global ones: a body question resolves
		// against the iteration map first, then falls back to the outer answers.
		iterAdapter := &scriptedAsker{answers: mergeAnswerMaps(a.answers, iterMap)}
		iterEng := NewEngine(iterAdapter)

		iterAnswers := make(map[string]Answer)
		for _, bodyRoot := range loop.Body {
			bodyCtx, err := iterEng.Run(bodyRoot)
			if err != nil {
				return nil, fmt.Errorf("loop %q iter %d body: %w", loop.ID(), i, err)
			}
			for k, v := range bodyCtx.Answers {
				iterAnswers[k] = v
			}
		}
		out[i] = iterAnswers
	}
	return out, nil
}

func mergeAnswerMaps(base, overlay map[string]Answer) map[string]Answer {
	merged := make(map[string]Answer, len(base)+len(overlay))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overlay {
		merged[k] = v
	}
	return merged
}
