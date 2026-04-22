package generators_typescript_backend_architecture_clean

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var cleanFiles embed.FS

const generatorName = "typescript-backend-clean-architecture"

// CleanArchitectureTS scaffolds a TypeScript Express clean-architecture project structure.
var CleanArchitectureTS = &scaffold.Generator{
	Name:     generatorName,
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		data := struct{ Name string }{s.Project.Name}

		ops, err := genfs.RenderDir(cleanFiles, "files", generatorName, data)
		if err != nil {
			return nil, fmt.Errorf("typescript-backend-clean-architecture: %w", err)
		}
		return ops, nil
	},

	PostApplyFunction: func(s spec.Spec) []generator.PostOp {
		return []generator.PostOp{
			{Command: "pnpm", Args: []string{"install"}, Dir: ".", Generator: generatorName},
		}
	},
}
