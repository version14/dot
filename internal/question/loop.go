package question

// LoopAction repeats Question N times. The runner first prompts for the
// count (using Label / Description / Value as the integer input prompt),
// stores it under Value, then runs Question N times.
type LoopAction struct {
	Label       string
	Description string
	Value       string
	Question    *Question
}

// Builder

// Loop creates a loop question. The runner first asks for the count (label/value),
// then repeats body N times.
func Loop(label, value string, body *Question) *Question {
	return &Question{LoopQuestion: &LoopAction{Label: label, Value: value, Question: body}}
}
