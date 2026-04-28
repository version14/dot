package backendArchitectureHexagonal

import "github.com/version14/dot/pkg/dotapi"

// Manifest declares backend-architecture-hexagonal — the generator that scaffolds an
// Hexagonal Architecture-based backend structure.
var Manifest = dotapi.Manifest{
	Name:        "backend_architecture_hexagonal",
	Version:     "1.0.0",
	Description: "Scaffolds a base structure for a backend architecture using Hexagonal Architecture",
	Outputs: []string{
		"src/adapters/primary/http/controllers/.gitkeep",
		"src/adapters/primary/http/routes/.gitkeep",
		"src/adapters/secondary/external/.gitkeep",
		"src/adapters/secondary/persistance/repositories/.gitkeep",

		"src/core/application/ports/in/.gitkeep",
		"src/core/application/use-cases/.gitkeep",
		"src/core/domain/entities/.gitkeep",
		"src/core/domain/errors/.gitkeep",
		"src/core/domain/value-objects/.gitkeep",

		"src/shared/errors/.gitkeep",
		"src/shared/libs/.gitkeep",
		"src/shared/middlewares/.gitkeep",
	},
	Validators: []dotapi.Validator{
		{
			Name: "backend_architecture_hexagonal",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/adapters/primary/http/controllers/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/adapters/primary/http/routes/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/adapters/secondary/external/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/adapters/secondary/persistance/repositories/.gitkeep"},

				{Type: dotapi.CheckFileExists, Path: "src/core/application/ports/in/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/core/application/use-cases/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/core/domain/entities/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/core/domain/errors/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/core/domain/value-objects/.gitkeep"},

				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/libs/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/middlewares/.gitkeep"},
			},
		},
	},
}
