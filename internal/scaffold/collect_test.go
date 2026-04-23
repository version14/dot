package scaffold

import (
	"testing"

	"github.com/version14/dot/internal/generator"
	q "github.com/version14/dot/internal/question"
	"github.com/version14/dot/internal/spec"
)

func resultWith(kvs ...string) *Result {
	r := &Result{}
	for i := 0; i < len(kvs)-1; i += 2 {
		r.add(AnswerEntry{Key: kvs[i], Value: kvs[i+1]})
	}
	return r
}

func TestCollectText(t *testing.T) {
	t.Parallel()

	var called bool
	gen := func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
		called = true
		return nil, nil, nil
	}

	flow := &q.Question{
		TextQuestion: &q.TextQuestion{
			Value: "name",
			Next:  &q.Next{Generator: gen, Question: nil},
		},
	}
	result := resultWith("name", "alice")
	acts := Collect(flow, result, spec.Spec{})
	if len(acts) != 1 {
		t.Fatalf("want 1 activation, got %d", len(acts))
	}
	if acts[0].QuestionKey != "name" || acts[0].AnswerValue != "alice" {
		t.Errorf("unexpected activation: %+v", acts[0])
	}
	if _, _, err := acts[0].Fn(acts[0].Spec); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("generator was not called")
	}
}

func TestCollectOptionTakenBranch(t *testing.T) {
	t.Parallel()

	var activated string
	makeGen := func(name string) q.GeneratorFunc {
		return func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
			activated = name
			return nil, nil, nil
		}
	}

	flow := &q.Question{
		OptionQuestion: &q.OptionQuestion{
			Value: "lang",
			Options: []*q.Option{
				{Value: "go", Next: &q.Next{Generator: makeGen("go")}},
				{Value: "ts", Next: &q.Next{Generator: makeGen("ts")}},
			},
		},
	}
	result := resultWith("lang", "go")
	acts := Collect(flow, result, spec.Spec{})
	if len(acts) != 1 {
		t.Fatalf("want 1 activation, got %d", len(acts))
	}
	if _, _, err := acts[0].Fn(acts[0].Spec); err != nil {
		t.Fatal(err)
	}
	if activated != "go" {
		t.Errorf("want go generator, got %q", activated)
	}
}

func TestCollectOptionNoGenerator(t *testing.T) {
	t.Parallel()

	flow := &q.Question{
		OptionQuestion: &q.OptionQuestion{
			Value: "lang",
			Options: []*q.Option{
				{Value: "go", Next: &q.Next{Question: nil}},
			},
		},
	}
	result := resultWith("lang", "go")
	acts := Collect(flow, result, spec.Spec{})
	if len(acts) != 0 {
		t.Fatalf("want 0 activations, got %d", len(acts))
	}
}

func TestCollectIfThen(t *testing.T) {
	t.Parallel()

	var branch string
	flow := &q.Question{
		IfQuestion: &q.IfAction{
			Key:        "env",
			Comparison: q.ComparisonEqual,
			Value:      "prod",
			Then: &q.Question{
				TextQuestion: &q.TextQuestion{
					Value: "domain",
					Next: &q.Next{Generator: func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
						branch = "then"
						return nil, nil, nil
					}},
				},
			},
			Else: &q.Question{
				TextQuestion: &q.TextQuestion{
					Value: "port",
					Next: &q.Next{Generator: func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
						branch = "else"
						return nil, nil, nil
					}},
				},
			},
		},
	}

	result := resultWith("env", "prod", "domain", "example.com")
	acts := Collect(flow, result, spec.Spec{})
	if len(acts) != 1 {
		t.Fatalf("want 1 activation, got %d", len(acts))
	}
	if _, _, err := acts[0].Fn(acts[0].Spec); err != nil {
		t.Fatal(err)
	}
	if branch != "then" {
		t.Errorf("want then branch, got %q", branch)
	}
}

func TestCollectNilFlow(t *testing.T) {
	t.Parallel()
	acts := Collect(nil, &Result{}, spec.Spec{})
	if len(acts) != 0 {
		t.Errorf("want empty, got %v", acts)
	}
}

