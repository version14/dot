package generator

import (
	"errors"
	"testing"

	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type fakeGen struct {
	name   string
	called int
	answer map[string]interface{}
	err    error
}

func (f *fakeGen) Name() string    { return f.name }
func (f *fakeGen) Version() string { return "1.0.0" }
func (f *fakeGen) Generate(ctx *dotapi.Context) error {
	if f.err != nil {
		return f.err
	}
	f.called++
	f.answer = ctx.Answers
	return ctx.State.CreateFile("touched-by-"+f.name, []byte("ok"))
}

func TestExecutor_RunsGeneratorsInOrder(t *testing.T) {
	reg := NewRegistry()
	a := &fakeGen{name: "a"}
	b := &fakeGen{name: "b"}
	if err := reg.Register(Manifest{Name: "a"}, a); err != nil {
		t.Fatalf("register a: %v", err)
	}
	if err := reg.Register(Manifest{Name: "b"}, b); err != nil {
		t.Fatalf("register b: %v", err)
	}

	s := &spec.ProjectSpec{Answers: map[string]spec.AnswerNode{"project_name": "x"}}
	vstate := state.NewVirtualProjectState(spec.ProjectMetadata{})
	exec := NewExecutor(reg, dotapi.DiscardLogger{})

	err := exec.Execute([]Invocation{{Name: "a"}, {Name: "b"}}, s, vstate)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if a.called != 1 || b.called != 1 {
		t.Errorf("call counts a=%d b=%d", a.called, b.called)
	}
	if !vstate.FileExists("touched-by-a") || !vstate.FileExists("touched-by-b") {
		t.Errorf("missing generator outputs")
	}
}

func TestExecutor_ScopesAnswersPerLoopFrame(t *testing.T) {
	reg := NewRegistry()
	g := &fakeGen{name: "svc"}
	_ = reg.Register(Manifest{Name: "svc"}, g)

	s := &spec.ProjectSpec{Answers: map[string]spec.AnswerNode{"project_name": "p"}}
	vstate := state.NewVirtualProjectState(spec.ProjectMetadata{})
	exec := NewExecutor(reg, dotapi.DiscardLogger{})

	frame := []flow.LoopFrame{{
		QuestionID: "services",
		Index:      0,
		Answers:    map[string]flow.AnswerNode{"name": "auth"},
	}}
	if err := exec.Execute([]Invocation{{Name: "svc", LoopStack: frame}}, s, vstate); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got := g.answer["name"]; got != "auth" {
		t.Errorf("scoped name = %v, want auth", got)
	}
	if got := g.answer["project_name"]; got != "p" {
		t.Errorf("global project_name = %v, want p", got)
	}
}

func TestExecutor_HaltsOnGeneratorError(t *testing.T) {
	reg := NewRegistry()
	good := &fakeGen{name: "good"}
	bad := &fakeGen{name: "bad", err: errors.New("explode")}
	never := &fakeGen{name: "never"}
	_ = reg.Register(Manifest{Name: "good"}, good)
	_ = reg.Register(Manifest{Name: "bad"}, bad)
	_ = reg.Register(Manifest{Name: "never"}, never)

	s := &spec.ProjectSpec{}
	vstate := state.NewVirtualProjectState(spec.ProjectMetadata{})
	err := NewExecutor(reg, dotapi.DiscardLogger{}).Execute(
		[]Invocation{{Name: "good"}, {Name: "bad"}, {Name: "never"}},
		s, vstate,
	)
	if err == nil {
		t.Fatal("expected error")
	}
	var failed *ErrGeneratorFailed
	if !errors.As(err, &failed) {
		t.Errorf("error type = %T, want *ErrGeneratorFailed", err)
	}
	if never.called != 0 {
		t.Errorf("expected never to be skipped after failure")
	}
}

func TestExecutor_UnknownGenerator(t *testing.T) {
	reg := NewRegistry()
	s := &spec.ProjectSpec{}
	vstate := state.NewVirtualProjectState(spec.ProjectMetadata{})
	err := NewExecutor(reg, nil).Execute([]Invocation{{Name: "ghost"}}, s, vstate)
	var unknown *ErrUnknownGenerator
	if !errors.As(err, &unknown) {
		t.Errorf("error = %v, want ErrUnknownGenerator", err)
	}
}
