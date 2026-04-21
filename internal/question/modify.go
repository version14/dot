package question

// DeepCopy returns a fully independent copy of question.
// Use this before mutating a shared package-level var.
func DeepCopy(question *Question) *Question {
	if question == nil {
		return nil
	}
	clone := cloneQuestion(*question)
	return &clone
}

// AddOption finds the OptionQuestion with questionKey and appends option.
func (question *Question) AddOption(questionKey string, option *Option) *Question {
	walk(question, func(next *Question) bool {
		if next.OptionQuestion != nil && next.OptionQuestion.Value == questionKey {
			next.OptionQuestion.Options = append(next.OptionQuestion.Options, option)
			return true
		}
		return false
	})
	return question
}

// RemoveOption removes the option with optionValue from the OptionQuestion
// identified by questionKey.
func (question *Question) RemoveOption(questionKey, optionValue string) *Question {
	walk(question, func(next *Question) bool {
		if next.OptionQuestion == nil || next.OptionQuestion.Value != questionKey {
			return false
		}
		kept := next.OptionQuestion.Options[:0]
		for _, option := range next.OptionQuestion.Options {
			if option.Value != optionValue {
				kept = append(kept, option)
			}
		}
		next.OptionQuestion.Options = kept
		return true
	})
	return question
}

// RemoveQuestion bypasses the question with the given key, linking its
// predecessor directly to its successor.
func (question *Question) RemoveQuestion(key string) *Question {
	splice(question, key)
	return question
}

// AttachFlow adds an option with a sub-flow at the OptionQuestion identified
// by questionKey. Shorthand for AddOption with a Next-wrapped flow.
func (question *Question) AttachFlow(questionKey, label, value string, flow *Question) *Question {
	return question.AddOption(questionKey, &Option{
		Label: label,
		Value: value,
		Next:  &Next{Question: flow},
	})
}
