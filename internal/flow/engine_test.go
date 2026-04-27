package flow

import (
	"errors"
	"reflect"
	"testing"
)

// scriptedAdapter returns answers from a map keyed by question ID.
type scriptedAdapter struct {
	answers map[string]Answer
	asked   []string
	err     error
}

func (s *scriptedAdapter) Ask(q Question, ctx *FlowContext) (Answer, error) {
	if s.err != nil {
		return nil, s.err
	}
	s.asked = append(s.asked, q.ID())
	if v, ok := s.answers[q.ID()]; ok {
		return v, nil
	}
	return "", nil
}

func TestEngine_LinearFlow(t *testing.T) {
	last := &TextQuestion{
		QuestionBase: QuestionBase{ID_: "name", Next_: &Next{End: true}},
		Label:        "name?",
	}
	first := &OptionQuestion{
		QuestionBase: QuestionBase{ID_: "type"},
		Label:        "type?",
		Options: []*Option{
			{Label: "Mono", Value: "mono", Next: &Next{Question: last}},
		},
	}

	ad := &scriptedAdapter{answers: map[string]Answer{"type": "mono", "name": "alpha"}}
	eng := NewEngine(ad)
	ctx, err := eng.Run(first)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got, want := ctx.Answers["type"], "mono"; got != want {
		t.Errorf("type = %v, want %v", got, want)
	}
	if got, want := ctx.Answers["name"], "alpha"; got != want {
		t.Errorf("name = %v, want %v", got, want)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"type", "name"}) {
		t.Errorf("visited = %v", ctx.VisitedNodes)
	}
}

func TestEngine_BranchingByOption(t *testing.T) {
	branchA := &TextQuestion{QuestionBase: QuestionBase{ID_: "a", Next_: &Next{End: true}}}
	branchB := &TextQuestion{QuestionBase: QuestionBase{ID_: "b", Next_: &Next{End: true}}}
	root := &OptionQuestion{
		QuestionBase: QuestionBase{ID_: "pick"},
		Options: []*Option{
			{Value: "a", Next: &Next{Question: branchA}},
			{Value: "b", Next: &Next{Question: branchB}},
		},
	}

	ad := &scriptedAdapter{answers: map[string]Answer{"pick": "b", "b": "ok"}}
	ctx, err := NewEngine(ad).Run(root)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"pick", "b"}) {
		t.Errorf("visited = %v, want [pick b]", ctx.VisitedNodes)
	}
	if _, ok := ctx.Answers["a"]; ok {
		t.Errorf("unvisited branch a should not appear in answers")
	}
}

func TestEngine_IfQuestionRoutesWithoutAsking(t *testing.T) {
	thenQ := &TextQuestion{QuestionBase: QuestionBase{ID_: "yes", Next_: &Next{End: true}}}
	elseQ := &TextQuestion{QuestionBase: QuestionBase{ID_: "no", Next_: &Next{End: true}}}
	gate := &IfQuestion{
		QuestionBase: QuestionBase{ID_: "gate"},
		Condition:    func(ctx *FlowContext) bool { return ctx.Answers["enabled"] == true },
		Then:         &Next{Question: thenQ},
		Else:         &Next{Question: elseQ},
	}

	ad := &scriptedAdapter{answers: map[string]Answer{"yes": "y"}}
	eng := NewEngine(ad)
	ctx, err := eng.Run(gate)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// gate has no answer recorded but is visited; "enabled" was never set,
	// so the condition is false and Else routes to "no".
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"gate", "no"}) {
		t.Errorf("visited = %v, want [gate no]", ctx.VisitedNodes)
	}
	for _, asked := range ad.asked {
		if asked == "gate" {
			t.Errorf("IfQuestion should not be asked")
		}
	}
}

