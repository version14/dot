package cli

import (
	"testing"

	"github.com/version14/dot/flows"
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

	// Both A options and B's single option converge to C, so C remains
	// unconditionally visible after same-target branch collapse.
	if len(cSlot.conditions) != 1 {
		t.Fatalf("Expected 1 collapsed condition for C, got %d", len(cSlot.conditions))
	}
	if buildHideFunc(cSlot.conditions, newLiveStore())() {
		t.Fatalf("Expected collapsed condition to keep C visible, got %+v", cSlot.conditions[0])
	}
}

func TestFormWalkerFrontendStylingBeforeState(t *testing.T) {
	def := flows.FrontendFlow()
	walker := newFormWalker(nil, nil)
	walker.walk(def.Root)

	index := make(map[string]int, len(walker.slots))
	for i, slot := range walker.slots {
		index[slot.question.ID()] = i
	}

	stylingIdx, ok := index["frontend-styling"]
	if !ok {
		t.Fatal("frontend-styling not found in walker slots")
	}
	stateIdx, ok := index["frontend-state"]
	if !ok {
		t.Fatal("frontend-state not found in walker slots")
	}
	if stylingIdx > stateIdx {
		t.Fatalf("frontend-styling should be before frontend-state, got styling=%d state=%d", stylingIdx, stateIdx)
	}
}
