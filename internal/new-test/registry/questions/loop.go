package questions

// LoopAction repeats Question N times. The runner first prompts for the
// count (using Label / Description / Value as the integer input prompt),
// stores it under Value, then runs Question N times.
type LoopAction struct {
	Label       string
	Description string
	Value       string
	Question    *Question
}
