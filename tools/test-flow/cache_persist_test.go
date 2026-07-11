package main

import (
	"os"
	"path/filepath"
	"testing"
)

// write creates path (and its parent dirs) under root with the given
// content, failing the test on any error.
func write(t *testing.T, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
}

// newCoreRepo builds a minimal synthetic repo tree exercising every
// coreExcludeDirs entry plus one "real" core file, and returns its root.
func newCoreRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write(t, root, "internal/cli/runner.go", "package cli")
	write(t, root, "README.md", "# hello")
	write(t, root, "docs/guide.md", "# guide")
	write(t, root, "bin/dot", "binary")
	write(t, root, "generators/react_app/generator.go", "package reactapp")
	write(t, root, "tools/test-flow/main.go", "package main")
	write(t, root, "tools/test-flow/testdata/some_case.json", `{"name":"x"}`)
	write(t, root, ".test-flow-cache/some_case.json", `{"fingerprint":"x"}`)
	write(t, root, ".git/HEAD", "ref: refs/heads/main")
	return root
}

func TestHashCore_ExcludesDenylistedDirs(t *testing.T) {
	root := newCoreRepo(t)
	before, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore: %v", err)
	}

	// Editing anything under an excluded dir must NOT change the hash.
	excludedEdits := map[string]string{
		"docs/guide.md":                           "# guide, but different now",
		"bin/dot":                                 "a completely different binary",
		"generators/react_app/generator.go":       "package reactapp\n\n// edited",
		"tools/test-flow/testdata/some_case.json": `{"name":"y"}`,
		".test-flow-cache/some_case.json":         `{"fingerprint":"y"}`,
		".git/HEAD":                               "ref: refs/heads/other",
		"README.md":                               "# hello, but different now",
		".claude/skills/foo/SKILL.md":             "# foo skill, but different now",
		".zed/settings.json":                      `{"theme":"dark"}`,
		".githooks/pre-commit":                    "#!/bin/sh\necho edited",
	}
	for path, content := range excludedEdits {
		write(t, root, path, content)
		after, err := hashCore(root)
		if err != nil {
			t.Fatalf("hashCore after editing %s: %v", path, err)
		}
		if after != before {
			t.Errorf("editing excluded path %q changed the core hash; expected it to be ignored", path)
		}
	}
}

func TestHashCore_MarkdownExcludedAnywhereInCore(t *testing.T) {
	root := newCoreRepo(t)
	before, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore: %v", err)
	}

	// A new .md file outside any excluded dir (e.g. a root-level doc) must
	// not affect the core hash either.
	write(t, root, "CONTRIBUTING.md", "# contributing")
	after, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore after adding markdown: %v", err)
	}
	if after != before {
		t.Errorf("adding a .md file changed the core hash; markdown should never affect it")
	}
}

func TestHashCore_DetectsRealCoreChanges(t *testing.T) {
	root := newCoreRepo(t)
	before, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore: %v", err)
	}

	// A change to a genuine core file (engine code, not under any excluded
	// dir, not markdown) MUST change the hash.
	write(t, root, "internal/cli/runner.go", "package cli\n\n// behaviour changed")
	after, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore after editing core file: %v", err)
	}
	if after == before {
		t.Errorf("editing internal/cli/runner.go did not change the core hash; core changes must invalidate every fixture")
	}
}

func TestHashCore_DetectsTestFlowToolChangesOutsideTestdata(t *testing.T) {
	root := newCoreRepo(t)
	before, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore: %v", err)
	}

	// tools/test-flow/testdata/ is excluded, but the rest of tools/test-flow
	// (the runner itself) is not — this is the nested-exclusion case.
	write(t, root, "tools/test-flow/main.go", "package main\n\n// behaviour changed")
	after, err := hashCore(root)
	if err != nil {
		t.Fatalf("hashCore after editing tools/test-flow/main.go: %v", err)
	}
	if after == before {
		t.Errorf("editing tools/test-flow/main.go did not change the core hash; only testdata/ should be excluded, not the whole tool")
	}
}
