package typescript_linters_biome_generator

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var expressFiles embed.FS

const generatorName = "biome-base-config"

var Generator = &scaffold.Generator{
	Name:     generatorName,
	Version:  "1.0.0",
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		ops, err := genfs.RenderDir(expressFiles, "files", generatorName, "")
		if err != nil {
			return nil, fmt.Errorf("biome-base-config : %w", err)
		}

		return ops, nil
	},

	PostApplyFunction: func(s spec.Spec) []generator.PostOp {
		return []generator.PostOp{
			{
				Command: "pnpm", Args: []string{"biome", "format", "--diagnostic-level=error", "--write"},
				Dir: ".", Generator: generatorName, Phase: generator.PhaseTypeCheck,
			},
		}
	},
}
