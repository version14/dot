package dotdir

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/version14/dot/internal/spec"
)

func TestSpecRoundTrip(t *testing.T) {
	root := t.TempDir()
	want := &spec.ProjectSpec{
		FlowID:    "monorepo",
		CreatedAt: time.Date(2026, 4, 25, 10, 30, 0, 0, time.UTC),
		Metadata:  spec.ProjectMetadata{ProjectName: "x", ToolVersion: "0.1.0"},
		Answers:   map[string]spec.AnswerNode{"linter": "biome"},
	}
	if err := SaveSpec(root, want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := LoadSpec(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.FlowID != want.FlowID || got.Metadata.ProjectName != "x" || got.Answers["linter"] != "biome" {
		t.Errorf("round-trip mismatch: %+v", got)
	}
}

func TestManifestRoundTrip(t *testing.T) {
	root := t.TempDir()
	want := &Manifest{
		ToolVersion:    "0.1.0",
		LastExecutedAt: time.Now().UTC().Truncate(time.Second),
		GeneratorsExecuted: []ExecutedGenerator{{
			Name:              "base-project",
			VersionConstraint: "^1.0.0",
			ResolvedVersion:   "1.2.5",
			InvocationCount:   1,
		}},
	}
	if err := SaveManifest(root, want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := LoadManifest(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(got.GeneratorsExecuted) != 1 || got.GeneratorsExecuted[0].Name != "base-project" {
		t.Errorf("round-trip mismatch: %+v", got)
	}
}

func TestWriteIgnore_HasExactSpecContents(t *testing.T) {
	root := t.TempDir()
	if err := WriteIgnore(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	want := IgnoreContents
	got, err := readFile(filepath.Join(root, DirName, IgnoreFile))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got != want {
		t.Errorf(".gitignore mismatch:\nwant=%q\n got=%q", want, got)
	}
}

func readFile(p string) (string, error) {
	b, err := readAll(p)
	return string(b), err
}
