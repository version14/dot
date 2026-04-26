package spec

// FlowContext → ProjectSpec — serializes the answer tree + attaches generator constraints

// Converts flat Result into the recursive answer tree
func Build(result *scaffold.Result, flowID string) *ProjectSpec {
	answers := map[string]AnswerNode{}
	for _, entry := range result.Entries {
		if entry.Iterations != nil {
			answers[entry.Key] = entry.Iterations // []map[string]interface{}
		} else if entry.Multi != nil {
			answers[entry.Key] = entry.Multi
		} else {
			answers[entry.Key] = entry.Value
		}
	}
	return &ProjectSpec{
		FlowID:  flowID,
		Answers: answers,
		// ...
	}
}
