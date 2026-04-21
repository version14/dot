package questions

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
