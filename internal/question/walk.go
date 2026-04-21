package question

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
