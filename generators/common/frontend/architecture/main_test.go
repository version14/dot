package frontend_architecture_generator

import (
	"strings"
	"testing"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/spec"
)

func TestArchitectureTSKnownPatterns(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		arch     string
		prefixes []string // expected path prefixes in output
	}{
		{
			arch:     "feature-sliced",
			prefixes: []string{"src/app", "src/features", "src/pages", "src/widgets", "src/entities", "src/shared"},
		},
		{
			arch:     "atomic",
			prefixes: []string{"src/components/atoms", "src/components/molecules", "src/components/organisms"},
		},
		{
			arch:     "container-presentational",
			prefixes: []string{"src/components", "src/containers", "src/pages"},
		},
	}

	for _, tc := range patterns {
		tc := tc
		t.Run(tc.arch, func(t *testing.T) {
			t.Parallel()

			s := spec.Spec{Extensions: map[string]any{"architecture": tc.arch}}
			ops, err := ArchitectureTS.Apply(s)
			if err != nil {
				t.Fatalf("Apply(%q): %v", tc.arch, err)
			}
			if len(ops) == 0 {
				t.Fatalf("Apply(%q): no FileOps returned", tc.arch)
			}

			for _, op := range ops {
				if op.Kind != generator.Create {
					t.Errorf("op.Kind: want Create, got %v (path=%s)", op.Kind, op.Path)
				}
				if op.Generator != "frontend-architecture-ts" {
					t.Errorf("op.Generator: want frontend-architecture-ts, got %q", op.Generator)
				}
			}

			for _, prefix := range tc.prefixes {
				found := false
				for _, op := range ops {
					if strings.HasPrefix(op.Path, prefix) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("no FileOp with prefix %q found in ops for arch=%q", prefix, tc.arch)
				}
			}
		})
	}
}

func TestArchitectureTSMissingExtension(t *testing.T) {
	t.Parallel()

	s := spec.Spec{}
	_, err := ArchitectureTS.Apply(s)
	if err == nil {
		t.Fatal("want error when architecture extension is missing, got nil")
	}
}

func TestArchitectureTSUnknownArch(t *testing.T) {
	t.Parallel()

	s := spec.Spec{Extensions: map[string]any{"architecture": "nonexistent-pattern"}}
	_, err := ArchitectureTS.Apply(s)
	if err == nil {
		t.Fatal("want error for unknown architecture pattern, got nil")
	}
}

func TestWalkDirPathsAreRelative(t *testing.T) {
	t.Parallel()

	ops, err := walkDir(archFiles, "files/atomic")
	if err != nil {
		t.Fatalf("walkDir: %v", err)
	}
	for _, op := range ops {
		if strings.HasPrefix(op.Path, "files/") {
			t.Errorf("path should not contain 'files/' prefix: %q", op.Path)
		}
	}
}