// TestCollectLoopScopedSpecs verifies that each loop iteration produces an
// activation whose Spec.Extensions carries that iteration's answers only.
func TestCollectLoopScopedSpecs(t *testing.T) {
	t.Parallel()

	makeBodyGen := func(seen *[]string) q.GeneratorFunc {
		return func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
			arch, _ := s.Extensions["arch"].(string)
			*seen = append(*seen, arch)
			return nil, nil, nil
		}
	}

	var seen []string
	body := &q.Question{
		OptionQuestion: &q.OptionQuestion{
			Value: "arch",
			Options: []*q.Option{
				{Value: "mvc", Next: &q.Next{Generator: makeBodyGen(&seen)}},
				{Value: "clean", Next: &q.Next{Generator: makeBodyGen(&seen)}},
			},
		},
	}

	flow := q.Loop("services", "services", body)

	result := &Result{}
	result.add(AnswerEntry{
		Key:   "services",
		Value: "2",
		Iterations: [][]AnswerEntry{
			{{Key: "arch", Value: "mvc"}},
			{{Key: "arch", Value: "clean"}},
		},
	})

	base := spec.Spec{Extensions: map[string]any{"project-name": "demo"}}
	acts := Collect(flow, result, base)
	if len(acts) != 2 {
		t.Fatalf("want 2 activations (one per iteration), got %d", len(acts))
	}

	// Each activation's Spec must carry its own iteration's answer, plus base.
	iter0, _ := acts[0].Spec.Extensions["arch"].(string)
	iter1, _ := acts[1].Spec.Extensions["arch"].(string)
	if iter0 != "mvc" || iter1 != "clean" {
		t.Errorf("want [mvc clean], got [%s %s]", iter0, iter1)
	}
	if name, _ := acts[0].Spec.Extensions["project-name"].(string); name != "demo" {
		t.Errorf("want outer answer visible inside loop, got %q", name)
	}

	// Running the activations must see scoped specs, not shared state.
	for _, a := range acts {
		if _, _, err := a.Fn(a.Spec); err != nil {
			t.Fatal(err)
		}
	}
	if len(seen) != 2 || seen[0] != "mvc" || seen[1] != "clean" {
		t.Errorf("generators saw wrong scope: %v", seen)
	}
}

// TestCollectLoopIterationsDoNotLeak verifies that iteration specs don't
// share the same underlying map: mutating one must not affect the others or
// the base.
func TestCollectLoopIterationsDoNotLeak(t *testing.T) {
	t.Parallel()

	body := &q.Question{
		TextQuestion: &q.TextQuestion{
			Value: "name",
			Next: &q.Next{Generator: func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
				return nil, nil, nil
			}},
		},
	}
	flow := q.Loop("apps", "apps", body)

	result := &Result{}
	result.add(AnswerEntry{
		Key:   "apps",
		Value: "2",
		Iterations: [][]AnswerEntry{
			{{Key: "name", Value: "a"}},
			{{Key: "name", Value: "b"}},
		},
	})

	base := spec.Spec{Extensions: map[string]any{"shared": "base"}}
	acts := Collect(flow, result, base)
	if len(acts) != 2 {
		t.Fatalf("want 2 activations, got %d", len(acts))
	}

	// Mutating iteration 0's map must not appear in iteration 1 or base.
	acts[0].Spec.Extensions["injected"] = "leak"
	if _, ok := acts[1].Spec.Extensions["injected"]; ok {
		t.Error("iteration 1 shares a map with iteration 0")
	}
	if _, ok := base.Extensions["injected"]; ok {
		t.Error("base.Extensions was mutated via an iteration")
	}
}

func TestEvaluateCondition(t *testing.T) {
	t.Parallel()

	cases := []struct {
		actual, expected string
		comp             q.Comparison
		want             bool
	}{
		{"5", "5", q.ComparisonEqual, true},
		{"5", "6", q.ComparisonEqual, false},
		{"5", "6", q.ComparisonNotEqual, true},
		{"5", "5", q.ComparisonNotEqual, false},
		{"7", "5", q.ComparisonGreaterThan, true},
		{"5", "7", q.ComparisonGreaterThan, false},
		{"5", "5", q.ComparisonGreaterEqual, true},
		{"4", "5", q.ComparisonGreaterEqual, false},
		{"3", "5", q.ComparisonLessThan, true},
		{"5", "3", q.ComparisonLessThan, false},
		{"5", "5", q.ComparisonLessEqual, true},
		{"6", "5", q.ComparisonLessEqual, false},
		// lexicographic fallback
		{"b", "a", q.ComparisonGreaterThan, true},
		{"a", "b", q.ComparisonLessThan, true},
	}
	for _, tc := range cases {
		got := evaluateCondition(tc.actual, tc.comp, tc.expected)
		if got != tc.want {
			t.Errorf("evaluateCondition(%q, %s, %q) = %v, want %v",
				tc.actual, tc.comp, tc.expected, got, tc.want)
		}
	}
}
