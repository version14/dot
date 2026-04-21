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
	textQeustion *TextQuestion
}

func Text(label, value string) *TextBuilder {
	return &TextBuilder{textQeustion: &TextQuestion{Label: label, Value: value}}
}

func (b *TextBuilder) Description(d string) *TextBuilder {
	b.textQeustion.Description = d
	return b
}

func (b *TextBuilder) Placeholder(p string) *TextBuilder {
	b.textQeustion.Placeholder = p
	return b
}

func (b *TextBuilder) Then(next *Question) *TextBuilder {
	b.textQeustion.Next = &Next{Question: next}
	return b
}

func (b *TextBuilder) Q() *Question {
	return &Question{TextQuestion: b.textQeustion}
}
