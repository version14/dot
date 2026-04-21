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

// Internal traversal

// walk visits every Question node depth-first. callback returns true to stop early.
func walk(question *Question, callback func(*Question) bool) bool {
	if question == nil {
		return false
	}
	if callback(question) {
		return true
	}
	if question.TextQuestion != nil && question.TextQuestion.Next != nil {
		if walk(question.TextQuestion.Next.Question, callback) {
			return true
		}
	}
	if question.OptionQuestion != nil {
		for _, option := range question.OptionQuestion.Options {
			if option.Next != nil && walk(option.Next.Question, callback) {
				return true
			}
		}
		if question.OptionQuestion.Next != nil && walk(question.OptionQuestion.Next.Question, callback) {
			return true
		}
	}
	if question.LoopQuestion != nil && walk(question.LoopQuestion.Question, callback) {
		return true
	}
	if question.IfQuestion != nil {
		if walk(question.IfQuestion.Then, callback) {
			return true
		}
		if walk(question.IfQuestion.Else, callback) {
			return true
		}
	}
	return false
}

func splice(question *Question, key string) {
	if question == nil {
		return
	}
	if question.TextQuestion != nil && question.TextQuestion.Next != nil {
		if cutNext(&question.TextQuestion.Next, key) {
			return
		}
		splice(question.TextQuestion.Next.Question, key)
	}
	if question.OptionQuestion != nil {
		for _, option := range question.OptionQuestion.Options {
			if cutNext(&option.Next, key) {
				return
			}
			if option.Next != nil {
				splice(option.Next.Question, key)
			}
		}
		if cutNext(&question.OptionQuestion.Next, key) {
			return
		}
		if question.OptionQuestion.Next != nil {
			splice(question.OptionQuestion.Next.Question, key)
		}
	}
	if question.LoopQuestion != nil {
		splice(question.LoopQuestion.Question, key)
	}
	if question.IfQuestion != nil {
		splice(question.IfQuestion.Then, key)
		splice(question.IfQuestion.Else, key)
	}
}

func cutNext(next **Next, key string) bool {
	if next == nil || *next == nil || (*next).Question == nil {
		return false
	}
	if questionKey((*next).Question) == key {
		*next = questionSuccessor((*next).Question)
		return true
	}
	return false
}

func questionKey(question *Question) string {
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

func questionSuccessor(question *Question) *Next {
	switch {
	case question.TextQuestion != nil:
		return question.TextQuestion.Next
	case question.OptionQuestion != nil:
		return question.OptionQuestion.Next
	}
	return nil
}

// Deep copy

func cloneQuestion(question Question) Question {
	out := Question{}
	if question.TextQuestion != nil {
		textQuestion := *question.TextQuestion
		if textQuestion.Next != nil {
			next := cloneNext(textQuestion.Next)
			textQuestion.Next = &next
		}
		out.TextQuestion = &textQuestion
	}
	if question.OptionQuestion != nil {
		optionQuestion := *question.OptionQuestion
		options := make([]*Option, len(optionQuestion.Options))
		for i, option := range optionQuestion.Options {
			choice := *option
			if choice.Next != nil {
				next := cloneNext(choice.Next)
				choice.Next = &next
			}
			options[i] = &choice
		}
		optionQuestion.Options = options
		if optionQuestion.Next != nil {
			next := cloneNext(optionQuestion.Next)
			optionQuestion.Next = &next
		}
		out.OptionQuestion = &optionQuestion
	}
	if question.LoopQuestion != nil {
		loopQuestion := *question.LoopQuestion
		loopQuestion.Question = DeepCopy(loopQuestion.Question)
		out.LoopQuestion = &loopQuestion
	}
	if question.IfQuestion != nil {
		ifAction := *question.IfQuestion
		ifAction.Then = DeepCopy(ifAction.Then)
		ifAction.Else = DeepCopy(ifAction.Else)
		out.IfQuestion = &ifAction
	}
	return out
}

func cloneNext(next *Next) Next {
	out := Next{}
	out.Question = DeepCopy(next.Question)
	return out
}
