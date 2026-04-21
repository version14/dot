package question

// Question is a sum-type: exactly one of the pointer fields is set.
// The runner dispatches on whichever is non-nil.
type Question struct {
	TextQuestion   *TextQuestion
	OptionQuestion *OptionQuestion
	LoopQuestion   *LoopAction
	IfQuestion     *IfAction
}
