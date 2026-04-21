package question

// Type

type TextQuestion struct {
	Label       string
	Description string
	Placeholder string
	Value       string
	Next        *Next
}

// Builder

type TextBuilder struct {
	textQuestion *TextQuestion
}

func Text(label, value string) *TextBuilder {
	return &TextBuilder{textQuestion: &TextQuestion{Label: label, Value: value}}
}

func (b *TextBuilder) Description(d string) *TextBuilder {
	b.textQuestion.Description = d
	return b
}

func (b *TextBuilder) Placeholder(p string) *TextBuilder {
	b.textQuestion.Placeholder = p
	return b
}

func (b *TextBuilder) Then(next *Question) *TextBuilder {
	b.textQuestion.Next = &Next{Question: next}
	return b
}

func (b *TextBuilder) Q() *Question {
	return &Question{TextQuestion: b.textQuestion}
}
