package flow

// FlowContext — accumulates answers as recursive tree + visited node list during traversal

type FlowContext struct {
	Answers       map[string]AnswerNode // recursive answer tree (built from Result after Run)
	LoopStack     []LoopFrame           // scope stack — one frame per active loop level
	VisitedNodes  []string              // traversal path (used for generator resolution)
	LoadedPlugins []string              // which plugins contributed to the flow
}

// Note: no History field — back navigation within non-loop questions is handled
// natively by Huh's single-form model. Cross-loop back is not supported in v1.

type LoopFrame struct {
	QuestionID string                // which LoopQuestion we're inside
	Index      int                   // current iteration (0-based)
	Answers    map[string]AnswerNode // answers for this iteration
}
