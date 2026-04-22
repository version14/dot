package frontend_architecture_generator

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	newtest_generator "github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var archFiles embed.FS

const generatorName = "frontend-architecture-ts"

// ArchitectureTS creates the src/ folder structure for a given frontend
// architecture pattern. It is framework-agnostic: any TypeScript frontend
// generator (React, Next.js, etc.) can call ArchitectureTS.Apply(s) to get
// the right directories without duplicating the folder manifests.
//
// Reads spec.Extensions["architecture"]: "feature-sliced" | "atomic" | "container-presentational".
var ArchitectureTS = &newtest_generator.Generator{
	Name:     generatorName,
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		arch, _ := s.Extensions["architecture"].(string)
		if arch == "" {
			return nil, fmt.Errorf("frontend-architecture-ts: spec.Extensions[\"architecture\"] is not set")
		}

		ops, err := genfs.WalkDir(archFiles, "files/"+arch, generatorName)
		if err != nil {
			return nil, fmt.Errorf("frontend-architecture-ts [%s]: %w", arch, err)
		}
		return ops, nil
	},
}

// walkDir is kept for tests that call it directly.
func walkDir(fsys embed.FS, root string) ([]generator.FileOp, error) {
	return genfs.WalkDir(fsys, root, generatorName)
}
