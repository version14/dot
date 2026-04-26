package flow

// Answer is the value returned by a single question node.
// Concrete types: string, bool, int, []string (multi-select)
type Answer = interface{}

// AnswerNode is a recursive tree node in the ProjectSpec.
// string | bool | int | []string | map[string]AnswerNode | []map[string]AnswerNode
type AnswerNode = interface{}
