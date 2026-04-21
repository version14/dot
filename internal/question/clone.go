package question

func cloneQuestion(question Question) Question {
	out := Question{}
	if question.TextQuestion != nil {
		tq := *question.TextQuestion
		if tq.Next != nil {
			next := cloneNext(tq.Next)
			tq.Next = &next
		}
		out.TextQuestion = &tq
	}
	if question.OptionQuestion != nil {
		oq := *question.OptionQuestion
		options := make([]*Option, len(oq.Options))
		for i, option := range oq.Options {
			choice := *option
			if choice.Next != nil {
				next := cloneNext(choice.Next)
				choice.Next = &next
			}
			options[i] = &choice
		}
		oq.Options = options
		if oq.Next != nil {
			next := cloneNext(oq.Next)
			oq.Next = &next
		}
		out.OptionQuestion = &oq
	}
	if question.LoopQuestion != nil {
		lq := *question.LoopQuestion
		lq.Question = DeepCopy(lq.Question)
		out.LoopQuestion = &lq
	}
	if question.IfQuestion != nil {
		ia := *question.IfQuestion
		ia.Then = DeepCopy(ia.Then)
		ia.Else = DeepCopy(ia.Else)
		out.IfQuestion = &ia
	}
	return out
}

func cloneNext(next *Next) Next {
	return Next{
		Generator: next.Generator,
		Question:  DeepCopy(next.Question),
	}
}
