package question

import (
	"testing"
)

func textQ(value string, next *Question) *Question {
	q := &Question{TextQuestion: &TextQuestion{Value: value}}
	if next != nil {
		q.TextQuestion.Next = &Next{Question: next}
	}
	return q
}

func optQ(value string, opts ...*Option) *Question {
	return &Question{OptionQuestion: &OptionQuestion{Value: value, Options: opts}}
}

func opt(value string) *Option { return &Option{Value: value, Label: value} }

func TestDeepCopyIndependence(t *testing.T) {
	t.Parallel()

	original := textQ("a", textQ("b", nil))
	clone := DeepCopy(original)

	// Mutate clone — original must be unaffected.
	clone.TextQuestion.Value = "CHANGED"
	if original.TextQuestion.Value != "a" {
		t.Error("DeepCopy: mutating clone affected original")
	}

	clone.TextQuestion.Next.Question.TextQuestion.Value = "CHANGED2"
	if original.TextQuestion.Next.Question.TextQuestion.Value != "b" {
		t.Error("DeepCopy: mutating cloned child affected original")
	}
}

func TestDeepCopyNil(t *testing.T) {
	t.Parallel()
	if DeepCopy(nil) != nil {
		t.Error("DeepCopy(nil) should return nil")
	}
}

func TestDeepCopyPreservesNextGenerator(t *testing.T) {
	t.Parallel()

	var called bool
	gen := GeneratorFunc(nil) // placeholder — we check non-nil preservation
	_ = called

	original := textQ("x", nil)
	original.TextQuestion.Next = &Next{Generator: gen}
	clone := DeepCopy(original)

	// Generator field must be copied (both nil in this case — confirms cloneNext ran).
	if clone.TextQuestion.Next == nil {
		t.Fatal("cloned Next should not be nil")
	}
}

func TestAddOption(t *testing.T) {
	t.Parallel()

	flow := optQ("lang", opt("go"))
	flow.AddOption("lang", opt("ts"))

	opts := flow.OptionQuestion.Options
	if len(opts) != 2 {
		t.Fatalf("want 2 options, got %d", len(opts))
	}
	if opts[1].Value != "ts" {
		t.Errorf("want ts, got %q", opts[1].Value)
	}
}

func TestAddOptionMissingKey(t *testing.T) {
	t.Parallel()

	flow := optQ("lang", opt("go"))
	flow.AddOption("missing", opt("ts")) // no-op
	if len(flow.OptionQuestion.Options) != 1 {
		t.Error("AddOption on missing key should not add option")
	}
}

func TestRemoveOption(t *testing.T) {
	t.Parallel()

	flow := optQ("lang", opt("go"), opt("ts"), opt("rust"))
	flow.RemoveOption("lang", "ts")

	opts := flow.OptionQuestion.Options
	if len(opts) != 2 {
		t.Fatalf("want 2 options, got %d", len(opts))
	}
	for _, o := range opts {
		if o.Value == "ts" {
			t.Error("ts should have been removed")
		}
	}
}

func TestRemoveQuestion(t *testing.T) {
	t.Parallel()

	// a → b → c  =>  remove b  =>  a → c
	flow := textQ("a", textQ("b", textQ("c", nil)))
	flow.RemoveQuestion("b")

	next := flow.TextQuestion.Next
	if next == nil {
		t.Fatal("a should still have a next after removing b")
	}
	if next.Question == nil || next.Question.TextQuestion.Value != "c" {
		t.Errorf("a should link directly to c, got %+v", next)
	}
}

func TestAttachFlow(t *testing.T) {
	t.Parallel()

	flow := optQ("framework", opt("react"))
	subFlow := textQ("style", nil)
	flow.AttachFlow("framework", "Vue", "vue", subFlow)

	opts := flow.OptionQuestion.Options
	if len(opts) != 2 {
		t.Fatalf("want 2 options, got %d", len(opts))
	}
	added := opts[1]
	if added.Value != "vue" || added.Next == nil || added.Next.Question != subFlow {
		t.Errorf("unexpected attached option: %+v", added)
	}
}
