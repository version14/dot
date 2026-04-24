package typescript_base_generator

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var baseFiles embed.FS

const generatorName = "typescript-base"

// BasePriority is lower than framework generators (100) so that React, MVC,
// clean-architecture, etc. can override package.json and tsconfig.json when
// they supply their own version.
const BasePriority = 50

// BaseTypescriptTS generates the foundational TypeScript project files:
// package.json, tsconfig.json, .gitignore, and src/index.ts.
//
// It runs when the user picks TypeScript as language, before any
// framework-specific generator. Framework generators at priority 100 override
// any conflicting files (e.g. React overrides package.json and tsconfig.json).
var BaseTypescriptTS = &scaffold.Generator{
	Name:     generatorName,
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		data := struct{ Name string }{s.Project.Name}

		ops, err := genfs.RenderDir(baseFiles, "files", generatorName, data)
		if err != nil {
			return nil, fmt.Errorf("typescript-base: %w", err)
		}

		for i := range ops {
			ops[i].Priority = BasePriority
		}
		return ops, nil
	},
}