func TestEngine_FragmentResolution(t *testing.T) {
	tail := &TextQuestion{QuestionBase: QuestionBase{ID_: "tail", Next_: &Next{End: true}}}
	root := &OptionQuestion{
		QuestionBase: QuestionBase{ID_: "linter"},
		Options: []*Option{
			{Value: "biome", Next: &Next{Fragment: "after-linter"}},
		},
	}

	ad := &scriptedAdapter{answers: map[string]Answer{"linter": "biome", "tail": "x"}}
	eng := NewEngine(ad)
	_ = eng.RegisterFragment(&FlowFragment{
		ID: "after-linter",
		Resolve: func(ctx *FlowContext) *Next {
			return &Next{Question: tail}
		},
	})

	ctx, err := eng.Run(root)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"linter", "tail"}) {
		t.Errorf("visited = %v", ctx.VisitedNodes)
	}
}

func TestEngine_Inject_Replace_SwapsTargetAndDiscardsOriginalChain(t *testing.T) {
	originalTail := &TextQuestion{QuestionBase: QuestionBase{ID_: "original_tail", Next_: &Next{End: true}}}
	original := &TextQuestion{
		QuestionBase: QuestionBase{ID_: "linter", Next_: &Next{Question: originalTail}},
	}

	pluginTail := &TextQuestion{QuestionBase: QuestionBase{ID_: "biome.tail", Next_: &Next{End: true}}}
	replacement := &TextQuestion{
		QuestionBase: QuestionBase{ID_: "biome.linter", Next_: &Next{Question: pluginTail}},
	}

	ad := &scriptedAdapter{answers: map[string]Answer{"biome.linter": "x", "biome.tail": "y"}}
	eng := NewEngine(ad)
	if err := eng.RegisterPlugin("biome"); err != nil {
		t.Fatalf("RegisterPlugin: %v", err)
	}
	if err := eng.Inject(&Injection{
		Plugin:      "biome",
		TargetID:    "linter",
		Kind:        InjectReplace,
		Replacement: replacement,
	}); err != nil {
		t.Fatalf("Inject: %v", err)
	}

	ctx, err := eng.Run(original)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"biome.linter", "biome.tail"}) {
		t.Errorf("visited = %v", ctx.VisitedNodes)
	}
	if _, ok := ctx.Answers["original_tail"]; ok {
		t.Errorf("original chain should be discarded by Replace")
	}
}

func TestEngine_Inject_AddOption_AppendsAndCanRoute(t *testing.T) {
	pluginBranch := &TextQuestion{QuestionBase: QuestionBase{ID_: "biome.config", Next_: &Next{End: true}}}
	root := &OptionQuestion{
		QuestionBase: QuestionBase{ID_: "linter", Next_: &Next{End: true}},
		Options: []*Option{
			{Label: "Prettier", Value: "prettier", Next: &Next{End: true}},
		},
	}

	ad := &scriptedAdapter{answers: map[string]Answer{"linter": "biome.biome", "biome.config": "strict"}}
	eng := NewEngine(ad)
	_ = eng.RegisterPlugin("biome")
	if err := eng.Inject(&Injection{
		Plugin:   "biome",
		TargetID: "linter",
		Kind:     InjectAddOption,
		Option:   &Option{Label: "Biome", Value: "biome.biome", Next: &Next{Question: pluginBranch}},
	}); err != nil {
		t.Fatalf("Inject: %v", err)
	}

	ctx, err := eng.Run(root)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"linter", "biome.config"}) {
		t.Errorf("visited = %v", ctx.VisitedNodes)
	}
}

func TestEngine_Inject_InsertAfter_PreservesOriginalNext(t *testing.T) {
	tail := &TextQuestion{QuestionBase: QuestionBase{ID_: "tail", Next_: &Next{End: true}}}
	root := &TextQuestion{
		QuestionBase: QuestionBase{ID_: "linter", Next_: &Next{Question: tail}},
	}

	plugQ := &TextQuestion{QuestionBase: QuestionBase{ID_: "biome.strict", Next_: &Next{End: true}}}

	ad := &scriptedAdapter{answers: map[string]Answer{"linter": "x", "biome.strict": "yes", "tail": "done"}}
	eng := NewEngine(ad)
	_ = eng.RegisterPlugin("biome")
	if err := eng.Inject(&Injection{
		Plugin:   "biome",
		TargetID: "linter",
		Kind:     InjectInsertAfter,
		Question: plugQ,
	}); err != nil {
		t.Fatalf("Inject: %v", err)
	}

	ctx, err := eng.Run(root)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"linter", "biome.strict", "tail"}) {
		t.Errorf("visited = %v, want [linter biome.strict tail]", ctx.VisitedNodes)
	}
}

