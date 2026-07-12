package fingerprint

import (
	"testing"

	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/state"
)

// run drives the Collector the way the Executor does: BeforeGenerate, the
// generator's writes, AfterGenerate — sharing one root state across generators,
// exactly as a real scaffold does.
func run(t *testing.T, steps []step) map[string]string {
	t.Helper()
	root := state.NewVirtualProjectState(spec.ProjectMetadata{ProjectName: "p"})
	c := NewCollector()
	for _, s := range steps {
		root.SetCurrentGenerator(s.name)
		c.BeforeGenerate(s.name, root)
		s.write(root)
		c.AfterGenerate(s.name, root)
	}
	return c.Fingerprints()
}

type step struct {
	name  string
	write func(*state.VirtualProjectState)
}

func mergeDeps(deps map[string]interface{}) func(*state.VirtualProjectState) {
	return func(s *state.VirtualProjectState) {
		_ = s.UpdateJSON("package.json", func(d *state.JSONDoc) error {
			d.Merge(map[string]interface{}{"dependencies": deps})
			return nil
		})
	}
}

// A generator's Contribution must depend only on what *it* did — never on what
// another generator merged into the same shared file beforehand.
//
// This is the whole reason a Contribution is a diff rather than a file hash. If
// it were a file hash, bumping react would move the Fingerprint of every
// generator that touches package.json, and a react bump would demand a version
// bump on all of them. That is exactly the defect that produced commit d53c302:
// 15 generator manifests bumped for a change to one generator.
func TestContributionIsIndependentOfOtherGenerators(t *testing.T) {
	// arkui merges its own dep into package.json *after* reactapp merged react.
	withReact1 := run(t, []step{
		{"reactapp", mergeDeps(map[string]interface{}{"react": "^19.2.7"})},
		{"arkui", mergeDeps(map[string]interface{}{"@ark-ui/react": "^5.37.2"})},
	})

	// Now react's Pin moves. arkui did not change, and does not name react.
	withReact2 := run(t, []step{
		{"reactapp", mergeDeps(map[string]interface{}{"react": "^19.9.9"})},
		{"arkui", mergeDeps(map[string]interface{}{"@ark-ui/react": "^5.37.2"})},
	})

	if withReact1["reactapp"] == withReact2["reactapp"] {
		t.Error("reactapp names react: its Fingerprint MUST move when the react Pin moves")
	}
	if withReact1["arkui"] != withReact2["arkui"] {
		t.Error("arkui does not name react: its Fingerprint must NOT move when the react Pin moves " +
			"(this is the d53c302 over-triggering bug)")
	}
}

// Ordering must not leak either: a generator that runs second contributes the
// same edit it would have contributed running first.
func TestContributionIsIndependentOfOrder(t *testing.T) {
	a := run(t, []step{
		{"first", mergeDeps(map[string]interface{}{"a": "^1.0.0"})},
		{"second", mergeDeps(map[string]interface{}{"b": "^2.0.0"})},
	})
	b := run(t, []step{
		{"second", mergeDeps(map[string]interface{}{"b": "^2.0.0"})},
		{"first", mergeDeps(map[string]interface{}{"a": "^1.0.0"})},
	})
	if a["second"] != b["second"] {
		t.Error("a generator's Contribution must not depend on what ran before it")
	}
}

// A generator that writes nothing has no Contribution and no Fingerprint.
func TestNoWriteNoFingerprint(t *testing.T) {
	fps := run(t, []step{
		{"silent", func(*state.VirtualProjectState) {}},
	})
	if _, ok := fps["silent"]; ok {
		t.Error("a generator that contributes nothing must not get a Fingerprint")
	}
}

// Appending to a shared file (e.g. .env.example) contributes only the appended
// bytes — not the whole file, which would again couple generators together.
func TestAppendContributesOnlyTheTail(t *testing.T) {
	appendEnv := func(text string) func(*state.VirtualProjectState) {
		return func(s *state.VirtualProjectState) { s.AppendFile(".env.example", []byte(text)) }
	}

	first := run(t, []step{
		{"a", appendEnv("A=1\n")},
		{"b", appendEnv("B=2\n")},
	})
	// "a" now writes something different; "b" appends the same line as before.
	second := run(t, []step{
		{"a", appendEnv("A=999\n")},
		{"b", appendEnv("B=2\n")},
	})

	if first["a"] == second["a"] {
		t.Error("a changed what it appends: its Fingerprint must move")
	}
	if first["b"] != second["b"] {
		t.Error("b appends the same bytes regardless of what a wrote: its Fingerprint must not move")
	}
}

// Creating a file contributes its content.
func TestCreateContributesContent(t *testing.T) {
	write := func(body string) func(*state.VirtualProjectState) {
		return func(s *state.VirtualProjectState) {
			s.WriteFile("src/main.ts", []byte(body), state.ContentRaw)
		}
	}
	a := run(t, []step{{"g", write("console.log(1)")}})
	b := run(t, []step{{"g", write("console.log(2)")}})
	if a["g"] == b["g"] {
		t.Error("changing a created file's content must move the Fingerprint")
	}

	same := run(t, []step{{"g", write("console.log(1)")}})
	if a["g"] != same["g"] {
		t.Error("identical output must produce an identical Fingerprint")
	}
}
