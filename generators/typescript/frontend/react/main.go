package typescript_frontend_react_generator

import (
	"embed"
	"fmt"

	frontend_architecture_generator "github.com/version14/dot/generators/common/frontend/architecture"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	newtest_generator "github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

// all: is required to include dot-files (.gitkeep) in the embedded FS.
//
//go:embed all:files
var reactFiles embed.FS

const generatorName = "react-ts"

// ReactTS is the generator for React + TypeScript + Vite projects.
// Wire it into a question with ReactTS.Func():
//
//	question.Select(...).ChoiceWithGen("React", "react", ReactTS.Func())
var Generator = &newtest_generator.Generator{
	Name:     generatorName,
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		arch, _ := s.Extensions["frontend-architecture"].(string)
		data := struct {
			Name         string
			Architecture string
		}{s.Project.Name, arch}

		ops, err := genfs.RenderDir(reactFiles, "files", generatorName, data)
		if err != nil {
			return nil, fmt.Errorf("react-ts base: %w", err)
		}

		archOps, err := frontend_architecture_generator.ArchitectureTS.Apply(s)
		if err != nil {
			return nil, fmt.Errorf("react-ts arch: %w", err)
		}

		return append(ops, archOps...), nil
	},

	PostApplyFunction: func(s spec.Spec) []generator.PostOp {
		return []generator.PostOp{
			{Command: "pnpm", Args: []string{"install"}, Dir: ".", Generator: generatorName},
		}
	},
}