func TestEngine_Inject_InsertAfter_ChainsMultipleInOrder(t *testing.T) {
	tail := &TextQuestion{QuestionBase: QuestionBase{ID_: "tail", Next_: &Next{End: true}}}
	root := &TextQuestion{
		QuestionBase: QuestionBase{ID_: "core", Next_: &Next{Question: tail}},
	}
	q1 := &TextQuestion{QuestionBase: QuestionBase{ID_: "p.q1", Next_: &Next{End: true}}}
	q2 := &TextQuestion{QuestionBase: QuestionBase{ID_: "p.q2", Next_: &Next{End: true}}}

	ad := &scriptedAdapter{}
	eng := NewEngine(ad)
	_ = eng.RegisterPlugin("p")
	_ = eng.Inject(&Injection{Plugin: "p", TargetID: "core", Kind: InjectInsertAfter, Question: q1})
	_ = eng.Inject(&Injection{Plugin: "p", TargetID: "core", Kind: InjectInsertAfter, Question: q2})

	ctx, err := eng.Run(root)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !reflect.DeepEqual(ctx.VisitedNodes, []string{"core", "p.q1", "p.q2", "tail"}) {
		t.Errorf("visited = %v", ctx.VisitedNodes)
	}
}

func TestHookRegistry_RejectsUnregisteredPlugin(t *testing.T) {
	r := NewHookRegistry()
	err := r.Inject(&Injection{
		Plugin:      "ghost",
		TargetID:    "x",
		Kind:        InjectReplace,
		Replacement: &TextQuestion{QuestionBase: QuestionBase{ID_: "ghost.y"}},
	})
	if err == nil {
		t.Fatal("expected error for unregistered plugin")
	}
}

func TestHookRegistry_EnforcesPrefixOnContributedIDs(t *testing.T) {
	r := NewHookRegistry()
	_ = r.RegisterPlugin("biome")

	cases := []struct {
		name string
		inj  *Injection
	}{
		{
			name: "Replace without prefix",
			inj: &Injection{
				Plugin: "biome", TargetID: "linter", Kind: InjectReplace,
				Replacement: &TextQuestion{QuestionBase: QuestionBase{ID_: "naked"}},
			},
		},
		{
			name: "AddOption without prefix",
			inj: &Injection{
				Plugin: "biome", TargetID: "linter", Kind: InjectAddOption,
				Option: &Option{Value: "naked"},
			},
		},
		{
			name: "InsertAfter without prefix",
			inj: &Injection{
				Plugin: "biome", TargetID: "linter", Kind: InjectInsertAfter,
				Question: &TextQuestion{QuestionBase: QuestionBase{ID_: "naked"}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := r.Inject(tc.inj); err == nil {
				t.Errorf("expected prefix-enforcement error")
			}
		})
	}
}

func TestHookRegistry_RejectsDuplicatePlugin(t *testing.T) {
	r := NewHookRegistry()
	if err := r.RegisterPlugin("a"); err != nil {
		t.Fatal(err)
	}
	if err := r.RegisterPlugin("a"); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestEngine_AdapterErrorPropagates(t *testing.T) {
	root := &TextQuestion{QuestionBase: QuestionBase{ID_: "x", Next_: &Next{End: true}}}
	ad := &scriptedAdapter{err: errors.New("boom")}
	if _, err := NewEngine(ad).Run(root); err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_NilAdapterFails(t *testing.T) {
	root := &TextQuestion{QuestionBase: QuestionBase{ID_: "x", Next_: &Next{End: true}}}
	if _, err := NewEngine(nil).Run(root); err == nil {
		t.Fatal("expected error when no adapter is configured")
	}
}
