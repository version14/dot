package question

type Comparison string

const (
	ComparisonEqual        Comparison = "equal"
	ComparisonNotEqual     Comparison = "not_equal"
	ComparisonGreaterThan  Comparison = "greater_than"
	ComparisonLessThan     Comparison = "less_than"
	ComparisonGreaterEqual Comparison = "greater_equal"
	ComparisonLessEqual    Comparison = "less_equal"
)

// IfAction branches on a previous answer. The runner looks up the answer for
// Key, compares it to Value using Comparison, then traverses Then or Else.
// Either branch may be nil (nil = skip that branch).
type IfAction struct {
	Key        string
	Comparison Comparison
	Value      string
	Then       *Question
	Else       *Question
}

// Builder

type IfBuilder struct {
	ifAction *IfAction
}

// If starts a conditional branch on a previous answer.
func If(key string, comp Comparison, value string) *IfBuilder {
	return &IfBuilder{ifAction: &IfAction{Key: key, Comparison: comp, Value: value}}
}

func (b *IfBuilder) Then(q *Question) *IfBuilder {
	b.ifAction.Then = q
	return b
}

func (b *IfBuilder) Else(q *Question) *IfBuilder {
	b.ifAction.Else = q
	return b
}

func (b *IfBuilder) Q() *Question {
	return &Question{IfQuestion: b.ifAction}
}
