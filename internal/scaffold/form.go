package scaffold

import (
	"github.com/charmbracelet/huh"
	q "github.com/version14/dot/internal/question"
)

// nodeKey extracts the value key from a question node.
func nodeKey(question *q.Question) string {
	switch {
	case question.TextQuestion != nil:
		return question.TextQuestion.Value
	case question.OptionQuestion != nil:
		return question.OptionQuestion.Value
	case question.LoopQuestion != nil:
		return question.LoopQuestion.Value
	}
	return ""
}

// buildGroups converts flat nodes into huh groups.
// Conditional groups have a WithHideFunc that re-evaluates live.
func (runner *Runner) buildGroups(nodes []flatNode) []*huh.Group {
	groups := make([]*huh.Group, 0, len(nodes))
	for _, node := range nodes {
		field := runner.buildField(node.question)
		if field == nil {
			continue
		}
		conditions := node.conditions // captured — must not be aliased across iterations
		group := huh.NewGroup(field)
		if len(conditions) > 0 {
			group = group.WithHideFunc(func() bool {
				return !runner.condsMatch(conditions)
			})
		}
		groups = append(groups, group)
	}
	return groups
}

// buildField returns the huh field for a question node.
// Pointers are reused when the same key appears in multiple branches so that
// all huh fields for that key share one live value — whichever branch is
// actually visible will write the correct answer into the shared pointer.
func (runner *Runner) buildField(question *q.Question) huh.Field {
	if question.TextQuestion != nil {
		return runner.buildInputField(question.TextQuestion)
	}
	if question.OptionQuestion != nil {
		return runner.buildSelectField(question.OptionQuestion)
	}
	return nil
}

func (runner *Runner) buildInputField(tq *q.TextQuestion) huh.Field {
	part, ok := runner.strPtrs[tq.Value]
	if !ok {
		part = new(string)
		runner.strPtrs[tq.Value] = part
	}
	return huh.NewInput().
		Title(tq.Label).
		Description(tq.Description).
		Placeholder(tq.Placeholder).
		Value(part)
}

func (runner *Runner) buildSelectField(oq *q.OptionQuestion) huh.Field {
	opts := make([]huh.Option[string], len(oq.Options))
	for i, option := range oq.Options {
		opts[i] = huh.NewOption(option.Label, option.Value)
	}
	if oq.Multiple {
		part, ok := runner.multiPtrs[oq.Value]
		if !ok {
			part = new([]string)
			*part = []string{}
			runner.multiPtrs[oq.Value] = part
		}
		return huh.NewMultiSelect[string]().
			Title(oq.Label).
			Description(oq.Description).
			Options(opts...).
			Value(part)
	}
	part, ok := runner.strPtrs[oq.Value]
	if !ok {
		part = new(string)
		runner.strPtrs[oq.Value] = part
	}
	return huh.NewSelect[string]().
		Title(oq.Label).
		Description(oq.Description).
		Options(opts...).
		Value(part)
}
