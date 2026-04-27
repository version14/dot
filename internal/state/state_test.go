package state

import (
	"strings"
	"testing"

	"github.com/version14/dot/internal/spec"
)

func newState(t *testing.T) *VirtualProjectState {
	t.Helper()
	return NewVirtualProjectState(spec.ProjectMetadata{ProjectName: "test"})
}

func TestVirtualProjectState_CreateGetExists(t *testing.T) {
	s := newState(t)
	s.SetCurrentGenerator("base")

	if err := s.CreateFile("README.md", []byte("hi")); err != nil {
		t.Fatalf("create: %v", err)
	}
	if !s.FileExists("README.md") {
		t.Error("FileExists = false, want true")
	}
	got, ok := s.GetFile("README.md")
	if !ok || string(got.Content) != "hi" {
		t.Errorf("GetFile = %v ok=%v", got, ok)
	}
	if got.CreatedBy != "base" {
		t.Errorf("CreatedBy = %q, want base", got.CreatedBy)
	}
	if err := s.CreateFile("README.md", []byte("dup")); err == nil {
		t.Error("expected duplicate-create error")
	}
}

func TestVirtualProjectState_DeleteAndPaths(t *testing.T) {
	s := newState(t)
	_ = s.CreateFile("a.txt", []byte("a"))
	_ = s.CreateFile("b.txt", []byte("b"))
	s.DeleteFile("a.txt")
	if s.FileExists("a.txt") {
		t.Error("expected a.txt deleted")
	}
	paths := s.Paths()
	if len(paths) != 1 || paths[0] != "b.txt" {
		t.Errorf("Paths = %v", paths)
	}
}

func TestUpdateJSON_CreateAndModify(t *testing.T) {
	s := newState(t)
	if err := s.UpdateJSON("package.json", func(d *JSONDoc) error {
		return d.SetNested("name", "alpha")
	}); err != nil {
		t.Fatalf("first update: %v", err)
	}
	if err := s.UpdateJSON("package.json", func(d *JSONDoc) error {
		return d.AddDep("dependencies", "react", "^18.0.0")
	}); err != nil {
		t.Fatalf("second update: %v", err)
	}
	node, _ := s.GetFile("package.json")
	if !strings.Contains(string(node.Content), `"react": "^18.0.0"`) {
		t.Errorf("missing react dep:\n%s", node.Content)
	}
	if !strings.Contains(string(node.Content), `"name": "alpha"`) {
		t.Errorf("missing name:\n%s", node.Content)
	}
}

func TestJSONDoc_SetNestedAndDelete(t *testing.T) {
	d := NewJSONDoc()
	if err := d.SetNested("scripts.build", "tsc"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if v, ok := d.GetNested("scripts.build"); !ok || v != "tsc" {
		t.Errorf("get = %v ok=%v", v, ok)
	}
	d.DeleteKey("scripts.build")
	if _, ok := d.GetNested("scripts.build"); ok {
		t.Error("expected key deleted")
	}
}

func TestJSONDoc_Merge(t *testing.T) {
	d := NewJSONDoc()
	_ = d.SetNested("dependencies.react", "^18")
	d.Merge(map[string]interface{}{
		"dependencies": map[string]interface{}{"axios": "^1"},
		"name":         "x",
	})
	if v, _ := d.GetNested("dependencies.react"); v != "^18" {
		t.Errorf("react lost during merge: %v", v)
	}
	if v, _ := d.GetNested("dependencies.axios"); v != "^1" {
		t.Errorf("axios not merged: %v", v)
	}
	if v, _ := d.GetNested("name"); v != "x" {
		t.Errorf("name not set")
	}
}

func TestUpdateGoMod(t *testing.T) {
	s := newState(t)
	if err := s.UpdateGoMod(func(m *GoMod) error {
		m.AddModule("example.com/x")
		m.SetGoVersion("1.26")
		m.AddRequire("github.com/lib/pq", "v1.10.0")
		return nil
	}); err != nil {
		t.Fatalf("update: %v", err)
	}
	node, _ := s.GetFile("go.mod")
	out := string(node.Content)
	if !strings.Contains(out, "module example.com/x") {
		t.Errorf("missing module:\n%s", out)
	}
	if !strings.Contains(out, "go 1.26") {
		t.Errorf("missing go version:\n%s", out)
	}
	if !strings.Contains(out, "github.com/lib/pq v1.10.0") {
		t.Errorf("missing require:\n%s", out)
	}

	if err := s.UpdateGoMod(func(m *GoMod) error {
		m.RemoveRequire("github.com/lib/pq")
		return nil
	}); err != nil {
		t.Fatalf("remove: %v", err)
	}
	node, _ = s.GetFile("go.mod")
	if strings.Contains(string(node.Content), "lib/pq") {
		t.Errorf("require not removed:\n%s", node.Content)
	}
}

func TestUpdateYAML(t *testing.T) {
	s := newState(t)
	if err := s.UpdateYAML("docker-compose.yml", func(d *YAMLDoc) error {
		d.SetKey("version", "3.8")
		return d.Append("services", map[string]interface{}{"name": "auth"})
	}); err != nil {
		t.Fatalf("update: %v", err)
	}
	node, _ := s.GetFile("docker-compose.yml")
	out := string(node.Content)
	if !strings.Contains(out, "version: \"3.8\"") && !strings.Contains(out, "version: 3.8") {
		t.Errorf("missing version:\n%s", out)
	}
	if !strings.Contains(out, "services:") {
		t.Errorf("missing services key:\n%s", out)
	}
}
