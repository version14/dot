package frontend_architecture_generator

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/version14/dot/internal/generator"
	newtest_generator "github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var archFiles embed.FS

// ArchitectureTS creates the src/ folder structure for a given frontend
// architecture pattern. It is framework-agnostic: any TypeScript frontend
// generator (React, Next.js, etc.) can call ArchitectureTS.Apply(s) to get
// the right directories without duplicating the folder manifests.
//
// Reads spec.Extensions["architecture"]: "feature-sliced" | "atomic" | "container-presentational".
var ArchitectureTS = &newtest_generator.Generator{
	Name:     "frontend-architecture-ts",
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		arch, _ := s.Extensions["architecture"].(string)
		if arch == "" {
			return nil, fmt.Errorf("frontend-architecture-ts: spec.Extensions[\"architecture\"] is not set")
		}

		archDir := fmt.Sprintf("files/%s", arch)
		ops, err := walkDir(archFiles, archDir)
		if err != nil {
			return nil, fmt.Errorf("frontend-architecture-ts [%s]: %w", arch, err)
		}
		return ops, nil
	},
}

// walkDir returns a FileOp for every file found under root in fsys.
// Output paths are relative to root (root prefix is stripped).
func walkDir(fsys embed.FS, root string) ([]generator.FileOp, error) {
	var ops []generator.FileOp

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		outPath := strings.TrimPrefix(path, root+"/")
		ops = append(ops, generator.FileOp{
			Kind:      generator.Create,
			Path:      outPath,
			Content:   string(content),
			Generator: "frontend-architecture-ts",
			Priority:  100,
		})
		return nil
	})

	return ops, err
}
