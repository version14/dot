package generators_typescript_backend_architecture_hexagonal

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

//go:embed all:files
var hexagonalFiles embed.FS

const generatorName = "typescript-backend-hexagonal"

// HexagonalTS scaffolds a TypeScript Express hexagonal (ports & adapters) project structure.
//
// Layout:
//
//	src/
//	  core/
//	    domain/         — entities, value objects, domain errors
//	    application/
//	      ports/in/     — driving ports (use case interfaces)
//	      ports/out/    — driven ports (repository / service interfaces)
//	      use-cases/    — implementations of driving ports
//	  adapters/
//	    primary/http/   — HTTP controllers and routes (driving adapters)
//	    secondary/      — DB repositories, external service clients (driven adapters)
//	  shared/           — cross-cutting: errors, libs, middlewares
var HexagonalTS = &scaffold.Generator{
	Name:     generatorName,
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		data := struct{ Name string }{s.Project.Name}

		ops, err := genfs.RenderDir(hexagonalFiles, "files", generatorName, data)
		if err != nil {
			return nil, fmt.Errorf("typescript-backend-hexagonal: %w", err)
		}
		return ops, nil
	},

	PostApplyFunction: func(s spec.Spec) []generator.PostOp {
		return []generator.PostOp{
			{Command: "pnpm", Args: []string{"install"}, Dir: ".", Generator: generatorName},
		}
	},
}
