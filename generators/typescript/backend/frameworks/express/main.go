package typescript_express_generator

import (
	"embed"
	"fmt"

	express_shared "github.com/version14/dot/generators/common/typescript/express/shared"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var expressFiles embed.FS

const generatorName = "typescript-express"

// ExpressTS generates Express-specific source files layered on top of an
// architecture generator. It dispatches to the correct files subdirectory
// based on the "ts-architecture" answer recorded in spec.Extensions.
var Generator = &scaffold.Generator{
	Name:     generatorName,
	Version:  "1.0.0",
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		arch, _ := s.Extensions["ts-architecture"].(string)
		if arch == "" {
			arch = "mvc"
		}

		data := struct{ Name string }{s.Project.Name}

		ops, err := genfs.RenderDir(expressFiles, "files/"+arch, generatorName, data)
		if err != nil {
			return nil, fmt.Errorf("typescript-express (%s): %w", arch, err)
		}

		// MVC has a flat src/ layout; Clean and Hexagonal use src/shared/.
		sharedPrefix := "src/shared"
		if arch == "mvc" {
			sharedPrefix = "src"
		}
		sharedOps, err := express_shared.SharedOps(sharedPrefix, generatorName)
		if err != nil {
			return nil, err
		}

		return append(ops, sharedOps...), nil
	},

	PostApplyFunction: func(s spec.Spec) []generator.PostOp {
		return []generator.PostOp{
			{
				Command: "pnpm", Args: []string{"install"},
				Dir: ".", Generator: generatorName, Phase: generator.PhaseInstall,
			},
			{
				Command: "pnpm", Args: []string{"exec", "tsc", "--noEmit"},
				Dir: ".", Generator: generatorName, Phase: generator.PhaseTypeCheck,
			},
			// Smoke: start dev server in background, then curl health endpoint.
			{
				Command: "pnpm", Args: []string{"run", "dev"},
				Dir: ".", Generator: generatorName, Phase: generator.PhaseSmoke, Background: true,
			},
			{
				Command: "curl", Args: []string{"-sf", "--retry", "5", "--retry-connrefused", "http://localhost:3000/health"},
				Dir: ".", Generator: generatorName, Phase: generator.PhaseSmoke,
			},
		}
	},
}
