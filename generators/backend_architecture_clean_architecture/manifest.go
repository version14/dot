package backendArchitectureCleanArchitecture

import "github.com/version14/dot/pkg/dotapi"

// Manifest declares backend-architecture-clean-architecture — the generator that scaffolds an
// clean architecture-based backend structure.
var Manifest = dotapi.Manifest{
	Name:        "backend_architecture_clean_architecture",
	Version:     "1.0.0",
	Description: "Scaffolds a base structure for a backend architecture using Clean Architecture",
	Outputs: []string{
		"src/modules/example/application/controllers/.gitkeep",
		"src/modules/example/application/use-cases/.gitkeep",
		"src/modules/example/application/validators/.gitkeep",

		"src/modules/example/domain/entities/.gitkeep",
		"src/modules/example/domain/errors/.gitkeep",
		"src/modules/example/domain/interfaces/.gitkeep",

		"src/modules/example/infrastructure/database/repositories/.gitkeep",
		"src/modules/example/infrastructure/database/schemas/.gitkeep",
	},
	Validators: []dotapi.Validator{
		{
			Name: "backend_architecture_clean_architecture",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/modules/example/application/controllers/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/example/application/use-cases/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/example/application/validators/.gitkeep"},

				{Type: dotapi.CheckFileExists, Path: "src/modules/example/domain/entities/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/example/domain/errors/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/example/domain/interfaces/.gitkeep"},

				{Type: dotapi.CheckFileExists, Path: "src/modules/example/infrastructure/database/repositories/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/example/infrastructure/database/schemas/.gitkeep"},
			},
		},
	},
}
