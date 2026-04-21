package scaffold

import (
	"github.com/version14/dot/internal/question"
)

// Activation records a generator that was activated by the user's answer.
type Activation struct {
	QuestionKey string
	AnswerValue string
	Fn          question.GeneratorFunc
}

// Collect walks the question tree along the path recorded in result and returns
// every generator attached to the steps the user actually took.
//
// For loops, it recurses into each iteration's sub-result so that per-service
// or per-app generators are collected with their correct scoped answers.
func Collect(flow *question.Question, result *Result) []Activation {
	var acts []Activation
	collectStep(flow, result, &acts)
	return acts
}

func collectStep(question *question.Question, result *Result, acts *[]Activation) {
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
				})
			}
			collectStep(textQuestion.Next.Question, result, acts)
		}
		return
	}

	if question.OptionQuestion != nil {
		optionQuestion := question.OptionQuestion
		answer := result.Get(optionQuestion.Value)
		// Follow the taken option's branch.
		for _, opt := range optionQuestion.Options {
			if opt.Value != answer || opt.Next == nil {
				continue
			}
			if opt.Next.Generator != nil {
				*acts = append(*acts, Activation{
					QuestionKey: optionQuestion.Value,
					AnswerValue: answer,
					Fn:          opt.Next.Generator,
				})
			}
			collectStep(opt.Next.Question, result, acts)
		}
		// Follow the shared continuation (Multiple selects or .Then()).
		if optionQuestion.Next != nil {
			if optionQuestion.Next.Generator != nil {
				*acts = append(*acts, Activation{
					QuestionKey: optionQuestion.Value,
					AnswerValue: answer,
					Fn:          optionQuestion.Next.Generator,
				})
			}
			collectStep(optionQuestion.Next.Question, result, acts)
		}
		return
	}

	if question.LoopQuestion != nil {
		loopQuestion := question.LoopQuestion
		// Find the loop entry in result and recurse into each iteration.
		if idx, ok := result.index[loopQuestion.Value]; ok {
			for _, iteration := range result.Entries[idx].Iterations {
				subResult := resultFromEntries(iteration)
				collectStep(loopQuestion.Question, subResult, acts)
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
			collectStep(ifQuestion.Then, result, acts)
		} else {
			collectStep(ifQuestion.Else, result, acts)
		}
	}
}

// evaluateCondition mirrors the IfAction comparison logic.
func evaluateCondition(actual string, comp question.Comparison, expected string) bool {
	switch comp {
	case question.ComparisonEqual:
		return actual == expected
	case question.ComparisonNotEqual:
		return actual != expected
	default:
		return false
	}
}

// resultFromEntries rebuilds a Result from a flat entry list (used for loop iterations).
func resultFromEntries(entries []AnswerEntry) *Result {
	result := &Result{}
	for _, entity := range entries {
		result.add(entity)
	}
	return result
}
