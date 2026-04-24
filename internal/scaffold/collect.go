package scaffold

import (
	"strconv"

	"github.com/version14/dot/internal/question"
	"github.com/version14/dot/internal/spec"
)

// Activation records a generator that was activated by the user's answer.
//
// Spec is the scoped view the generator should see: the base spec with any
// outer-scope answers, overlaid by iteration answers when the activation lives
// inside a loop body. Each activation carries its own Spec so that two
// iterations of the same loop can't see each other's answers.
type Activation struct {
	QuestionKey string
	AnswerValue string
	Fn          question.GeneratorFunc
	Spec        spec.Spec
}

// Collect walks the question tree along the path recorded in result and returns
// every generator attached to the steps the user actually took.
//
// base is the spec that applies outside of any loop. For loop bodies, Collect
// projects a new spec per iteration by overlaying that iteration's answers on
// top of base.Extensions, so each generator activation sees only its own scope.
func Collect(flow *question.Question, result *Result, base spec.Spec) []Activation {
	var acts []Activation
	collectStep(flow, result, base, &acts)
	return acts
}

func collectStep(question *question.Question, result *Result, scope spec.Spec, acts *[]Activation) {
	if question == nil {
		return
	}

	if question.TextQuestion != nil {
		textQuestion := question.TextQuestion
		if textQuestion.Next != nil {
			if textQuestion.Next.Generator != nil {
				*acts = append(*acts, Activation{
					QuestionKey: textQuestion.Value,
					AnswerValue: result.Get(textQuestion.Value),
					Fn:          textQuestion.Next.Generator,
					Spec:        scope,
				})
			}
			collectStep(textQuestion.Next.Question, result, scope, acts)
		}
		return
	}

	if question.OptionQuestion != nil {
		optionQuestion := question.OptionQuestion
		answer := result.Get(optionQuestion.Value)
		for _, opt := range optionQuestion.Options {
			if opt.Value != answer || opt.Next == nil {
				continue
			}
			if opt.Next.Generator != nil {
				*acts = append(*acts, Activation{
					QuestionKey: optionQuestion.Value,
					AnswerValue: answer,
					Fn:          opt.Next.Generator,
					Spec:        scope,
				})
			}
			collectStep(opt.Next.Question, result, scope, acts)
		}
		if optionQuestion.Next != nil {
			if optionQuestion.Next.Generator != nil {
				*acts = append(*acts, Activation{
					QuestionKey: optionQuestion.Value,
					AnswerValue: answer,
					Fn:          optionQuestion.Next.Generator,
					Spec:        scope,
				})
			}
			collectStep(optionQuestion.Next.Question, result, scope, acts)
		}
		return
	}

	if question.LoopQuestion != nil {
		loopQuestion := question.LoopQuestion
		if idx, ok := result.index[loopQuestion.Value]; ok {
			for _, iteration := range result.Entries[idx].Iterations {
				subResult := resultFromEntries(iteration)
				subScope := overlayEntries(scope, iteration)
				collectStep(loopQuestion.Question, subResult, subScope, acts)
			}
		}
		return
	}

	if question.IfQuestion != nil {
		ifQuestion := question.IfQuestion
		actual := result.Get(ifQuestion.Key)
		if actual == "" && result.index == nil {
			return
		}
		taken := evaluateCondition(actual, ifQuestion.Comparison, ifQuestion.Value)
		if taken {
			collectStep(ifQuestion.Then, result, scope, acts)
		} else {
			collectStep(ifQuestion.Else, result, scope, acts)
		}
	}
}

// evaluateCondition mirrors the IfAction comparison logic. Ordered comparisons
// try numeric parse first and fall back to lexicographic comparison.
func evaluateCondition(actual string, comp question.Comparison, expected string) bool {
	switch comp {
	case question.ComparisonEqual:
		return actual == expected
	case question.ComparisonNotEqual:
		return actual != expected
	case question.ComparisonGreaterThan:
		return compare(actual, expected) > 0
	case question.ComparisonLessThan:
		return compare(actual, expected) < 0
	case question.ComparisonGreaterEqual:
		return compare(actual, expected) >= 0
	case question.ComparisonLessEqual:
		return compare(actual, expected) <= 0
	default:
		return false
	}
}

// compare returns -1, 0, or 1. Numeric when both sides parse as int,
// lexicographic otherwise.
func compare(a, b string) int {
	ai, aerr := strconv.Atoi(a)
	bi, berr := strconv.Atoi(b)
	if aerr == nil && berr == nil {
		switch {
		case ai < bi:
			return -1
		case ai > bi:
			return 1
		default:
			return 0
		}
	}
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// overlayEntries returns a new spec with entries merged on top of scope.Extensions.
// Project is preserved; iteration values shadow same-key outer values for this scope only.
func overlayEntries(scope spec.Spec, entries []AnswerEntry) spec.Spec {
	ext := make(map[string]any, len(scope.Extensions)+len(entries))
	for k, v := range scope.Extensions {
		ext[k] = v
	}
	for _, e := range entries {
		switch {
		case len(e.Iterations) > 0:
			// Nested loops: store iteration count under the loop key (as a string).
			ext[e.Key] = strconv.Itoa(len(e.Iterations))
		case len(e.Multi) > 0:
			ext[e.Key] = e.Multi
		default:
			ext[e.Key] = e.Value
		}
	}
	return spec.Spec{Project: scope.Project, Extensions: ext}
}

// resultFromEntries rebuilds a Result from a flat entry list (used for loop iterations).
func resultFromEntries(entries []AnswerEntry) *Result {
	result := &Result{}
	for _, entity := range entries {
		result.add(entity)
	}
	return result
}
