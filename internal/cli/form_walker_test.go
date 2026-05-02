package cli

import (
	"testing"

	"github.com/version14/dot/internal/flow"
)

func TestFormWalkerDiamondConvergence(t *testing.T) {
	c := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{ID_: "C"},
		Label:        "Question C",
	}

	b := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "B"},
		Label:        "Question B",
		Options: []*flow.Option{
			{Label: "B1", Value: "b1", Next: &flow.Next{Question: c}},
		},
	}

	a := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "A"},
		Label:        "Question A",
		Options: []*flow.Option{
			{Label: "A1", Value: "a1", Next: &flow.Next{Question: b}},
			{Label: "A2", Value: "a2", Next: &flow.Next{Question: b}},
		},
	}

	walker := newFormWalker(nil, nil)
	walker.walk(a)

	// Find slot for C
	var cSlot *formSlot
	for _, s := range walker.slots {
		if s.question.ID() == "C" {
			cSlot = s
			break
		}
	}

	if cSlot == nil {
		t.Fatal("Question C not found in walker slots")
	}

	// C should be reachable via (A=a1 AND B=b1) OR (A=a2 AND B=b1)
	if len(cSlot.conditions) < 2 {
		t.Errorf("Expected at least 2 conditions for C, got %d", len(cSlot.conditions))
		for i, cond := range cSlot.conditions {
			t.Logf("Condition %d: %v", i, cond)
		}
	}
}
