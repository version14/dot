package flow

// Scope chain resolver — flattens answer tree into ctx.Answers for a specific loop depth

// FlattenScope walks the LoopStack from outermost to innermost,
// merging answers at each level. Deeper scopes win on key conflicts.
func FlattenScope(global map[string]AnswerNode, stack []LoopFrame) map[string]interface{} {
	result := map[string]interface{}{}
	// Start with global
	for k, v := range global {
		result[k] = v
	}
	// Each loop frame overrides
	for _, frame := range stack {
		for k, v := range frame.Answers {
			result[k] = v
		}
	}
	return result
}
