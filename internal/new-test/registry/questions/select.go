package questions

// Builder

type SelectBuilder struct {
	optionQuestion *OptionQuestion
}

func Select(label, value string) *SelectBuilder {
	return &SelectBuilder{optionQuestion: &OptionQuestion{Label: label, Value: value}}
}

func (b *SelectBuilder) Description(d string) *SelectBuilder {
	b.optionQuestion.Description = d
	return b
}

func (b *SelectBuilder) Multiple() *SelectBuilder {
	b.optionQuestion.Multiple = true
	return b
}

// Choice adds an option. next is optional (0 or 1 argument).
func (b *SelectBuilder) Choice(label, value string, next ...*Question) *SelectBuilder {
	opt := &Option{Label: label, Value: value}
	if len(next) > 0 && next[0] != nil {
		opt.Next = &Next{Question: next[0]}
	}
	b.optionQuestion.Options = append(b.optionQuestion.Options, opt)
	return b
}

// ChoiceGen adds a choice that triggers a generator when selected.
func (b *SelectBuilder) ChoiceGen(label, value, generatorID string, next ...*Question) *SelectBuilder {
	opt := &Option{Label: label, Value: value, GeneratorID: generatorID}
	if len(next) > 0 && next[0] != nil {
		opt.Next = &Next{Question: next[0]}
	}
	b.optionQuestion.Options = append(b.optionQuestion.Options, opt)
	return b
}

// ChoiceWithGen adds a choice that calls fn when the path is taken.
// fn is invoked by core.Collect after the survey completes.
func (b *SelectBuilder) ChoiceWithGen(label, value string, fn GeneratorFunc, next ...*Question) *SelectBuilder {
	opt := &Option{Label: label, Value: value}
	opt.Next = &Next{Generator: fn}
	if len(next) > 0 && next[0] != nil {
		opt.Next.Question = next[0]
	}
	b.optionQuestion.Options = append(b.optionQuestion.Options, opt)
	return b
}

// Then sets the shared continuation taken after the question regardless of
// which option was chosen. Required for Multiple() questions.
func (b *SelectBuilder) Then(next *Question) *SelectBuilder {
	b.optionQuestion.Next = &Next{Question: next}
	return b
}

func (b *SelectBuilder) Q() *Question {
	return &Question{OptionQuestion: b.optionQuestion}
}
