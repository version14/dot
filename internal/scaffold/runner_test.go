package scaffold

import (
	"testing"

	q "github.com/version14/dot/internal/question"
)

func TestResultGetAndGetMulti(t *testing.T) {
	t.Parallel()

	r := &Result{}
	r.add(AnswerEntry{Key: "name", Value: "alice"})
	r.add(AnswerEntry{Key: "tags", Multi: []string{"a", "b"}})

	if got := r.Get("name"); got != "alice" {
		t.Errorf("Get: want alice, got %q", got)
	}
	if got := r.Get("missing"); got != "" {
		t.Errorf("Get missing: want empty, got %q", got)
	}
	if got := r.GetMulti("tags"); len(got) != 2 || got[0] != "a" {
		t.Errorf("GetMulti: want [a b], got %v", got)
	}
	if got := r.GetMulti("missing"); got != nil {
		t.Errorf("GetMulti missing: want nil, got %v", got)
	}
}

func TestResultGetOnNilIndex(t *testing.T) {
	t.Parallel()
	r := &Result{}
	if got := r.Get("x"); got != "" {
		t.Errorf("empty Result.Get: want empty, got %q", got)
	}
	if got := r.GetMulti("x"); got != nil {
		t.Errorf("empty Result.GetMulti: want nil, got %v", got)
	}
}

func TestFlattenLinearChain(t *testing.T) {
	t.Parallel()

	// a → b → c
	flow := q.Text("A", "a").Then(
		q.Text("B", "b").Then(
			q.Text("C", "c").Q(),
		).Q(),
	).Q()

	var nodes []flatNode
	var loops []loopEntry
	if err := flatten(flow, nil, &nodes, &loops); err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 3 {
		t.Fatalf("want 3 nodes, got %d", len(nodes))
	}
	for i, want := range []string{"a", "b", "c"} {
		if got := nodeKey(nodes[i].question); got != want {
			t.Errorf("node[%d]: want %q, got %q", i, want, got)
		}
	}
}

func TestFlattenDuplicateKeyError(t *testing.T) {
	t.Parallel()

	flow := q.Text("A", "dup").Then(
		q.Text("B", "dup").Q(),
	).Q()

	var nodes []flatNode
	var loops []loopEntry
	err := flatten(flow, nil, &nodes, &loops)
	if err == nil {
		t.Fatal("want error for duplicate key, got nil")
	}
}

func TestCondsMatch(t *testing.T) {
	t.Parallel()

	runner := &Runner{Result: &Result{}}
	runner.Result.add(AnswerEntry{Key: "lang", Value: "go"})

	tests := []struct {
		conds []conditionType
		want  bool
	}{
		{[]conditionType{{key: "lang", value: "go"}}, true},
		{[]conditionType{{key: "lang", value: "ts"}}, false},
		{[]conditionType{{key: "lang", value: "go", negate: true}}, false},
		{[]conditionType{{key: "lang", value: "ts", negate: true}}, true},
		{nil, true},
	}
	for _, tc := range tests {
		if got := runner.condsMatch(tc.conds); got != tc.want {
			t.Errorf("conds=%v: want %v, got %v", tc.conds, tc.want, got)
		}
	}
}

func TestWithCond(t *testing.T) {
	t.Parallel()

	base := []conditionType{{key: "a", value: "1"}}
	result := withCond(base, conditionType{key: "b", value: "2"})
	if len(result) != 2 {
		t.Fatalf("want 2 conds, got %d", len(result))
	}
	// original must not be mutated
	if len(base) != 1 {
		t.Error("withCond mutated original slice")
	}
}
