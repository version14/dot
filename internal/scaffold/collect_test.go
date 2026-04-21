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
	acts := Collect(flow, result)
	if len(acts) != 1 {
		t.Fatalf("want 1 activation, got %d", len(acts))
	}
	if acts[0].QuestionKey != "name" || acts[0].AnswerValue != "alice" {
		t.Errorf("unexpected activation: %+v", acts[0])
	}
	// Fn must be callable
	if _, _, err := acts[0].Fn(spec.Spec{}); err != nil {
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
	acts := Collect(flow, result)
	if len(acts) != 1 {
		t.Fatalf("want 1 activation, got %d", len(acts))
	}
	if _, _, err := acts[0].Fn(spec.Spec{}); err != nil {
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
	acts := Collect(flow, result)
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
	acts := Collect(flow, result)
	if len(acts) != 1 {
		t.Fatalf("want 1 activation, got %d", len(acts))
	}
	if _, _, err := acts[0].Fn(spec.Spec{}); err != nil {
		t.Fatal(err)
	}
	if branch != "then" {
		t.Errorf("want then branch, got %q", branch)
	}
}

func TestCollectNilFlow(t *testing.T) {
	t.Parallel()
	acts := Collect(nil, &Result{})
	if len(acts) != 0 {
		t.Errorf("want empty, got %v", acts)
	}
}
